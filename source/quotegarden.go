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

	res := &QGQuotes{}
	if err := qg.sendRequest(req, res); err != nil {
		return nil, err
	}
	return res.Data[0].ToQuote(), nil
}

func (qg *QuoteGarden) AllGenres(ctx context.Context) ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/genres", qg.BaseURL), nil)
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

func (qg *QuoteGarden) Quotes(ctx context.Context, options *QueryOptions) ([]*Quote, *Pagination, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/quotes?%s", qg.BaseURL, options.Sprint()),
		nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)
	res := &QGQuotes{}
	if err := qg.sendRequest(req, res); err != nil {
		return nil, nil, err
	}
	return res.DataToQuotes(), &res.Pagination, nil
}

func (qg *QuoteGarden) PrintQuotesPage(title string, quotes []*Quote) {
	maxNumColumns := 3
	panels := make(pterm.Panels, 9)

	row := 0
	col := 0
	panels[row] = make([]pterm.Panel, maxNumColumns)
	for _, q := range quotes {
		p := pterm.DefaultBox.Sprint(q.Sprint())
		panel := pterm.Panel{Data: p}

		panels[row][col] = panel
		col += 1
		if col >= maxNumColumns {
			row += 1
			col = 0
			panels[row] = make([]pterm.Panel, maxNumColumns)
		}
	}

	pRender, _ := pterm.DefaultPanel.WithPanels(panels).Srender()
	pterm.DefaultBox.WithTitle(title).WithTitleBottomRight().WithRightPadding(0).WithBottomPadding(0).Println(pRender)
}

func (qg *QuoteGarden) sendRequest(req *http.Request, v interface{}) (retErr error) {
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

	res, err := qg.HTTPClient.Do(req)
	if err != nil {
		return err
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
