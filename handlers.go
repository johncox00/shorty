package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"math/rand"
)

type LinkShortnerAPI struct {
	myconnection *MongoConnection
}

type UrlMap struct {
  ShortUrl   string  `json:shorturl`
  LongUrl    string  `json:longurl`
}

type UrlClickCount struct {
	ClickCount  int `json:click_count`
}

type ApiErrorResponse struct {
  ErrorMsg  string  `json:error_message`
}

func NewUrlLinkShortenerAPI() *LinkShortnerAPI {
	LS := &LinkShortnerAPI{
		myconnection: NewDBConnection(),
	}
	return LS
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (LS *LinkShortnerAPI) GenerateShortCode() string {
	n := 8
	b := make([]byte, n)
  for i := range b {
      b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
  }
  return string(b)
}

func (Ls *LinkShortnerAPI) UrlRoot(w http.ResponseWriter, r *http.Request) {
  fmt.Print(w, "Hello and welcome to the Go link shortner API \n"+
    	"Do a Get request with the short Link to get the long Link \n"+
		"Do a POST request with long Link to get a short Link \n")
}

func (Ls *LinkShortnerAPI)UrlShow(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sUrl := vars["shorturl"]
	responseEncoder := json.NewEncoder(w)
  if len(sUrl) > 0 {
    fmt.Print("retrieving long url...")
		lUrl, err := Ls.myconnection.FindUrl(sUrl)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
	    responseEncoder.Encode(&ApiErrorResponse{"No long URL found."})
	    return
		}
		responseEncoder.Encode(&UrlMap{ShortUrl: sUrl, LongUrl: lUrl})
  }
}

func (Ls *LinkShortnerAPI)UrlClickShow(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sUrl := vars["shorturl"]
	responseEncoder := json.NewEncoder(w)
  if len(sUrl) > 0 {
    fmt.Print("retrieving click count...")
		count, err := Ls.myconnection.FindCount(sUrl)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
	    responseEncoder.Encode(&ApiErrorResponse{"No long URL found."})
	    return
		}
		responseEncoder.Encode(&UrlClickCount{count})
  }
}

func (Ls *LinkShortnerAPI)UrlCreate(w http.ResponseWriter, r *http.Request) {
  reqBody := new(UrlMap)
  responseEncoder := json.NewEncoder(w)
  if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
    w.WriteHeader(http.StatusBadRequest)
    responseEncoder.Encode(&ApiErrorResponse{"Invalid JSON body."})
    return
  }
  fmt.Print("adding url to the database")
	sCode := Ls.GenerateShortCode()
	_, err := Ls.myconnection.InsertUrl(reqBody.LongUrl, sCode)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		if err := responseEncoder.Encode(&ApiErrorResponse{err.Error()}); err != nil {
			fmt.Fprintf(w, "Error %s occured while trying to add the url \n", err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	responseEncoder.Encode(&UrlMap{ShortUrl: sCode, LongUrl: reqBody.LongUrl})
}
