package models

type PostBookCheckout struct {
	BookID int `json:"book_id"`
	UserID int `json:"user_id"`
}
