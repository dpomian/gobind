package models

type Notebook struct {
	Id      string
	Title   string `json:"title"`
	Content string `json:"content"`

	// user_id
}
