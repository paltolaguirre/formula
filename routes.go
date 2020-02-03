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
		"FormulaList",
		"GET",
		"/api/formula/formulas",
		FormulaList,
	},
	Route{
		"FormulaShow",
		"GET",
		"/api/formula/formulas/{id}",
		FormulaShow,
	},
	Route{
		"FormulaAdd",
		"POST",
		"/api/formula/formulas",
		FormulaAdd,
	},
	Route{
		"FormulaUpdate",
		"PUT",
		"/api/formula/formulas/{id}",
		FormulaUpdate,
	},
	Route{
		"FormulaRemove",
		"DELETE",
		"/api/formula/formulas/{id}",
		FormulaRemove,
	},
	Route{
		"FormulaRemoveMasivo",
		"DELETE",
		"/api/formula/formulas",
		FormulaRemoveMasivo,
	},
}
