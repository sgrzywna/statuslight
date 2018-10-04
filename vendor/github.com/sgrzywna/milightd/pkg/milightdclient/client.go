package milightdclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sgrzywna/milightd/pkg/models"
)

// Client represents HTTP client for the milightd daemon.
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

// SetLight controls mi-light device through milightd daemon.
func (c *Client) SetLight(l models.Light) error {
	url := fmt.Sprintf("%s/api/v1/light", c.url)

	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
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
		return fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetSequences returns list of defined sequences from milightd daemon.
func (c *Client) GetSequences() ([]models.Sequence, error) {
	url := fmt.Sprintf("%s/api/v1/sequence", c.url)

	req, err := http.NewRequest("GET", url, nil)
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
		return nil, fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	var sequences []models.Sequence

	err = json.NewDecoder(resp.Body).Decode(&sequences)
	if err != nil {
		return nil, err
	}

	return sequences, nil
}

// AddSequence adds sequence through milightd daemon.
func (c *Client) AddSequence(seq models.Sequence) error {
	url := fmt.Sprintf("%s/api/v1/sequence", c.url)

	data, err := json.Marshal(seq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetSequence return sequence definition from milightd daemon.
func (c *Client) GetSequence(name string) (*models.Sequence, error) {
	url := fmt.Sprintf("%s/api/v1/sequence/%s", c.url, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	var seq models.Sequence

	err = json.NewDecoder(resp.Body).Decode(&seq)
	if err != nil {
		return nil, err
	}

	return &seq, nil
}

// DeleteSequence deletes sequence through milightd daemon.
func (c *Client) DeleteSequence(name string) error {
	url := fmt.Sprintf("%s/api/v1/sequence/%s", c.url, name)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetSequenceState returns state of the running sequence from milightd daemon.
func (c *Client) GetSequenceState() (*models.SequenceState, error) {
	url := fmt.Sprintf("%s/api/v1/seqctrl", c.url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	var state models.SequenceState

	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

// SetSequenceState control state of the running sequence through milightd daemon.
func (c *Client) SetSequenceState(state models.SequenceState) error {
	url := fmt.Sprintf("%s/api/v1/seqctrl", c.url)

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
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
		return fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
