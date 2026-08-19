package paystack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	secretKey string
	baseURL   string
	client    *http.Client
}

func New(secretKey string) *Client {
	return &Client{
		secretKey: secretKey,
		baseURL:   "https://api.paystack.co",
		client:    &http.Client{Timeout: 10 * time.Second},
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
		return fmt.Errorf("paystack error: %s", resp.Status)
	}
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// CreateSubscription creates a Paystack subscription
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

// DisableSubscription cancels a subscription on Paystack
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

type SubscriptionResponse struct {
	SubscriptionCode string `json:"subscription_code"`
	EmailToken       string `json:"email_token"`
	Status           string `json:"status"`
}