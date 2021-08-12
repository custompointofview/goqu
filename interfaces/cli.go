package interfaces

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"

	"github.com/custompointofview/goqu/assets"
	"github.com/custompointofview/goqu/source"
)

const GO_BACK = "< Go back"

var DEFAULT_SOURCE = source.NewQuoteGarden()
var DEFAULT_SOURCE_LIMIT = 9

type Term struct {
	wg          sync.WaitGroup
	Error       chan error
	Done        chan bool
	source      source.Sources
	sourceLimit int
}

// NewTerm creates a Term object
func NewTerm() *Term {
	return &Term{
		Error:       make(chan error),
		Done:        make(chan bool),
		source:      DEFAULT_SOURCE,
		sourceLimit: DEFAULT_SOURCE_LIMIT,
	}
}

// Run executes the primary functionality of Term
func (t *Term) Run(ctx context.Context) {
	t.printIntro()
	go func() {
		for {
			t.selectCommand(ctx)
		}
	}()

	select {
	case err := <-t.Error:
		if !strings.Contains(err.Error(), promptui.ErrInterrupt.Error()) {
			t.printError(err)
		}
		t.printExit()
	case <-t.Done:
		t.printExit()
		return
	case <-ctx.Done():
		return
	}
}

func (t *Term) selectCommand(ctx context.Context) {
	pterm.DefaultSection.Println("Main menu")

	var cmdOptions = []string{"Configure", "Get Random Quote", "Get Based On Genres",
		"Get Based On Authors", "Search...", "Exit"}
	prompt := promptui.Select{
		Label: "What would you like?",
		Items: cmdOptions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		t.Error <- fmt.Errorf("prompt failed: %v", err)
		return
	}

	switch result {
	case cmdOptions[0]:
		t.configure()
	case cmdOptions[1]:
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			pterm.DefaultSection.Println("Retrieving Random Quote...")
			t.randomQuote(ctx)
		}()
	case cmdOptions[2]:
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			pterm.DefaultSection.Println("Retrieving Quotes From Genres...")
			t.selectGenre(ctx)
		}()
	case cmdOptions[3]:
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			pterm.DefaultSection.Println("Retrieving Quotes From Authors...")
			t.selectAuthor(ctx)
		}()
	case cmdOptions[4]:
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			pterm.DefaultSection.Println("Searching Quotes...")
			t.selectSearch(ctx, nil)
		}()
	case "Exit":
		t.Done <- true
	}
	t.wg.Wait()
}

func (t *Term) configure() {
	var cmdOptions = []string{"Select source", "Select quotes limit", GO_BACK}
	prompt := promptui.Select{
		Label: "What would you like?",
		Items: cmdOptions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		t.Error <- fmt.Errorf("prompt failed: %v", err)
		return
	}
	switch result {
	case cmdOptions[0]:
		t.configureSelectSource()
	case cmdOptions[1]:
		t.configureSelectSourceLimit()
	case GO_BACK:
		return
	}
}

func (t *Term) configureSelectSource() {
	var cmdOptions = []string{"QuoteGarden", GO_BACK}
	prompt := promptui.Select{
		Label: "Source for quotes",
		Items: cmdOptions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		t.Error <- fmt.Errorf("prompt failed: %v", err)
		return
	}
	switch result {
	case cmdOptions[0]:
		t.source = source.NewQuoteGarden()
	case GO_BACK:
		return
	}
}

func (t *Term) configureSelectSourceLimit() {
	validate := func(input string) error {
		_, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Limit (default=9)",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		t.Error <- fmt.Errorf("prompt failed: %v", err)
		return
	}
	t.sourceLimit, _ = strconv.Atoi(result)
}

func (t *Term) randomQuote(ctx context.Context) {
	quote, err := t.source.RandomQuote(ctx)
	if err != nil {
		t.Error <- fmt.Errorf("could not get random quote from source: %v", err)
		return
	}
	quote.Print()
}

func (t *Term) selectGenre(ctx context.Context) {
	// make HTTP request
	items, err := t.source.AllGenres(ctx)
	if err != nil {
		t.Error <- fmt.Errorf("could not get categories from source: %v", err)
		return
	}
	pterm.Info.Printfln("Number of items: %+v", len(items))

	// filter or not the items
	items = t.selectFilter(items)

	// select from the response genres
	prompt := promptui.Select{
		Label: "Select genre",
		Items: items,
	}
	_, selection, err := prompt.Run()
	if err != nil {
		t.Error <- fmt.Errorf("prompt failed: %v", err)
		return
	}
	// create query options & go further
	qo := &source.QueryOptions{
		Genre: selection,
	}
	t.goFurther(ctx, qo)
}

func (t *Term) selectFilter(items []string) []string {
	var cmdOptions = []string{"No filter", "Filter search", GO_BACK}
	prompt := promptui.Select{
		Label: "Source for quotes",
		Items: cmdOptions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		t.Error <- fmt.Errorf("prompt failed: %v", err)
		return nil
	}
	switch result {
	case cmdOptions[0]:
		return items
	case cmdOptions[1]:
		filter := t.showFilter()
		return t.filter(items, filter)
	case GO_BACK:
	}
	return items
}

func (t *Term) showFilter() string {
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
		Label:     "Filter:",
		Templates: templates,
		Validate:  validate,
	}

	selection, err := prompt.Run()
	if err != nil {
		t.Error <- fmt.Errorf("prompt failed %v", err)
		return ""
	}
	return selection
}

func (t *Term) filter(items []string, filter string) []string {
	var filtered []string
	for _, i := range items {
		if strings.Contains(strings.ToLower(i), strings.ToLower(filter)) {
			filtered = append(filtered, i)
		}
	}
	return filtered
}

