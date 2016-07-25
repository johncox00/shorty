package main

import (
	"github.com/gorilla/mux"
)

func NewLinkShortenerRouter(routes Routes) *mux.Router {
  router := mux.NewRouter().StrictSlash(true)
  for _, route := range routes {
    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(route.Handler)
  }
  return router
}