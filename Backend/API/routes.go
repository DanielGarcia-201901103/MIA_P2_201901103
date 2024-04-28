package API

import (
	"net/http"
	"github.com/gorilla/mux"
)

type Route struct {
	Name       string
	MethodType string
	Path       string
	Handler    http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Name(route.Name).
			Methods(route.MethodType).
			Path(route.Path).
			Handler(route.Handler)
	}
	return router
}

var routes = Routes{
	Route{
		"Home",
		"GET",
		"/daniel",
		datos,
	},
	{
		"dataFront",
		"POST",
		"/sendContent",
		_dataFront,
	},
	{
		"loginFront",
		"POST",
		"/login",
		doLogin,
	},
	{
		"loogutFront",
		"GET",
		"/logout",
		doLogout,
	},
	{
		"loginFront2",
		"POST",
		"/loginFront",
		doLoginFront,
	},
}