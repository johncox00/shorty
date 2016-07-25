package main

import "net/http"

func main() {
  LinkShortener := NewUrlLinkShortenerAPI()
  routes := CreateRoutes(LinkShortener)
  router := NewLinkShortenerRouter(routes)
  http.ListenAndServe(":5100", router)
}
