// Package source implements different sources
package source

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pterm/pterm"
)

const (
	QUOTEGARDEN_URI = "https://quote-garden.herokuapp.com/api/v3"
)

type QuoteGarden struct {
	BaseURL    string
	HTTPClient *http.Client
}

type qgOptions struct {
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

func (qg *QuoteGarden) RandomQuote(ctx context.Context) (*Quote, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/quotes/random", qg.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res := &QGQuote{}
	if err := qg.sendRequest(req, res); err != nil {
		return nil, err
	}
	return res.Data[0], nil
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

func (qg *QuoteGarden) AllAuthors(ctx context.Context) ([]string, error) {
	limit := 100
	page := 1

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/authors?limit=%d&page=%d", qg.BaseURL, limit, page), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res := &QGAuthors{}
	if err := qg.sendRequest(req, res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *QuoteGarden) sendRequest(req *http.Request, v interface{}) (retErr error) {
	spinnerSuccess, _ := pterm.DefaultSpinner.Start("Sending request...")
	defer func() {
		if retErr != nil {
			spinnerSuccess.Fail()
			return
		}
		spinnerSuccess.Success()
	}()

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return retErr
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if retErr = json.NewDecoder(res.Body).Decode(&v); retErr != nil {
		return retErr
	}

	return nil
}
