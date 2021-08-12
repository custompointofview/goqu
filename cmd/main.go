package main

import (
	"context"

	"github.com/custompointofview/goqu/interfaces"
)

func main() {
	// create application context
	ctx := context.Background()

	// run
	t := interfaces.NewTerm()
	t.Run(ctx)
}