func (t *Term) selectAuthor(ctx context.Context) {
	// make HTTP request
	items, err := t.source.AllAuthors(ctx)
	if err != nil {
		t.Error <- fmt.Errorf("could not get categories from source: %v", err)
		return
	}
	pterm.Info.Printfln("Number of items: %+v", len(items))

	// filter or not the items
	items = t.selectFilter(items)

	// select from the response genres
	prompt := promptui.Select{
		Label: "Select author",
		Items: items,
	}

	_, selection, err := prompt.Run()
	if err != nil {
		t.Error <- fmt.Errorf("prompt failed %v", err)
		return
	}
	// create query options & go further
	qo := &source.QueryOptions{
		Author: selection,
	}
	t.goFurther(ctx, qo)
}

func (t *Term) selectSearch(ctx context.Context, qo *source.QueryOptions) {
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
		t.Error <- fmt.Errorf("prompt failed %v", err)
		return
	}

	// create query options & go further
	qoTemp := &source.QueryOptions{
		Query: selection,
	}
	if qo != nil {
		qoTemp.Author = qo.Author
		qoTemp.Genre = qo.Genre
	}
	t.goFurther(ctx, qoTemp)
}

func (t *Term) goFurther(ctx context.Context, qo *source.QueryOptions) {
	for {
		t.printSection(qo)
		itemSelection := []string{"Show all quotes", "Get random quote", "Add a filter...", GO_BACK}
		prompt := promptui.Select{
			Label: "What would you like?",
			Items: itemSelection,
		}

		_, result, err := prompt.Run()
		if err != nil {
			t.Error <- fmt.Errorf("prompt failed: %v", err)
			return
		}

		switch result {
		case itemSelection[0]:
			t.showAllQuotes(ctx, qo)
		case itemSelection[1]:
			t.showRandomQuote(ctx, qo)
		case itemSelection[2]:
			t.selectSearch(ctx, qo)
		case GO_BACK:
			return
		}
	}
}

func (t *Term) showAllQuotes(ctx context.Context, qo *source.QueryOptions) {
	pageSelection := 1

	for {
		// query based on selection
		opt := &source.QueryOptions{
			Limit:  int32(t.sourceLimit),
			Page:   int32(pageSelection),
			Genre:  qo.Genre,
			Author: qo.Author,
			Query:  qo.Query,
		}
		quotes, pag, err := t.source.Quotes(ctx, opt)
		if err != nil {
			t.Error <- fmt.Errorf("quotes query failed: %v", err)
			return
		}
		title := fmt.Sprintf("PAGE %d/%d", pageSelection, pag.TotalPages)
		t.source.PrintQuotesPage(title, quotes)

		itemSelection := []string{"Next Page", "Previous Page", GO_BACK}
		prompt := promptui.Select{
			Label: "Select action",
			Items: itemSelection,
		}

		_, result, err := prompt.Run()
		if err != nil {
			t.Error <- fmt.Errorf("prompt failed: %v", err)
			return
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
		case GO_BACK:
			return
		}
	}
}

func (t *Term) showRandomQuote(ctx context.Context, qo *source.QueryOptions) {
	pageSelection := 1
	for {
		// query based on selection
		opt := &source.QueryOptions{
			Limit:  int32(t.sourceLimit),
			Page:   int32(pageSelection),
			Genre:  qo.Genre,
			Author: qo.Author,
			Query:  qo.Query,
		}
		quotes, pag, err := t.source.Quotes(ctx, opt)
		if err != nil {
			t.Error <- fmt.Errorf("quotes query failed: %v", err)
			return
		}
		// get random quote
		rand.Seed(time.Now().Unix())
		randQ := quotes[rand.Intn(len(quotes))]
		randQ.Print()
		// change page
		pageSelection = rand.Intn(pag.TotalPages) + 1

		// after menu
		itemSelection := []string{"Get Another", GO_BACK}
		prompt := promptui.Select{
			Label: "Random quote",
			Items: itemSelection,
		}

		_, result, err := prompt.Run()
		if err != nil {
			t.Error <- fmt.Errorf("prompt failed: %v", err)
			return
		}
		switch result {
		case itemSelection[0]:
			// just continue with another request
		case GO_BACK:
			return
		}
	}
}

func (t *Term) printIntro() {
	pterm.Println()
	newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgBlack),
		BackgroundStyle: pterm.NewStyle(pterm.BgGreen),
		Margin:          20,
	}
	newHeader.Println("Yo! I'm GoQu!")
}

func (t *Term) printError(err error) {
	pterm.Println()
	newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgWhite),
		BackgroundStyle: pterm.NewStyle(pterm.BgRed),
		Margin:          10,
	}
	newHeader.Println("ERROR:", err)
}

func (t *Term) printExit() {
	pterm.Println()
	newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgWhite),
		BackgroundStyle: pterm.NewStyle(pterm.BgBlue),
		Margin:          20,
	}
	rand.Seed(time.Now().Unix())
	para := pterm.DefaultParagraph.WithMaxWidth(60).Sprintln(assets.QUOTES_GOKU[rand.Intn(len(assets.QUOTES_GOKU))])
	newHeader.Printfln("%s\n-- Goku, 'Dragon Ball Z'", para)
}

func (t *Term) printSection(qo *source.QueryOptions) {
	var selection []string
	options := []string{qo.Genre, qo.Author, qo.Query}
	for _, o := range options {
		if strings.TrimSpace(o) != "" {
			selection = append(selection, o)
		}
	}
	pterm.DefaultSection.Printfln("Selected options: %s", strings.Join(selection, " & "))
}
