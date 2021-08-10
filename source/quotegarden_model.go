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

type QGQuotes struct {
	StatusCode  int        `json:"statusCode"`
	Message     string     `json:"message"`
	Pagination  Pagination `json:"pagination"`
	TotalQuotes int        `json:"totalQuotes"`
	Data        []*QGQuote `json:"data"`
}

func (qgq *QGQuotes) DataToQuotes() (retQ []*Quote) {
	for _, q := range qgq.Data {
		retQ = append(retQ, q.ToQuote())
	}
	return retQ
}

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	NextPage    int `json:"nextPage"`
	TotalPages  int `json:"totalPages"`
}

type QGQuote struct {
	ID          string `json:"_id"`
	QuoteText   string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
	QuoteGenre  string `json:"quoteGenre"`
}

func (qgq *QGQuote) ToQuote() *Quote {
	return &Quote{
		ID:     qgq.ID,
		Author: qgq.QuoteAuthor,
		Text:   qgq.QuoteText,
		Genre:  qgq.QuoteGenre,
	}
}
