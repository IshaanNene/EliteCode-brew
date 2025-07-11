package api

import (
	"encoding/json"
	"fmt"
	"time"

	"elitecode/internal/storage"
	"elitecode/internal/utils"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client  *resty.Client
	baseURL string
	token   string
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewClient() *Client {
	config := storage.GetConfig()
	
	client := resty.New().
		SetBaseURL(config.APIBaseURL).
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetHeader("User-Agent", "Elitecode-CLI/1.0").
		SetHeader("Content-Type", "application/json")

	// Add auth token if available
	if config.AuthToken != "" {
		client.SetAuthToken(config.AuthToken)
	}

	// Add request/response logging in debug mode
	if config.Debug {
		client.SetDebug(true)
	}

	// Add error handling
	client.OnError(func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			utils.Logger.Error("API request failed", 
				"status", v.Response.Status(),
				"body", string(v.Response.Body()),
			)
		} else {
			utils.Logger.Error("API request failed", "error", err)
		}
	})

	return &Client{
		client:  client,
		baseURL: config.APIBaseURL,
		token:   config.AuthToken,
	}
}

func (c *Client) SetAuthToken(token string) {
	c.token = token
	c.client.SetAuthToken(token)
}

func (c *Client) Get(endpoint string, result interface{}) error {
	resp, err := c.client.R().
		SetResult(result).
		Get(endpoint)

	if err != nil {
		return fmt.Errorf("GET request failed: %w", err)
	}

	if resp.IsError() {
		return c.handleErrorResponse(resp)
	}

	return nil
}

func (c *Client) Post(endpoint string, body interface{}, result interface{}) error {
	req := c.client.R()
	
	if body != nil {
		req.SetBody(body)
	}
	
	if result != nil {
		req.SetResult(result)
	}

	resp, err := req.Post(endpoint)
	if err != nil {
		return fmt.Errorf("POST request failed: %w", err)
	}

	if resp.IsError() {
		return c.handleErrorResponse(resp)
	}

	return nil
}

func (c *Client) Put(endpoint string, body interface{}, result interface{}) error {
	req := c.client.R()
	
	if body != nil {
		req.SetBody(body)
	}
	
	if result != nil {
		req.SetResult(result)
	}

	resp, err := req.Put(endpoint)
	if err != nil {
		return fmt.Errorf("PUT request failed: %w", err)
	}

	if resp.IsError() {
		return c.handleErrorResponse(resp)
	}

	return nil
}

func (c *Client) Delete(endpoint string, result interface{}) error {
	req := c.client.R()
	
	if result != nil {
		req.SetResult(result)
	}

	resp, err := req.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %w", err)
	}

	if resp.IsError() {
		return c.handleErrorResponse(resp)
	}

	return nil
}

func (c *Client) handleErrorResponse(resp *resty.Response) error {
	var errorResp ErrorResponse
	if err := json.Unmarshal(resp.Body(), &errorResp); err != nil {
		return fmt.Errorf("API error (%d): %s", resp.StatusCode(), string(resp.Body()))
	}

	return fmt.Errorf("API error (%d): %s", errorResp.Code, errorResp.Message)
}

func (c *Client) HealthCheck() error {
	var result map[string]interface{}
	return c.Get("/health", &result)
}