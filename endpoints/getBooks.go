package endpoints

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgaldamez77/oloc/dal"
	"github.com/dgaldamez77/oloc/dal/dto"
	"github.com/dgaldamez77/oloc/log"
)

func getBooks(w http.ResponseWriter, r *http.Request) {
	db := dal.NewDAL()

	title := r.URL.Query().Get("title")
	publisher := r.URL.Query().Get("publisher")
	genre := r.URL.Query().Get("genre")
	author := r.URL.Query().Get("author")
	isbn := r.URL.Query().Get("isbn")

	useByAuthor := false
	q := []dal.QueryParams{
		{
			FieldName: "b.active",
			Value:     true,
		},
	}

	if title != "" {
		q = append(q, dal.QueryParams{
			FieldName:  "lower(b.title)",
			Value:      "%" + strings.ToLower(title) + "%",
			Comparison: "LIKE",
		})
	}
	if isbn != "" {
		q = append(q, dal.QueryParams{
			FieldName: "b.isbn",
			Value:     isbn,
		})
	}
	if genre != "" {
		q = append(q, dal.QueryParams{
			FieldName:  "lower(g.genre)",
			Value:      "%" + strings.ToLower(genre) + "%",
			Comparison: "LIKE",
		})
	}
	if publisher != "" {
		q = append(q, dal.QueryParams{
			FieldName:  "lower(p.publisher)",
			Value:      "%" + strings.ToLower(publisher) + "%",
			Comparison: "LIKE",
		})
	}
	if author != "" {
		q = append(q, dal.QueryParams{
			FieldName:  "lower(a.first_name || ' ' || a.last_name)",
			Value:      "%" + strings.ToLower(author) + "%",
			Comparison: "LIKE",
		})
		useByAuthor = true
	}

	var err error
	var books []dto.Book
	if useByAuthor {
		books, err = db.GetBooksUsingAuthor(q)
	} else {
		books, err = db.GetBooks(q)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	b, err := json.Marshal(books)
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
