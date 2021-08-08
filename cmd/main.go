package main

import (
	"context"

	"github.com/custompointofview/goqu/interfaces"
)

func main() {
	ctx := context.Background()

	t := interfaces.NewTerm()
	if err := t.SelectCommand(ctx); err != nil {
		panic(err)
	}
}
