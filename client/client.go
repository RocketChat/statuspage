package client

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/RocketChat/statuscentral/client/oauthclient"
)

var defaultBaseURL = "http://localhost:5050"

// Config client config
type Config struct {
	BaseURL     string
	Token       string
	JwtKey      string
	OAuthClient *oauthclient.OAuthClient
}

// Client fleetcommand client
type Client struct {
	baseURL   *url.URL
	token     string
	userAgent string
	client    *http.Client
	jwtKey    []byte

	debug bool

	oAuthClient *oauthclient.OAuthClient
}

// New creates new fleetcommand client from config
func New(cfg Config) (*Client, error) {
	transport := &http.Transport{
		MaxIdleConns:    10,               //nolint:gomnd // Tech debt
		IdleConnTimeout: 30 * time.Second, //nolint:gomnd // Tech debt
	}

	if os.Getenv("http_proxy") != "" {
		url, err := url.Parse(os.Getenv("http_proxy"))
		if err != nil {
			log.Println(err)
		}

		transport.Proxy = http.ProxyURL(url)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}

	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL: baseURL,
		token:   cfg.Token,
		jwtKey:  []byte(cfg.JwtKey),

		userAgent: "StatusCentral Client",

		oAuthClient: cfg.OAuthClient,

		client: &http.Client{Transport: transport},
	}

	return c, nil
}

// DebugMode turns debug mode on
func (c *Client) DebugMode() {
	c.debug = true
}

func (c *Client) debugLog(v ...interface{}) {
	if c.debug {
		log.Println(v...)
	}
}

// Incidents incident methods
func (c *Client) Incidents() IncidentsInterface {
	return &incidents{client: c}
}

// ScheduledMaintenance maintenance methods
func (c *Client) ScheduledMaintenance() ScheduledMaintenanceInterface {
	return &scheduledMaintenance{client: c}
}

// Services service methods
func (c *Client) Services() ServicesInterface {
	return &services{client: c}
}
