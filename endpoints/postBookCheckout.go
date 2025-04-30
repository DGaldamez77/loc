package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/dgaldamez77/oloc/dal"
	"github.com/dgaldamez77/oloc/log"
)

func postBookCheckout(w http.ResponseWriter, r *http.Request) {
	db := dal.NewDAL()

	q := []dal.QueryParams{
		{
			FieldName: "b.active",
			Value:     true,
		},
	}
	books, err := db.GetBooks(q)
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
