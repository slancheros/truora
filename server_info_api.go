package main

import (
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"truora/handlers"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/serverInfo", func(r chi.Router) {
		r.Get("/{domain:}", handlers.RetrieveDomainInfo)
		r.Get("/", handlers.RetrieveDomainInfo)
		r.Get("/list", handlers.ListDomainsQueried)
	})

	http.ListenAndServe(":3344", r)
}
