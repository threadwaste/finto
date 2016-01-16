package finto

import (
	"net/http"

	"github.com/gorilla/mux"
)

type fintoHandlerFunc func(fc *fintoContext) http.Handler

type Route struct {
	Handler fintoHandlerFunc
	Method  string
	Name    string
	Pattern string
}

type Routes []Route

var routes = Routes{
	Route{
		Handler: rolesList,
		Name:    "list-role",
		Method:  "GET",
		Pattern: "/roles",
	},
	Route{
		Handler: rolesSetActive,
		Name:    "set-active-role",
		Method:  "PUT",
		Pattern: "/roles",
	},
	Route{
		Handler: rolesShow,
		Name:    "show-role",
		Method:  "GET",
		Pattern: "/roles/{alias}",
	},
	Route{
		Handler: mockProfileCreds,
		Name:    "get-role-credentials",
		Method:  "GET",
		Pattern: "/roles/{alias}/credentials",
	},
	Route{
		Handler: mockProfile,
		Name:    "metadata-iam-secreds",
		Method:  "GET",
		Pattern: "/latest/meta-data/iam/security-credentials/",
	},
	Route{
		Handler: mockProfileCreds,
		Name:    "metadata-iam-secreds-role",
		Method:  "GET",
		Pattern: "/latest/meta-data/iam/security-credentials/{alias}",
	},
}

func FintoRouter(fc *fintoContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Name(route.Name).
			Path(route.Pattern).
			Handler(route.Handler(fc))
	}

	return router
}
