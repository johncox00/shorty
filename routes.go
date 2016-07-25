package main

import "net/http"

type Route struct {
  Name    string
  Method  string
  Pattern string
  Handler http.HandlerFunc
}

type Routes []Route

func CreateRoutes(LS *LinkShortnerAPI) Routes{
  return Routes{
    Route{
      "Root",
      "GET",
      "/",
      LS.UrlRoot,
    },
    Route{
      "UrlShow",
      "GET",
      "/{shorturl}",
      LS.UrlShow,
    },
    Route{
      "UrlCreate",
      "POST",
      "/shorturls",
      LS.UrlCreate,
    },
    Route{
      "UrlClickShow",
      "GET",
      "/{shorturl}/click_count",
      LS.UrlClickShow,
    },
  }
}
