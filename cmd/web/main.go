package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/isichei/recipe-book/internal/database"
)

type application struct {
	db database.RecipeDatabase
}

func main() {

	listenAddr := flag.String("listenaddr", ":8000", "The address for the API to listen on")
	flag.Parse()

	app := &application{database.NewTestDatabase()}

	log.Printf("Starting Recipe App on %s...\n", *listenAddr)
	http.ListenAndServe(*listenAddr, app.routes())
}
