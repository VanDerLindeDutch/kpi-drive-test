package http_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type BasicAuth struct {
	Username string
	Password string
}

type Client struct {
	client    *http.Client
	baseUrl   string
	basicAuth *BasicAuth
}

func NewClient(baseUrl string, auth *BasicAuth) *Client {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &Client{
		client:    &http.Client{Transport: transport},
		baseUrl:   baseUrl,
		basicAuth: auth,
	}
}

func (c *Client) CookieAuthorize(ctx context.Context, path string, headers map[string]string, body url.Values, respPtr any) ([]*http.Cookie, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseUrl+path, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if c.basicAuth != nil {
		req.SetBasicAuth(c.basicAuth.Username, c.basicAuth.Password)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if respPtr != nil {
		err = json.NewDecoder(resp.Body).Decode(respPtr)
		if err != nil {
			return nil, 0, err
		}
	}
	return resp.Cookies(), resp.StatusCode, nil
}

func (c *Client) PostFormDataCookie(ctx context.Context, path string, headers map[string]string, body url.Values, respPtr any, cookie *http.Cookie) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseUrl+path, strings.NewReader(body.Encode()))
	if err != nil {
		return 0, err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if c.basicAuth != nil {
		req.SetBasicAuth(c.basicAuth.Username, c.basicAuth.Password)
	}
	req.AddCookie(cookie)
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if respPtr != nil {
		err = json.NewDecoder(resp.Body).Decode(respPtr)
		if err != nil {
			return 0, err
		}
	}
	return resp.StatusCode, nil
}

func (c *Client) GetJsonCookie(ctx context.Context, path string, headers map[string]string, body []byte, respPtr any, cookie *http.Cookie) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+path, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	req.Header.Add("content-type", "application/json")
	if c.basicAuth != nil {
		req.SetBasicAuth(c.basicAuth.Username, c.basicAuth.Password)
	}
	req.AddCookie(cookie)
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if respPtr != nil {
		err = json.NewDecoder(resp.Body).Decode(respPtr)
		if err != nil {
			return 0, err
		}
	}
	return resp.StatusCode, nil
}

func (c *Client) GetJson(ctx context.Context, path string, headers map[string]string, respPtr any) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+path, nil)
	if err != nil {
		return 0, err
	}
	//req.Header.Add("content-type", c.contentType)
	if c.basicAuth != nil {
		req.SetBasicAuth(c.basicAuth.Username, c.basicAuth.Password)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if respPtr != nil {
		err = json.NewDecoder(resp.Body).Decode(respPtr)
		if err != nil {
			return 0, err
		}
	}
	return resp.StatusCode, nil
}

func (c *Client) PostJson(ctx context.Context, path string, headers map[string]string, body []byte, respPtr any) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseUrl+path, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	//req.Header.Add("content-type", c.contentType)
	if c.basicAuth != nil {
		req.SetBasicAuth(c.basicAuth.Username, c.basicAuth.Password)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if respPtr != nil {
		err = json.NewDecoder(resp.Body).Decode(respPtr)
		if err != nil {
			return 0, err
		}
	}
	return resp.StatusCode, nil
}
