package endpoints

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dgaldamez77/loc/dal"
	"github.com/dgaldamez77/loc/log"
	"github.com/go-chi/chi"
)

func getBook(w http.ResponseWriter, r *http.Request) {
	db := dal.NewDAL()

	bookID := chi.URLParam(r, "bookID")

	q := []dal.QueryParams{
		{
			FieldName: "b.book_id",
			Value:     bookID,
		},
	}
	book, err := db.GetBook(q)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(w, "book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	b, err := json.Marshal(book)
	if err != nil {
		log.LogError(err, "failed to marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
