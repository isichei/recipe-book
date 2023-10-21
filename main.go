package main

import (
	"flag"
	"log"

	"github.com/isichei/recipe-book/api"
)

func main() {

	listenAddr := flag.String("listenaddr", ":8000", "The address for the API to listen on")

	flag.Parse()
	server := api.NewServer(*listenAddr)

	log.Fatal(server.Start())
}
