package dto

import "time"

type BookCheckout struct {
	ID        int        `json:"book_checkout_id"`
	BookID    int        `json:"book_id"`
	UserID    int        `json:"user_id"`
	Returned  *time.Time `json:"returned"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy string     `json:"updated_by,omitempty"`
	Created   time.Time  `json:"created"`
	Updated   *time.Time `json:"updated,omitempty"`
}
