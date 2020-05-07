package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)

	}

	return router
}

var routes = Routes{
	Route{
		"Healthy",
		"GET",
		"/api/formula/healthy",
		Healthy,
	},
	Route{
		"FunctionList",
		"GET",
		"/api/formula/formulas",
		FunctionList,
	},
	Route{
		"FunctionShow",
		"GET",
		"/api/formula/formulas/{id}",
		FunctionShow,
	},
	Route{
		"FunctionAdd",
		"POST",
		"/api/formula/formulas",
		FunctionAdd,
	},
	Route{
		"FunctionUpdate",
		"PUT",
		"/api/formula/formulas/{id}",
		FunctionUpdate,
	},
	Route{
		"FunctionRemove",
		"DELETE",
		"/api/formula/formulas/{id}",
		FunctionRemove,
	},
	Route{
		"FunctionRemoveMasivo",
		"DELETE",
		"/api/formula/formulas",
		FunctionRemoveMasivo,
	},
	Route{
		"FunctionExecute",
		"POST",
		"/api/formula/execute",
		FunctionExecute,
	},
	Route{
		"FunctionAddPublic",
		"POST",
		"/api/formula/public/formulas",
		FunctionAddPublic,
	},
}
