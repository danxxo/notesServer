package dto

type Note struct {
	Id              string `json:"id"`
	AuthorFirstName string `json:"authorFirstName"`
	AuthorLastName  string `json:"authorLastName"`
	Note            string `json:"note"`
}

type Notes []Note
