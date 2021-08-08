package source

type QGGenres struct {
	StatusCode  int                `json:"statusCode"`
	Message     string             `json:"message"`
	Pagination  QGGenresPagination `json:"pagination"`
	TotalQuotes int                `json:"totalQuotes"`
	Data        []string           `json:"data"`
}

type QGGenresPagination struct {
	CurrentPage int `json:"currentPage"`
	NextPage    int `json:"nextPage"`
	TotalPage   int `json:"totalPages"`
}
