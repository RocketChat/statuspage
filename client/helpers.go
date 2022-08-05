package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) getToken() (token string, err error) {
	if c.oAuthClient != nil && c.oAuthClient.HasActiveSession() {
		token, err = c.oAuthClient.GetAccessToken("", false)
		if err != nil {
			return "", err
		}

		return token, nil
	}

	return c.token, nil
}

func (c *Client) buildRequest(method string, resourceURI string, body interface{}) (*http.Request, error) {
	rel, _ := url.Parse(resourceURI)

	c.debugLog(rel)

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	if c.oAuthClient != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	} else {
		req.Header.Add("Authorization", token)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("Unauthorized")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("Invalid Route.  Make sure your client is up to date")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 204 {
		errResp := &ErrorResponse{}

		err = json.NewDecoder(resp.Body).Decode(errResp)
		if err != nil {
			return nil, err
		}

		return nil, errResp
	}

	if v != nil && resp.StatusCode != 204 {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}
