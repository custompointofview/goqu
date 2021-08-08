package source

import "context"

type Sources interface {
	AllGenres(ctx context.Context) ([]string, error)
}
