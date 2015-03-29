package main

import (
	"fmt"
	"github.com/ovaskevich/TxtRoulette/server"
	"log"
	"net/http"
	"os"
)

func main() {
	// Read the args.
	if len(os.Args) != 2 {
		log.Fatal("usage: server.go port")
	}
	port := ":" + os.Args[1]

	// Start the server.
	fmt.Printf("Starting TxtRoulette server on port %s...\n", port)
	http.HandleFunc("/receive/", server.Receive)
	log.Fatal(http.ListenAndServe(port, nil))
}
