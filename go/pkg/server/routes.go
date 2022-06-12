// server packages api routing and handling for go webscraping app.
package server

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"HomePage",
		"GET",
		"/",
		GetRecords,
	},
	Route{
		"Refresh",
		"GET",
		"/refresh",
		PutRecords,
	},
	Route{
		"GetRecord",
		"GET",
		"/Record/{id}",
		GetRecord,
	},
}
