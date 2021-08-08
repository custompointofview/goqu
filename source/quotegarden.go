// Package source implements different sources
package source

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	QUOTEGARDEN_URI = "https://quote-garden.herokuapp.com/api/v3"
)

type QuoteGarden struct {
	BaseURL    string
	HTTPClient *http.Client
}

type QCOptions struct {
	author string
	genre  string
	query  string
	page   int32
	limit  int32
}

func NewQuoteGarden() *QuoteGarden {
	return &QuoteGarden{
		BaseURL: QUOTEGARDEN_URI,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (qg *QuoteGarden) AllGenres(ctx context.Context) ([]string, error) {
	limit := 100
	page := 1

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/genres?limit=%d&page=%d", qg.BaseURL, limit, page), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res := &QGGenres{}
	if err := qg.sendRequest(req, res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *QuoteGarden) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
