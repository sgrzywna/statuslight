package statuslightclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sgrzywna/statuslight/internal/app/statuslight"
)

// Client represents HTTP client to control status light daemon.
type Client struct {
	url    string
	client *http.Client
}

// NewClient returns initialized Client object.
func NewClient(url string) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// SetStatus sets status on remote status light daemon.
func (c *Client) SetStatus(id string, status bool) error {
	s := statuslight.Status{
		State: status,
		ID:    id,
	}

	d, err := json.Marshal(s)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/status", c.url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statuslight client: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
