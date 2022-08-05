package common

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/RocketChat/statuscentral/client"
	"github.com/RocketChat/statuscentral/client/oauthclient"
)

const oauthHost string = ""
const clientID string = ""

const scopes string = "offline_access workspace"

var oClient *oauthclient.OAuthClient

var debugMode = false

var codeChan = make(chan string)

func Login(baseURL, token string) error {
	state := State{}

	if token != "" && baseURL != "" {
		state.BaseURL = baseURL
		state.LoginToken = token
	} else if oauthHost != "" && clientID != "" {
		if err := initOAuthClient(oauthHost); err != nil {
			return err
		}

		if err := doAuthentication(); err != nil {
			return err
		}

		s, err := oClient.GetSessionInfo()
		if err != nil {
			return err
		}

		state.Session = s
	}

	if err := SaveState(state); err != nil {
		return err
	}

	return nil
}

func Logout() error {
	if err := initOAuthClient(oauthHost); err != nil {
		return err
	}

	if err := oClient.Revoke(); err != nil {
		return err
	}

	if err := DeleteState(); err != nil {
		return err
	}

	return nil
}

func GetStatusCentralClient() *client.Client {
	scClient, err := getClient()
	if err != nil {
		if debugMode {
			log.Println("Error Restoring Session State:", err)
		}

		exitWithLoginMessage()
	}

	if scClient == nil {
		exitWithLoginMessage()
	}

	return scClient
}

func getClient() (*client.Client, error) {
	state, err := LoadState()
	if err != nil {
		return nil, err
	}

	if state == nil {
		return nil, errors.New("no state found")
	}

	clientConfig := client.Config{
		BaseURL: state.BaseURL,
	}

	if state.Session != nil {
		if err := initOAuthClient(oauthHost); err != nil {
			return nil, err
		}

		if err := oClient.RestoreSession(*state.Session); err != nil {
			return nil, err
		}

		if !oClient.HasActiveSession() {
			log.Println("Invalid session")
			return nil, errors.New("State doesn't contain active session")
		}

		clientConfig.OAuthClient = oClient
	} else if state.LoginToken != "" {
		clientConfig.Token = state.LoginToken
	}

	c, err := client.New(clientConfig)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func exitWithLoginMessage() {
	log.Println("No state to restore from.\n\nTo Login use: statusctl login")
	os.Exit(0)
}

func doAuthentication() error {
	if err := oClient.NewSession(); err != nil {
		return err
	}

	go startServer()

	openbrowser(oClient.BuildAuthorizeURL())

	code := <-codeChan
	verboseLog("Got code back from oauth", code)

	if err := oClient.CompleteAuthorization(code); err != nil {
		return err
	}

	return nil
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)

	Router := gin.New()

	Router.GET("/callback", callbackHandler)

	Router.Run(":11899") //nolint:errcheck // Tech debt
}

func callbackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if !oClient.IsValidState(state) {
		c.String(http.StatusBadRequest, "Invalid state returned")
		c.Abort()
		os.Exit(1)
	}

	c.String(http.StatusOK, "You can close your browser window now")
	c.Abort()
	codeChan <- code
}

func openbrowser(url string) {
	logger("Opening Browser to perform OAuth Authorization. \n\nIf your browser doesn't open please proceed to the route below:\n\n", url)

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	// case "windows":
	// err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func initOAuthClient(URL string) error {
	c, err := oauthclient.New(oauthclient.ClientConfig{
		URL:         URL,
		ClientID:    clientID,
		Scope:       scopes,
		RedirectURI: "http://localhost:11899/callback",
		PKCE:        true,
	})

	oClient = c

	if err != nil {
		return err
	}

	return nil
}

func verboseLog(v ...interface{}) {
	if debugMode {
		logger(v...)
	}
}

func logger(v ...interface{}) {
	all := append([]interface{}{fmt.Sprintf("[%s]", time.Now().Format("01/02/2006 15:04:05"))}, v...)
	log.Println(all...)
}
