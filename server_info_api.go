package main

import (

	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	 "encoding/json"
	"fmt"
)


type DomainInfo struct {
	Key []ServerDesc `json:"servers"`
	ServersChanged string `json:"servers_changed"`
	ServersSSLGrade string `json:"ssl_grade"`
	PreviousSSLGrade string `json:"previous_ssl_grade"`
	Logo  string `json:"logo"`
	Title string `json:"title"`
	IsDown bool `json:"is_down"`
}

type ServerDesc struct{
	ServerAddress string `json:"address"`
	SSLGrade      string `json:"ssl-grade"`
	Country       string `json:"country"`
	Owner         string `json:"owner"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("hello world"))
		respondwithJSON( w, 200,map[string]string{"message": "hello Sandra"})
	})
	http.ListenAndServe(":3334", r)
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

