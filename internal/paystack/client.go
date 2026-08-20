package paystack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/12sub/Reset/internal/cache"
)

type Client struct {
	secretKey string
	baseURL   string
	client    *http.Client
	cache     *cache.RedisCache
}

func New(secretKey string, c *cache.RedisCache) *Client {
	return &Client{
		secretKey: secretKey,
		baseURL:   "https://api.paystack.co",
		client:    &http.Client{Timeout: 10 * time.Second},
		cache:     c,
	}
}

func (c *Client) do(method, path string, body, result interface{}) error {
	var bodyReader *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("paystack error: status %d", resp.StatusCode)
	}
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) VerifyTransaction(ctx context.Context, reference string) (*VerifyResponse, error) {
	cacheKey := "tx:" + reference
	var cached VerifyResponse

	found, err := c.cache.Get(ctx, cacheKey, &cached)
	if err != nil {
		// Log cache error but continue to API
	}
	if found {
		return &cached, nil
	}

	var resp struct {
		Status  bool           `json:"status"`
		Message string         `json:"message"`
		Data    VerifyResponse `json:"data"`
	}
	if err := c.do("GET", "/transaction/verify/"+reference, nil, &resp); err != nil {
		return nil, err
	}
	if !resp.Status {
		return nil, fmt.Errorf("paystack: %s", resp.Message)
	}

	// Cache transaction for 30 minutes (immutable data)
	c.cache.Set(ctx, cacheKey, &resp.Data, 30*time.Minute)
	return &resp.Data, nil
}

func (c *Client) CreateSubscription(customerEmail, planCode, authorization string) (*SubscriptionResponse, error) {
	payload := map[string]string{
		"customer":      customerEmail,
		"plan":          planCode,
		"authorization": authorization,
	}
	var resp struct {
		Status  bool                 `json:"status"`
		Message string               `json:"message"`
		Data    SubscriptionResponse `json:"data"`
	}
	if err := c.do("POST", "/subscription", payload, &resp); err != nil {
		return nil, err
	}
	if !resp.Status {
		return nil, fmt.Errorf("paystack: %s", resp.Message)
	}
	return &resp.Data, nil
}

func (c *Client) DisableSubscription(subscriptionCode, token string) error {
	payload := map[string]string{
		"code":  subscriptionCode,
		"token": token,
	}
	var resp struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
	if err := c.do("POST", "/subscription/disable", payload, &resp); err != nil {
		return err
	}
	if !resp.Status {
		return fmt.Errorf("paystack: %s", resp.Message)
	}
	return nil
}

type VerifyResponse struct {
	Status        string `json:"status"`
	Reference     string `json:"reference"`
	Amount        int    `json:"amount"`
	Authorization struct {
		AuthorizationCode string `json:"authorization_code"`
	} `json:"authorization"`
	Customer struct {
		Email string `json:"email"`
	} `json:"customer"`
	Plan struct {
		PlanCode string `json:"plan_code"`
	} `json:"plan"`
}

type SubscriptionResponse struct {
	SubscriptionCode string `json:"subscription_code"`
	EmailToken       string `json:"email_token"`
	Status           string `json:"status"`
}