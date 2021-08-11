package interfaces

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

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
	var cmdOptions = []string{"Get Random Quote", "Get Based On Genres",
		"Get Based On Authors", "Search...", "Exit"}
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
		pterm.DefaultSection.Println("Retrieving Random Quote...")
		if err := t.RandomQuote(ctx); err != nil {
			return err
		}
	case cmdOptions[1]:
		pterm.DefaultSection.Println("Retrieving Quotes Based On Genres...")
		if err := t.SelectGenre(ctx); err != nil {
			return err
		}
	case cmdOptions[2]:
		pterm.DefaultSection.Println("Retrieving Quotes Based On Authors...")
		if err := t.SelectAuthor(ctx); err != nil {
			return err
		}
	case cmdOptions[3]:
		pterm.DefaultSection.Println("Retrieving Quotes Based On Input...")
		if err := t.PromptInput(ctx, nil); err != nil {
			return err
		}
	case "Exit":
		t.Done <- true
	}

	return nil
}

func (t *Term) RandomQuote(ctx context.Context) error {
	quote, err := t.source.RandomQuote(ctx)
	if err != nil {
		return fmt.Errorf("could not get random quote from source: %v", err)
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

	// TODO: refactor the things below
	// select from the response genres
	prompt := promptui.Select{
		Label: "Select Genre",
		Items: items,
	}

	_, selection, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}
	// create query options & go further
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

	// TODO: refactor the things below
	// select from the response genres
	prompt := promptui.Select{
		Label: "Select Category",
		Items: items,
	}

	_, selection, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}
	// create query options & go further
	qo := &source.QueryOptions{
		Author: selection,
	}
	t.GoFurther(ctx, qo)
	return nil
}

func (t *Term) PromptInput(ctx context.Context, qo *source.QueryOptions) error {
	validate := func(input string) error {
		if input == "" {
			return fmt.Errorf("search must not be empty")
		}
		space := regexp.MustCompile(`\s+`)
		if space.Match([]byte(input)) {
			return fmt.Errorf("search must contain a single term")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Search:",
		Templates: templates,
		Validate:  validate,
	}

	selection, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed %v", err)
	}

	// create query options & go further
	qoTemp := &source.QueryOptions{
		Query: selection,
	}
	if qo != nil {
		qoTemp.Author = qo.Author
		qoTemp.Genre = qo.Genre
	}
	t.GoFurther(ctx, qoTemp)

	return nil
}

func (t *Term) GoFurther(ctx context.Context, qo *source.QueryOptions) error {
	for {
		t.printSection(qo)
		itemSelection := []string{"Show All", "Show Random Quote", "Search...", "Exit"}
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
				return err
			}
		case itemSelection[1]:
			if err := t.ShowRandom(ctx, qo); err != nil {
				return err
			}
		case itemSelection[2]:
			if err := t.PromptInput(ctx, qo); err != nil {
				return err
			}
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
			Query:  qo.Query,
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

func (t *Term) ShowRandom(ctx context.Context, qo *source.QueryOptions) error {
	pageSelection := 1
	for {
		// query based on selection
		opt := &source.QueryOptions{
			Limit:  9,
			Page:   int32(pageSelection),
			Genre:  qo.Genre,
			Author: qo.Author,
			Query:  qo.Query,
		}
		quotes, pag, err := t.source.Quotes(ctx, opt)
		if err != nil {
			return fmt.Errorf("quotes query failed: %v", err)
		}
		// get random quote
		rand.Seed(time.Now().Unix())
		randQ := quotes[rand.Intn(len(quotes))]
		randQ.Print()
		// change page
		pageSelection = rand.Intn(pag.TotalPages) + 1

		// after menu
		itemSelection := []string{"Get Another", "Exit"}
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
			// just continue with another request
		case "Exit":
			return nil
		}
	}
}

func (t *Term) printSection(qo *source.QueryOptions) {
	// TODO: further refactoring needed
	options := []string{qo.Genre, qo.Author, qo.Query}
	pterm.DefaultSection.Printf("Selected option: %s", strings.Join(options, " "))
}
