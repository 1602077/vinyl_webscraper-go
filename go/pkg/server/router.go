// server packages api routing and handling for go webscraping app.
package server

import (
	"github.com/gorilla/mux"
)

// NewRouter create a gorilla mux Router using the route defined by the variable
// routes in routes.go.
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	/* Serve css files
	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("../static/")))
		)
	*/

	return router
}
