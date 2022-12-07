package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tie.prodigy9.co/config"
	"tie.prodigy9.co/domain"
)

type Client struct {
	http       *http.Client
	apiPrefix  string
	adminToken string
}

func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{
		http:       &http.Client{},
		apiPrefix:  cfg.APIPrefix(),
		adminToken: cfg.AdminToken(),
	}, nil
}

func (c *Client) GetTies() ([]*domain.Tie, error) {
	list := &domain.List[*domain.Tie]{}
	if err := c.do(list, "GET", c.apiPrefix+"/ties", nil); err != nil {
		return nil, err
	} else {
		return list.Data, nil
	}
}

func (c *Client) CreateTie(slug string, target string) (*domain.Tie, error) {
	buf, err := json.Marshal(&domain.CreateTie{
		Slug:      slug,
		TargetURL: target,
	})
	if err != nil {
		return nil, err
	}

	tie := &domain.Tie{}
	if err := c.do(tie, "POST", c.apiPrefix+"/ties", bytes.NewBuffer(buf)); err != nil {
		return nil, err
	} else {
		return tie, nil
	}
}

func (c *Client) DeleteTie(slug string) (*domain.Tie, error) {
	tie := &domain.Tie{}
	if err := c.do(tie, "DELETE", c.apiPrefix+"/ties/"+slug, nil); err != nil {
		return nil, err
	} else {
		return tie, nil
	}
}

func (c *Client) do(out any, method, path string, payload io.Reader) error {
	req, err := http.NewRequest(method, path, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.adminToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return err
		}
	}
	return nil
}
