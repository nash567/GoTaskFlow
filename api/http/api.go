package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Q map[string]string
type handlerFunc func(http.ResponseWriter, *http.Request)

func RegisterHandlers(
	router *mux.Router,
	method string,
	path string,
	query Q,
	handler handlerFunc,
	middlewares ...mux.MiddlewareFunc,
) *mux.Route {
	// new route for the path.
	route := router.Path(path)

	// specify the method unless it's an empty string or "*".
	if method != "" && method != "*" {
		route = route.Methods(method)
	}
	// specify the query params (if any)
	for key, value := range query {
		route.Queries(key, value)
	}
	// specify middlewares (if any) on a dedicated sub-Router
	// route type do not have use method so converting it to type router
	if len(middlewares) > 0 {
		router = route.Subrouter()
		router.Use(middlewares...)
		route = router.NewRoute()
	}
	// setup the handler which gets us a new route
	if handler != nil {
		route.HandlerFunc(handler)
	}

	return route

}
