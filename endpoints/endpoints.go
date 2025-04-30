package endpoints

import (
	"fmt"
	"net/http"

	"github.com/dgaldamez77/oloc/log"
	"github.com/go-chi/chi"
)

var (
	r = chi.NewRouter()
)

func RegisterEndpoints() chi.Router {
	registerEndpoing(http.MethodGet, "/books", getBooks)
	registerEndpoing(http.MethodGet, "/book", getBook)
	registerEndpoing(http.MethodPost, "/book/checkout", postBookCheckout)
	registerEndpoing(http.MethodPost, "/books/checkin", postBookCheckin)

	return r
}

func registerEndpoing(method, path string, handler http.HandlerFunc) {
	r.MethodFunc(method, path, handler)
	log.LogInfo(fmt.Sprintf("Registering endpoint %s %s", method, path))
}
