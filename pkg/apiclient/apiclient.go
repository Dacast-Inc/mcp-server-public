package apiclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiClient struct {
	hostname string
	timeout  time.Duration
}

func NewApiClient(hostname string, timeout time.Duration) *ApiClient {
	return &ApiClient{
		hostname: hostname,
		timeout:  timeout,
	}
}

func (c *ApiClient) DoRequest(method string, endpoint string, params url.Values, bodyString *string, headers map[string]string) (*http.Response, []byte, error) {
	uri := c.hostname + endpoint

	client := &http.Client{
		Timeout: c.timeout,
	}

	var body io.Reader
	if bodyString != nil && (method == http.MethodPost || method == http.MethodPut) {
		body = strings.NewReader(*bodyString)
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.URL.RawQuery = params.Encode()

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, respBody, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, respBody, nil
}
