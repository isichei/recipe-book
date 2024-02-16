package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/isichei/recipe-book/internal/database"
)

type application struct {
	db     database.RecipeDatabase
	logger *slog.Logger
}

func main() {
	listenAddr := flag.String("listenaddr", ":8000", "The address for the API to listen on")
	flag.Parse()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &application{database.NewTestDatabaseFromDir("./static/recipe_mds/"), logger}

	app.logger.Info(fmt.Sprintf("Starting Recipe App on %s...", *listenAddr))
	http.ListenAndServe(*listenAddr, app.routes())
}
