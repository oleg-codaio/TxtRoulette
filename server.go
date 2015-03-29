package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type User struct {
	phoneNumber string
	recent      [5]*User
	blocked     []*User
}

func receive(w http.ResponseWriter, r *http.Request) {

	wholeurl := r.URL.String()
	body := r.URL.Query()["Body"]
	phone := r.URL.Query()["From"]

	fmt.Printf("wholeurl:\n%s\n\nPhone: %s\nBody: %s,\n\n", wholeurl, phone, body)
}

func sendSMS(phonenumber, message string) {

	apiusr := os.Getenv("TWILIO_APIUSR")
	apikey := os.Getenv("TWILIO_APIKEY")

	u := "https://api.twilio.com/2010-04-01/Accounts/AC7dbbd979132aeb252095fa79059a5de4/Messages.json"

	hc := http.Client{}
	form := url.Values{}
	form.Add("To", phonenumber)
	form.Add("From", "+13208398785")
	form.Add("Body", message)

	req, err := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(apiusr, apikey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//fmt.Printf("the request was: \n%v\n\n",req)

	resp, err := hc.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 201 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(body)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: server.go port")
	}

	port := ":" + os.Args[1]

	// Start the server.
	fmt.Printf("Starting TxtRoulette server on port %s...\n", port)
	http.HandleFunc("/receive/", receive)
	log.Fatal(http.ListenAndServe(port, nil))
}
