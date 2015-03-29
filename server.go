package main

import (
	"fmt"
	//"io"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	//"strings"
)

func receive(w http.ResponseWriter, r *http.Request) {

	wholeurl := r.URL.String()
	body := r.URL.Query()["Body"]
	phone := r.URL.Query()["From"]

	fmt.Printf("wholeurl:\n%s\n\nPhone: %s\nBody: %s,\n\n", wholeurl, phone, body)
}

func main() {
	// Read the port from the file.
	//fileContents, err := ioutil.ReadFile("port.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//port := strings.TrimSpace(string(fileContents))

	if len(os.Args) != 2 {
		log.Fatal("usage: server.go port")
	}

	port := ":" + os.Args[1]

	// Start the server.
	fmt.Printf("Starting TxtRoulette server on port %s...\n", port)
	http.HandleFunc("/receive/", receive)
	log.Fatal(http.ListenAndServe(port, nil))
}
