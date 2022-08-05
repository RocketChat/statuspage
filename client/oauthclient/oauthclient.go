package oauthclient

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuthClient a new oauth client to talk with fleetcommand
type OAuthClient struct {
	config  ClientConfig
	session ClientSession
}

// ClientConfig configuration for the OAuthClient
type ClientConfig struct {
	URL          string
	ClientID     string
	ClientSecret string
	PKCE         bool
	Scope        string
	RedirectURI  string
}

type ClientSession struct {
	state            string    `json:"-"`
	codeVerifier     string    `json:"-"`
	codeChallenge    string    `json:"-"`
	ExpiresAt        time.Time `json:"expiresAt"`
	AccessToken      string    `json:"accessToken"`
	RefreshToken     string    `json:"refreshToken"`
	AuthorizedScopes string    `json:"authorizedScopes"`
}

type tokenResponse struct {
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

// New Returns a new OAuthClient
func New(config ClientConfig) (*OAuthClient, error) {
	client := &OAuthClient{
		config:  config,
		session: ClientSession{},
	}

	if client.config.URL == "" {
		client.config.URL = "http://localhost:5050"
	}

	if client.config.ClientID == "" {
		return nil, errors.New("ClientID must be set")
	}

	if client.config.Scope == "" {
		return nil, errors.New("Scope must be set")
	}

	if client.config.RedirectURI == "" {
		return nil, errors.New("RedirectURI must be set")
	}

	if client.config.ClientSecret == "" && !client.config.PKCE {
		return nil, errors.New("ClientSecret must be set unless using PKCE")
	}

	return client, nil
}

// NewSession initializes a new session
func (o *OAuthClient) NewSession() error {
	state, err := newUUID()
	if err != nil {
		return err
	}

	o.session = ClientSession{
		state: state,
	}

	if o.config.PKCE {
		codeVerifier, codeChallenge, err := generateCodeChallenge()
		if err != nil {
			return err
		}

		o.session.codeChallenge = codeChallenge
		o.session.codeVerifier = codeVerifier
	}

	return nil
}

func (o *OAuthClient) GetSessionInfo() (*ClientSession, error) {
	if !o.HasActiveSession() {
		return nil, errors.New("No Active Session")
	}

	return &o.session, nil
}

func (o *OAuthClient) RestoreSession(session ClientSession) error {
	o.session = session

	if !o.HasActiveSession() {
		return errors.New("Not Valid Session")
	}

	return nil
}

// BuildAuthorizeURL builds the authorize url the user would visit to grant authorization
func (o *OAuthClient) BuildAuthorizeURL() string {
	// Make spaces link safe
	scope := strings.ReplaceAll(o.config.Scope, " ", "%20")

	authorizeURL := o.config.URL + "/authorize" +
		"?client_id=" + o.config.ClientID +
		"&response_type=code" +
		"&access_type=offline" + // Actually not sure what cases this needs passed?
		"&redirect_uri=" + o.config.RedirectURI +
		"&scope=" + scope +
		"&state=" + o.session.state

	if o.config.PKCE {
		authorizeURL = authorizeURL + "&code_challenge_method=S256" +
			"&code_challenge=" + o.session.codeChallenge
	}

	return authorizeURL
}

// CompleteAuthorization exchanges authorization code for access token and refresh (if offline_access scope was used)
func (o *OAuthClient) CompleteAuthorization(code string) error {
	formBody := url.Values{
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {o.config.RedirectURI},
		"client_id":     {o.config.ClientID},
		"client_secret": {o.config.ClientSecret},
		"code":          {code},
	}

	if o.config.PKCE {
		formBody.Add("code_verifier", o.session.codeVerifier)
	}

	response, err := http.PostForm(o.config.URL+"/api/oauth/token", formBody)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		return errors.New(string(body))
	}

	var tokens tokenResponse

	if err := json.Unmarshal(body, &tokens); err != nil {
		return err
	}

	if tokens.AccessToken == "" {
		return errors.New("No accessToken returned")
	}

	// Add expires seconds to current time - minus 10 seconds so we never try to use an almost expired access token
	o.session.ExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn-10) * time.Second)

	o.session.RefreshToken = tokens.RefreshToken
	o.session.AccessToken = tokens.AccessToken
	o.session.AuthorizedScopes = tokens.Scope

	return nil
}

