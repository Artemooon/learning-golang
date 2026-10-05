package main

import (
	"log"
	"net/http"

	"github.com/artem/vintage-jazz-records-api/internal/httpapi"
	"github.com/artem/vintage-jazz-records-api/internal/records"
)

func main() {
	store := records.NewStore([]records.Record{
		{ID: "kind-of-blue", Title: "Kind of Blue", Artist: "Miles Davis", Year: 1959, Label: "Columbia"},
		{ID: "time-out", Title: "Time Out", Artist: "The Dave Brubeck Quartet", Year: 1959, Label: "Columbia"},
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpapi.NewHandler(store),
	}

	log.Printf("Vintage jazz records API listening on http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
