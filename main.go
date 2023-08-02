package main

import (
	"flag"
	"log"

	"github.com/isichei/recipe-book/api"
	"github.com/isichei/recipe-book/storage"
)

func main() {
	listenPort := flag.String("listen-port", ":8000", "What port to serve the app on")
	server := api.NewServer(*listenPort, storage.NewFakeStorage())
	log.Fatal(server.Start())
}
