package main

import (
	"flag"
	"fmt"
	"github.com/isichei/recipe-book/internal/filesyncer"
	"github.com/isichei/recipe-book/internal/database"
	"log"
	"log/slog"
	"net"
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

func startTCPFileServer(port, apiKey, directory string) error {
	// TODO: Could probably use the DB here if I end up using it as a cache for search
	fc, err := filesyncer.CreateFileCache(directory)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", port, err)
	}
	defer ln.Close()

	slog.Info("TCP Listening for authenticated connection", "port", port)

	AcceptConnErrCounter := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Warn("Failed to accept connection", "error", err)
			if AcceptConnErrCounter >= 5 {
				return err
			}
			AcceptConnErrCounter += 1
			continue
		}
		conn, err = filesyncer.AuthenticateListenerConnection(conn, apiKey)
		if err != nil {
			conn.Close()
			continue
		}
		// TODO: Think I need to reset the fc after this call to refresh it
		// I also might want to use a mutex here incase two syncs are called
		// but the later will be me doing it twice so not that much of a problem right

		// Now do the real work with the authenticated connection
		syncer := filesyncer.Syncer{Replica: true, Conn: conn, FileCache: fc}

		if err := syncer.Run(); err != nil {
			log.Printf("Replica failed error: %s\n", err)
		}
		log.Println("Replica sync complete")
	}
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
			if err := startTCPFileServer("9000", apiKey, path.Join(*staticPath, "recipe_mds")); err != nil {
				log.Printf("TCP File Server failed: %v", err)
			}
		}()
	}

	app.logger.Info(fmt.Sprintf("Starting Recipe App on %s...", *port))
	http.ListenAndServe(":"+*port, app.routes())
}
