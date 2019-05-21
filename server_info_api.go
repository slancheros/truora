package main

import (
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type DomainInfo struct {
	Servers          []ServerDesc `json:"servers"`
	ServersChanged   string       `json:"servers_changed"`
	ServersSSLGrade  string       `json:"ssl_grade"`
	PreviousSSLGrade string       `json:"previous_ssl_grade"`
	Logo             string       `json:"logo"`
	Title            string       `json:"title"`
	IsDown           bool         `json:"is_down"`
}

type ServerDesc struct {
	ServerAddress string `json:"address"`
	SSLGrade      string `json:"ssl-grade"`
	Country       string `json:"country"`
	Owner         string `json:"owner"`
}

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/serverInfo", func(r chi.Router) {
		r.Get("/{domain:}", retrieveDomainInfo)
		r.Get("/", retrieveDomainInfo)
	})

	http.ListenAndServe(":3334", r)
}
