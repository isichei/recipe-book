package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/isichei/recipe-book/internal/database"
	"github.com/isichei/recipe-book/internal/filesyncer"
)

// Because remembering when to use single or double dash args is not worth it
func acceptDoubleDashArgs(subArgs []string) []string {
	for i, arg := range subArgs {
		if arg[:2] == "--" {
			subArgs[i] = "-" + arg[2:]
		}
	}
	return subArgs
}

func main() {
	// Set log level to debug
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// DB creation command
	dbCreateCmd := flag.NewFlagSet("db-add-recipes", flag.ExitOnError)
	dbCreateDbPath := dbCreateCmd.String("dbpath", "./recipes.db", "Path to the recipe sqlite database")
	DbCreateRecipeFileDir := dbCreateCmd.String("recipes", "static/recipe_mds/", "Directory to the recipe md files")

	// Run DB migration
	dbMigrateCmd := flag.NewFlagSet("db-migrate", flag.ExitOnError)
	dbMigrateDbPath := dbMigrateCmd.String("dbpath", "./recipes.db", "Path to the recipe sqlite database")

	// AWS download (sync) assets command
	syncCmd := flag.NewFlagSet("sync-from-aws", flag.ExitOnError)
	bucket := syncCmd.String("bucket", "", "Name of the AWS bucket where assets are stored")
	dataPath := syncCmd.String("data-path", "", "Path to where static assets are stored locally")

	// Start TCP connection with fly app
	tcpCmd := flag.NewFlagSet("start-tcp", flag.ExitOnError)
	tcpAddress := tcpCmd.String("address", "", "Address to connect the sender to")
	directoryToSync := tcpCmd.String("directory", "", "Directory to sync md files to")
	pingOnly := tcpCmd.Bool("ping-only", false, "Only send the authenticatation check to the server to check it's recieving")
	replica := tcpCmd.Bool("replica", false, "Run this cmd as the replica")

	// AWS upload all assets command
	wrongSubCommandMsg := "expected the command 'db-add-recipes', 'db-migrate', 'start-tcp' or 'sync-from-aws'"

	if len(os.Args) < 2 {
		fmt.Println(wrongSubCommandMsg)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "db-add-recipes":
		dbCreateCmd.Parse(acceptDoubleDashArgs(os.Args[2:]))
		db, err := database.NewSqlDatabase(*dbCreateDbPath, true)
		if err != nil {
			log.Fatal(err)
		}
		db.AddFiles(*DbCreateRecipeFileDir)

	case "db-migrate":
		db, err := database.CreateDbConnection(*dbMigrateDbPath)
		if err != nil {
			log.Fatal("Failed to create db connection to %s: %s\n", *dbMigrateDbPath, err)
		}
		database.RunDbMigrations(db)

	case "sync-from-aws":
		syncCmd.Parse(acceptDoubleDashArgs(os.Args[2:]))
		if *bucket == "" {
			log.Fatal("--bucket the AWS bucket must be set")
		} else {
			syncFromAws(*bucket, *dataPath)
		}
	case "start-tcp":
		tcpCmd.Parse(acceptDoubleDashArgs(os.Args[2:]))
		apiKey, apiKeyExists := os.LookupEnv("TCP_API_KEY")
		if !apiKeyExists {
			log.Fatal("No TCP_API_KEY set as an env")
		}
		if *replica {
			log.Println("Running replica tcp file server")
			err := filesyncer.StartReplicaTCPFileServer(*tcpAddress, apiKey, *directoryToSync)
			if err != nil {
				log.Fatalf("Failed to run the replica server: %s\n", err)
			}
		} else {

			tls := true
			if strings.HasPrefix(*tcpAddress, "0.0.0.0") || strings.HasPrefix(*tcpAddress, "127.0.0.1") {
				tls = false
			}
			msg := "Running main tcp file client"
			if *pingOnly {
				msg += " (ping only)"
			}
			if !tls {
				msg += " (no tls)"
			}
			log.Println(msg)
			startMainTCP(*tcpAddress, apiKey, *directoryToSync, *pingOnly, tls)
		}

	default:
		fmt.Println(wrongSubCommandMsg)
		os.Exit(1)
	}

}
