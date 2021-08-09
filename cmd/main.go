package main

import (
	"context"

	"github.com/custompointofview/goqu/interfaces"
)

func main() {
	ctx := context.Background()

	t := interfaces.NewTerm()

	go func() {
		t.Run(ctx)
	}()

	select {
	case err := <-t.Error:
		panic(err)
	case <-t.Done:
		return
	case <-ctx.Done():
		return
	}
}
