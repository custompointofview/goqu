package source

import (
	"context"
)

type Sources interface {
	RandomQuote(ctx context.Context) (*Quote, error)
	AllGenres(ctx context.Context) ([]string, error)
	AllAuthors(ctx context.Context) ([]string, error)
	Quotes(ctx context.Context, options *QueryOptions) ([]*Quote, *Pagination, error)
	PrintQuotesPage(title string, quotes []*Quote, columns int)
}
