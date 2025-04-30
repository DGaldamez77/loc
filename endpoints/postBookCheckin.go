package endpoints

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgaldamez77/loc/dal"
	"github.com/dgaldamez77/loc/dal/dto"
	"github.com/dgaldamez77/loc/endpoints/models"
	"github.com/dgaldamez77/loc/log"
)

func postBookCheckin(w http.ResponseWriter, r *http.Request) {
	db := dal.NewDAL()

	var bookCheckout models.PostBookCheckin
	if err := json.NewDecoder(r.Body).Decode(&bookCheckout); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// verify book is checked out to user
	q := []dal.QueryParams{
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

	// find the book checkout that has not been returned
	for i := range bookCheckouts {
		bco := bookCheckouts[i]
		if bco.Returned == nil {
			checkedout = &bco
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")

	// no checkout found... respond with error
	if checkedout == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("book is not checked out to this user"))
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

	// Update book checkout table to create the association between book and user
	var outCheckout *dto.BookCheckout
	changes := map[string]interface{}{
		"returned":   time.Now().UTC(),
		"updated_by": bookCheckout.CreatedBy,
	}
	if outCheckout, err = db.UpdateBookCheckout(checkedout.ID, changes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		if err = db.RollbackTransaction(tx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_ = db.RollbackTransaction(tx)
		return
	}

	// get current book availability
	q = []dal.QueryParams{
		{
			FieldName: "bi.book_id",
			Value:     bookCheckout.BookID,
		},
	}
	bookInventory, err := db.GetBookInventory(q)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(w, "no book availability found in book inventory", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_ = db.RollbackTransaction(tx)
		return
	}

	// reduce count of book availability
	changes = map[string]interface{}{
		"book_count": bookInventory.BookCount + 1,
		"updated_by": bookCheckout.CreatedBy,
	}
	if _, err = db.UpdateBookInventory(bookInventory.BookID, changes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		_ = db.RollbackTransaction(tx)
		return
	}

	if err = db.Commit(tx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		_ = db.RollbackTransaction(tx)
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
