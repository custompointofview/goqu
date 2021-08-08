package interfaces

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/custompointofview/goqu/source"
)

type Term struct {
	source source.Sources
}

func NewTerm() *Term {
	return &Term{
		source: source.NewQuoteGarden(),
	}
}

func (t *Term) SelectCommand(ctx context.Context) error {
	prompt := promptui.Select{
		Label: "Command",
		Items: []string{"Get All Genres"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	fmt.Printf("You choose %q\n", result)
	switch result {
	case "Get All Genres":
		if err := t.SelectCategory(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (t *Term) SelectCategory(ctx context.Context) error {
	// make HTTP request

	items, err := t.source.AllGenres(ctx)
	if err != nil {
		return fmt.Errorf("could not get categories from source: %v", err)
	}

	// select from the response genres
	prompt := promptui.Select{
		Label: "Select Category",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	fmt.Printf("You choose %q\n", result)
	return nil
}

func (t *Term) Default() {
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
			"Saturday", "Sunday"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
