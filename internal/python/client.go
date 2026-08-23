package python

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type ClassifyRequest struct {
	Text string `json:"text"`
}

type ClassifyResponse struct {
	RawText        string  `json:"raw_text"`
	Merchant       string  `json:"merchant"`
	Category       string  `json:"category"`
	Amount         float64 `json:"amount"`
	IsSubscription bool    `json:"is_subscription"`
	Confidence     float64 `json:"confidence"`
	MatchType      string  `json:"match_type"`
}

func (c *Client) Classify(ctx context.Context, text string) (*ClassifyResponse, error) {
	fmt.Printf("DEBUG: Calling Python at %s/classify-transaction\n", c.baseURL)
	body, _ := json.Marshal(ClassifyRequest{Text: text})
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/classify-transaction", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python service returned %d", resp.StatusCode)
	}

	var result ClassifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}