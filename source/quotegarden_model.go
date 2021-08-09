package source

type QGGenres struct {
	StatusCode  int        `json:"statusCode"`
	Message     string     `json:"message"`
	Pagination  Pagination `json:"pagination"`
	TotalQuotes int        `json:"totalQuotes"`
	Data        []string   `json:"data"`
}

type QGAuthors struct {
	StatusCode  int        `json:"statusCode"`
	Message     string     `json:"message"`
	Pagination  Pagination `json:"pagination"`
	TotalQuotes int        `json:"totalQuotes"`
	Data        []string   `json:"data"`
}

type QGQuote struct {
	StatusCode  int        `json:"statusCode"`
	Message     string     `json:"message"`
	Pagination  Pagination `json:"pagination"`
	TotalQuotes int        `json:"totalQuotes"`
	Data        []*Quote   `json:"data"`
}

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	NextPage    int `json:"nextPage"`
	TotalPage   int `json:"totalPages"`
}

// TODO: quote needs to be generic
type Quote struct {
	ID          string `json:"_id"`
	QuoteText   string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
	QuoteGenre  string `json:"quoteGenre"`
}
