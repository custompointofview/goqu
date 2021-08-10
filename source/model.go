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
	// prefPrint := pterm.PrefixPrinter{
	// 	MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
	// 	Prefix: pterm.Prefix{
	// 		Style: &pterm.ThemeDefault.InfoPrefixStyle,
	// 		Text:  strings.ToUpper(q.Genre),
	// 	},
	// }

	// // TODO: this is good, but not good enough, long text is not multiline - needs separation
	// pterm.DefaultParagraph.WithMaxWidth(120).Println(prefPrint.Sprintf("%s", q.Text))
	// pterm.DefaultParagraph.WithMaxWidth(120).Printf("-- %s\n", pterm.LightRed(q.Author))
	// pterm.Println()

	newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgWhite),
		BackgroundStyle: pterm.NewStyle(pterm.BgGray),
		Margin:          5,
	}

	// Print header.
	newHeader.Println(q.Sprint())
}
