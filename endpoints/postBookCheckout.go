package endpoints

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dgaldamez77/loc/dal"
	"github.com/dgaldamez77/loc/dal/dto"
	"github.com/dgaldamez77/loc/endpoints/models"
	"github.com/dgaldamez77/loc/log"
)

func postBookCheckout(w http.ResponseWriter, r *http.Request) {
	db := dal.NewDAL()

	var bookCheckout models.PostBookCheckout
	if err := json.NewDecoder(r.Body).Decode(&bookCheckout); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify book availability
	q := []dal.QueryParams{
		{
			FieldName: "bi.book_id",
			Value:     bookCheckout.BookID,
		},
	}
	bookInventory, err := db.GetBookInventory(q)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(w, "book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if bookInventory.BookCount <= 0 {
		http.Error(w, "book is not available in inventory", http.StatusInternalServerError)
		return
	}

	q = []dal.QueryParams{
		{
			FieldName: "bco.book_id",
			Value:     bookCheckout.BookID,
		},
		{
			FieldName: "bco.user_id",
			Value:     bookCheckout.UserID,
		},
	}
	var checkedout *dto.BookCheckout
	bookCheckouts, err := db.GetBookCheckouts(q)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			// not needed but only used to indicate that there is
			// no checkout when no data is found
			checkedout = nil
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	for i := range bookCheckouts {
		bco := bookCheckouts[i]
		if bco.Returned == nil {
			checkedout = &bco
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if checkedout != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("book is already checked out to this user"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	tx, err := db.BeginTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var outCheckout *dto.BookCheckout
	newCheckout := dto.BookCheckout{
		BookID:    bookCheckout.BookID,
		UserID:    bookCheckout.UserID,
		CreatedBy: bookCheckout.CreatedBy,
	}
	if outCheckout, err = db.InsertBookCheckout(newCheckout); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		if err = db.RollbackTransaction(tx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db.RollbackTransaction(tx)
		return
	}

	// reduce count of book availability
	changes := map[string]interface{}{
		"book_count": bookInventory.BookCount - 1,
		"updated_by": bookCheckout.CreatedBy,
	}
	if _, err = db.UpdateBookInventory(bookInventory.BookID, changes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		db.RollbackTransaction(tx)
		return
	}

	if err = db.Commit(tx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		db.RollbackTransaction(tx)
		return
	}

	b, err := json.Marshal(outCheckout)
	if err != nil {
		log.LogError(err, "failed to marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}
