package endpoints

import (
	"fmt"
	"net/http"

	"github.com/dgaldamez77/loc/log"
	"github.com/go-chi/chi"
)

var (
	r = chi.NewRouter()
)

func RegisterEndpoints() chi.Router {
	registerEndpoing(http.MethodGet, "/books", getBooks)
	registerEndpoing(http.MethodGet, "/book/{bookID}", getBook)
	registerEndpoing(http.MethodPost, "/book/checkout", postBookCheckout)
	registerEndpoing(http.MethodPost, "/book/checkin", postBookCheckin)

	return r
}

func registerEndpoing(method, path string, handler http.HandlerFunc) {
	r.MethodFunc(method, path, handler)
	log.LogInfo(fmt.Sprintf("Registering endpoint %s %s", method, path))
}
