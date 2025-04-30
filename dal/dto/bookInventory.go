package dto

import "time"

type BookInventory struct {
	BookID    int        `json:"book_id"`
	BookCount int        `json:"book_count"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy string     `json:"updated_by,omitempty"`
	Created   time.Time  `json:"created"`
	Updated   *time.Time `json:"updated,omitempty"`
}
