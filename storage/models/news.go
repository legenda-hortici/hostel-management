package models

type News struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Annotation string `json:"annotation"`
	Text       string `json:"text"`
	Date       string `json:"date"`
	NewsType   string `json:"newsType"`
}

type Notice struct {
	ID         int
	Title      string
	Annotation string
	Text       string
	Date       string
	NewsType   string
}
