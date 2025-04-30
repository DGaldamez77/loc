package dto

import "time"

type Book struct {
	ID          int        `json:"book_id"`
	Title       string     `json:"title"`
	PublisherID int        `json:"publisher_id"`
	Publisher   Publisher  `json:"publisher"`
	ISBN        string     `json:"isbn"`
	Synopsis    string     `json:"synopsis"`
	GenreID     int        `json:"genre_id"`
	Genre       Genre      `json:"genre"`
	Authors     []Author   `json:"author"`
	Active      bool       `json:"active"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by,omitempty"`
	Created     time.Time  `json:"created"`
	Updated     *time.Time `json:"updated,omitempty"`
}
