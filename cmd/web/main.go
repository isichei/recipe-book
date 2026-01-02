package main

import (
	"flag"
	"fmt"
	"github.com/isichei/recipe-book/internal/database"
	"github.com/isichei/recipe-book/internal/filesyncer"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
)

type creds struct {
	username string
	password string
}

type application struct {
	creds        creds
	db           database.RecipeDatabase
	staticFolder string
	logger       *slog.Logger
	enableWrite  bool
}

func main() {
	defaultPort, portExists := os.LookupEnv("PORT")
	if !portExists {
		defaultPort = "8000"
	}

	defaultRecipeDir, defaultDirExists := os.LookupEnv("RECIPE_FILES")
	if !defaultDirExists {
		defaultRecipeDir = "./static/recipe_mds/"
	}

	user, userExists := os.LookupEnv("RECIPE_USER")
	if !userExists {
		user = "user"
	}
	password, passwordExists := os.LookupEnv("RECIPE_PASSWORD")
	if !passwordExists {
		password = "password"
	}

	apiKey, apiKeyExists := os.LookupEnv("TCP_API_KEY")

	port := flag.String("port", defaultPort, fmt.Sprintf("The address for the API to listen on. (Default %s)", defaultPort))
	recipeDir := flag.String("recipe-dir", defaultRecipeDir, fmt.Sprintf("Path to the recipe files (if a directory then expects it to contain markdown files for each recipe. If a filepath expects a database file of recipes. (Default %s)", defaultRecipeDir))
	dbPath := flag.String("db", "", "Path to the recipe db file")
	staticPath := flag.String("static-path", "/static/", "Path to the static asset folder")
	enableWrite := flag.Bool("enable-write", false, "Enable the /add-recipe path in the app")
	enableFileSync := flag.Bool("enable-filesync", false, "Enable the ability to sync files over TCP")
	flag.Parse()

	if *enableFileSync && !apiKeyExists {
		panic("TCP_API_KEY not set")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var db database.RecipeDatabase
	var err error

	if *dbPath == "" && *recipeDir == "" {
		log.Fatal("Both db and recipe-dir were set to empty strings need to specify one")
	} else if *dbPath != "" {
		db, err = database.NewSqlDatabase(*dbPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db = database.NewTestDatabaseFromDir(*recipeDir)
	}
	appCreds := creds{user, password}
	app := &application{appCreds, db, *staticPath, logger, *enableWrite}

	if *enableFileSync {
		app.logger.Info("Starting TCP File Sync on 9000")
		go func() {
			if err := filesyncer.StartReplicaTCPFileServer("0.0.0.0:9000", apiKey, path.Join(*staticPath, "recipe_mds")); err != nil {
				log.Printf("TCP File Server failed: %v", err)
			}
		}()
	}

	app.logger.Info(fmt.Sprintf("Starting Recipe App on %s...", *port))
	http.ListenAndServe(":"+*port, app.routes())
}
