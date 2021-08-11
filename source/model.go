package source

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pterm/pterm"
)

type QueryOptions struct {
	Author string
	Genre  string
	Query  string
	Page   int32
	Limit  int32
}

func (qgp *QueryOptions) Sprint() string {
	var retElem []string
	if qgp.Author != "" {
		retElem = append(retElem, fmt.Sprintf("author=%s", url.QueryEscape(qgp.Author)))
	}
	if qgp.Genre != "" {
		retElem = append(retElem, fmt.Sprintf("genre=%s", url.QueryEscape(qgp.Genre)))
	}
	if qgp.Query != "" {
		retElem = append(retElem, fmt.Sprintf("query=%s", url.QueryEscape(qgp.Query)))
	}
	if qgp.Limit != 0 {
		retElem = append(retElem, fmt.Sprintf("limit=%d", qgp.Limit))
	}
	if qgp.Page != 0 {
		retElem = append(retElem, fmt.Sprintf("page=%d", qgp.Page))
	}
	return strings.Join(retElem, "&")
}

type Quote struct {
	ID     string
	Text   string
	Author string
	Genre  string
}

func (q *Quote) Sprint() string {
	return fmt.Sprintf("%s \n---------------\n%s \n-- %s", strings.ToUpper(q.Genre),
		pterm.DefaultParagraph.WithMaxWidth(60).Sprintln(q.Text),
		q.Author)
}

func (q *Quote) Print() {
	header := q.header()
	header.Println(q.Sprint())
}

func (q *Quote) HSprint() string {
	header := q.header()
	return header.Sprintln(q.Sprint())
}

func (q *Quote) header() pterm.HeaderPrinter {
	return pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgWhite),
		BackgroundStyle: pterm.NewStyle(pterm.BgGray),
		Margin:          1,
	}
}
