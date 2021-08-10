package interfaces

import (
	"context"
	"fmt"

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
	var cmdOptions = []string{"Get Random Quote", "Get Based On Genres", "Get Based On Authors", "Exit"}
	prompt := promptui.Select{
		Label: "Command",
		Items: cmdOptions,
	}

	pterm.Println()
	newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgBlack),
		BackgroundStyle: pterm.NewStyle(pterm.BgRed),
		Margin:          20,
	}
	newHeader.Println("Make your choice...")
	_, result, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	switch result {
	case cmdOptions[0]:
		pterm.DefaultSection.Println("Getting Random Quote...")
		if err := t.RandomQuote(ctx); err != nil {
			return err
		}
	case cmdOptions[1]:
		pterm.DefaultSection.Println("Getting Quotes Based On Genres...")
		if err := t.SelectGenre(ctx); err != nil {
			return err
		}
	case cmdOptions[2]:
		pterm.DefaultSection.Println("Getting Quotes Based On Authors...")
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
	quote.Print()
	return nil
}

func (t *Term) SelectGenre(ctx context.Context) error {
	// make HTTP request
	items, err := t.source.AllGenres(ctx)
	if err != nil {
		return fmt.Errorf("could not get categories from source: %v", err)
	}

	// select from the response genres
	prompt := promptui.Select{
		Label: "Select Genre",
		Items: items,
	}

	_, selection, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	qo := &source.QueryOptions{
		Genre: selection,
	}
	t.GoFurther(ctx, qo)
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

	_, selection, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	qo := &source.QueryOptions{
		Author: selection,
	}
	t.GoFurther(ctx, qo)
	return nil
}

func (t *Term) GoFurther(ctx context.Context, qo *source.QueryOptions) error {
	for {
		t.printSection(qo)
		itemSelection := []string{"Show All", "Show Random Quote", "Exit"}
		prompt := promptui.Select{
			Label: "Select Action:",
			Items: itemSelection,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %v", err)
		}

		switch result {
		case itemSelection[0]:
			if err := t.ShowAll(ctx, qo); err != nil {
				return nil
			}
		case itemSelection[1]:
			// TODO: implement get random quote from selection
		case "Exit":
			return nil
		}
	}

}

func (t *Term) ShowAll(ctx context.Context, qo *source.QueryOptions) error {
	pageSelection := 1

	for {
		// query based on selection
		opt := &source.QueryOptions{
			Limit:  9,
			Page:   int32(pageSelection),
			Genre:  qo.Genre,
			Author: qo.Author,
		}
		quotes, pag, err := t.source.Quotes(ctx, opt)
		if err != nil {
			return fmt.Errorf("quotes query failed: %v", err)
		}
		title := fmt.Sprintf("PAGE %d/%d", pageSelection, pag.TotalPages)
		t.source.PrintQuotesPage(title, quotes)

		itemSelection := []string{"Next Page", "Previous Page", "Exit"}
		prompt := promptui.Select{
			Label: "Select Action:",
			Items: itemSelection,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %v", err)
		}
		switch result {
		case itemSelection[0]:
			pageSelection += 1
			if pageSelection >= pag.TotalPages {
				pageSelection = 1
			}
		case itemSelection[1]:
			if pageSelection == 1 {
				pageSelection = pag.TotalPages
			} else {
				pageSelection -= 1
			}
		case "Exit":
			return nil
		}
	}

}

func (t *Term) printSection(qo *source.QueryOptions) {
	// TODO: further refactoring needed
	if qo.Genre != "" {
		pterm.DefaultSection.Printf("Selected option: %s", qo.Genre)
	} else {
		pterm.DefaultSection.Printf("Selected option: %s", qo.Author)
	}
}