func (o *OAuthClient) refreshAccessToken(scope string, store bool) (accessToken string, err error) {
	if o.session.RefreshToken == "" {
		return "", errors.New("No refresh token available to use for refresh")
	}

	response, err := http.PostForm(o.config.URL+"/api/oauth/token", url.Values{
		"grant_type":    {"refresh_token"},
		"redirect_uri":  {o.config.RedirectURI},
		"client_id":     {o.config.ClientID},
		"client_secret": {o.config.ClientSecret},
		"scope":         {scope},
		"refresh_token": {o.session.RefreshToken},
	})

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 201 {
		return "", errors.New(string(body))
	}

	var token tokenResponse

	if err := json.Unmarshal(body, &token); err != nil {
		return "", err
	}

	if store {
		o.session.AccessToken = token.AccessToken

		// Add expires seconds to current time - minus 10 seconds so we never try to use an almost expired access token
		o.session.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn-10) * time.Second)
	}

	return token.AccessToken, err
}

// GetAccessToken if current token has desired scopes and hasn't expired will be returned, otherwise refreshToken if available will be used to obtain a valid accessToken
func (o *OAuthClient) GetAccessToken(limitScope string, forceRefresh bool) (accessToken string, err error) {
	refresh := forceRefresh
	store := true
	scope := o.session.AuthorizedScopes

	if limitScope != "" && o.session.AuthorizedScopes != limitScope {
		if !strings.Contains(scope, limitScope) {
			return "", errors.New("Requesting a scope new scope")
		}

		scope = limitScope
		store = false
		refresh = true
	}

	if o.session.ExpiresAt.Before(time.Now()) || refresh {
		accessToken, err = o.refreshAccessToken(scope, store)
		if err != nil {
			return "", err
		}

		return accessToken, nil
	}

	return o.session.AccessToken, nil
}

// Revoke invalidates the token or logs out
func (o *OAuthClient) Revoke() error {
	// Nothing really to do since they don't have a refresh token
	if o.session.RefreshToken == "" {
		return nil
	}

	response, err := http.PostForm(o.config.URL+"/api/oauth/revoke", url.Values{
		"client_id":       {o.config.ClientID},
		"client_secret":   {o.config.ClientSecret},
		"token":           {o.session.RefreshToken},
		"token_type_hint": {"refresh_token"},
	})

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New(string(body))
	}

	return nil
}

// IsValidState compares state to stored state to see if valid
func (o *OAuthClient) IsValidState(state string) bool {
	return o.session.state == state
}

// HasActiveSession returns true if there is an active session
func (o *OAuthClient) HasActiveSession() bool {
	if o.session.AccessToken == "" {
		return false
	}

	if o.session.ExpiresAt.Before(time.Now()) && o.session.RefreshToken == "" {
		return false
	}

	return true
}

func generateCodeChallenge() (codeVerifier string, codeChallenge string, err error) {
	codeVerifier, err = newUUID()
	if err != nil {
		return "", "", err
	}

	sha256Hash := sha256.New()
	if _, err := sha256Hash.Write([]byte(codeVerifier)); err != nil {
		return "", "", err
	}

	sha := sha256Hash.Sum(nil)

	codeChallenge = base64.StdEncoding.EncodeToString(sha)

	// Base64Url - https://base64.guru/standards/base64url
	codeChallenge = strings.ReplaceAll(codeChallenge, "+", "-")
	codeChallenge = strings.ReplaceAll(codeChallenge, "/", "_")
	codeChallenge = strings.ReplaceAll(codeChallenge, "=", "")

	return codeVerifier, codeChallenge, nil
}

// newUUID generates a random UUID according to the RFC 4122, https://play.golang.org/p/4FkNSiUDMg
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)

	if n != len(uuid) || err != nil {
		return "", err
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
