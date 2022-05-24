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
		HomePage,
	},
	Route{
		"Refresh",
		"GET",
		"/refresh",
		GetRecordPrices,
	},
	Route{
		"RecordList",
		"GET",
		"/Record/{id}",
		GetRecord,
	},
}
