package interfaces

import (
	"context"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"

	"github.com/custompointofview/goqu/source"
)

type Term struct {
	source source.Sources
	Error  chan error
	Done   chan bool
}

func NewTerm() *Term {
	return &Term{
		source: source.NewQuoteGarden(),
		Error:  make(chan error),
		Done:   make(chan bool),
	}
}

func (t *Term) Run(ctx context.Context) {
	for {
		err := t.SelectCommand(ctx)
		if err != nil {
			t.Error <- err
			return
		}
		select {
		case <-t.Done:
			t.Done <- true
			return
		default:
			continue
		}
	}
}

func (t *Term) SelectCommand(ctx context.Context) error {
	// TODO: Commands should be constants or in an object. Too many hardcoded locations.
	prompt := promptui.Select{
		Label: "Command",
		Items: []string{"Get Random Quote", "Get All Genres", "Get All Authors", "Exit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	switch result {
	case "Get Random Quote":
		pterm.DefaultSection.Println("Getting Random Quote...")
		if err := t.RandomQuote(ctx); err != nil {
			return err
		}
	case "Get All Genres":
		pterm.DefaultSection.Println("Getting Genres...")
		if err := t.SelectCategory(ctx); err != nil {
			return err
		}
	case "Get All Authors":
		pterm.DefaultSection.Println("Getting Authors...")
		if err := t.SelectAuthor(ctx); err != nil {
			return err
		}
	case "Exit":
		t.Done <- true
	}

	return nil
}

func (t *Term) RandomQuote(ctx context.Context) error {
	// make HTTP request
	quote, err := t.source.RandomQuote(ctx)
	if err != nil {
		return fmt.Errorf("could not get categories from source: %v", err)
	}

	// TODO: quote styling needs to be refactored!
	prefPrint := pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.InfoPrefixStyle,
			Text:  strings.ToUpper(quote.QuoteGenre),
		},
	}

	// TODO: this is good, but not good enough, long text is not multiline - needs separation
	prefPrint.Printf("%s\n-----\n%s\n", quote.QuoteText, quote.QuoteAuthor)
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
		return fmt.Errorf("prompt failed: %v", err)
	}

	fmt.Printf("You choose %q\n", result)
	return nil
}

func (t *Term) SelectAuthor(ctx context.Context) error {
	// make HTTP request
	items, err := t.source.AllAuthors(ctx)
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
		return fmt.Errorf("prompt failed: %v", err)
	}

	fmt.Printf("You choose %q\n", result)
	return nil
}
