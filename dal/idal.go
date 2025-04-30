package dal

import (
	"database/sql"

	"github.com/dgaldamez77/loc/dal/dto"
)

type IDAL interface {
	query(string, ...interface{}) (*sql.Rows, error)
	queryRow(string, ...interface{}) *sql.Row

	BeginTransaction() (*sql.Tx, error)
	Commit(*sql.Tx) error
	RollbackTransaction(*sql.Tx) error

	GetBooks([]QueryParams) ([]dto.Book, error)
	GetBooksUsingAuthor([]QueryParams) ([]dto.Book, error)
	GetBook([]QueryParams) (*dto.Book, error)

	GetBookCheckouts([]QueryParams) ([]dto.BookCheckout, error)
	InsertBookCheckout(dto.BookCheckout) (*dto.BookCheckout, error)

	GetBookInventory([]QueryParams) (*dto.BookInventory, error)
	UpdateBookInventory(int, map[string]interface{}) (*dto.BookInventory, error)
}
