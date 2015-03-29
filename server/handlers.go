package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

// Map from phone number <=> User.
var users = make(map[string]*User)

// Pairs of currently connected Users.
var pairs = make(map[*User]*User)

// Lobby of unpaired Users.
var lobby = make(map[*User]bool, 0)

// Mutex for writing to state.
var mutex = sync.Mutex

func MatchLobbiedUsers() {
	if len(lobby) < 2 {
		return
	}
	for u1 := range lobby {
		for u2 := range lobby {
			if u1 == u2 { //don't match with self
				continue
			}
			for _, block := range u1.blocked {
				if block == u2 {
					continue //u1 blocked u2
				}
			}
			for _, block := range u2.blocked {
				if block == u1 {
					continue //u2 blocked u1
				}
			}

			pairs[u1] = u2
			pairs[u2] = u1
			sendSMS(u1.phoneNumber, "you've been paired with another person, say hi!")
			sendSMS(u2.phoneNumber, "you've been paired with another person, say hi!")
			fmt.Printf("paired %s & %s\n", u1.phoneNumber, u2.phoneNumber)
			//remove the users from the lobby
			delete(lobby, u1)
			delete(lobby, u2)
			return
		}
	}
	return
}

func Receive(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	msg := strings.TrimSpace(r.URL.Query()["Body"][0])
	num := r.URL.Query()["From"][0]

	user, isRegistered := users[num]
	_, isInLobby := lobby[user]
	_, isPaired := pairs[user]

	if !isRegistered && msg != "CONNECT" {
		// TODO: reply using ResponseWriter instead of calling SendSMS
		sendInstructions(num)
		return
	}

	switch msg {
	case "CONNECT":
		if !isRegistered {
			// Add user to users
			u := NewUser(num)
			users[num] = u
		} else if isPaired {
			//sendSMS(num, "Invalid command. Use NEXT, DISCONNECT, or BLOCK.")
			sendSMS(num, "You're already chatting with someone! txt DISCONNECT to leave or try talking to someone new with NEXT.")
			return
		} else if isInLobby {
			sendSMS(num, "Please wait! We're still trying to connect you...")
			return
		}

		// Add the user to the lobby
		lobby[users[num]] = true
		fmt.Printf("added  %s to lobby\n", num)
		sendSMS(num, "Hang tight, we're trying to connect you...")
		// Try to connect user if there is someone free in the lobby
		//check if another user in lobby
		MatchLobbiedUsers()
		return

	case "DISCONNECT":
		if isPaired {
			lobby[pairs[users[num]]] = true
			sendSMS(num, "You have successfully disconnected.")
			sendSMS(pairs[users[num]].phoneNumber, "You're partner left the chat, please hang tight while we find someone new to chat with.")
			delete(pairs, pairs[users[num]])
			delete(pairs, users[num])
			// Unpair them
			MatchLobbiedUsers()
		} else if isInLobby {
			delete(lobby, users[num])
			sendSMS(num, "You have successfully disconnected.")
			// Remove them from the lobby
		} else {
			sendSMS(num, "You're already disconnected! You can reconnect by texting CONNECT.")
		}
		return
	case "NEXT":
		if isPaired {
			lobby[pairs[users[num]]] = true
			lobby[users[num]] = true
			// Unpair them
			sendSMS(num, "You left the chat, please hang tight while we find someone new to chat with.")
			sendSMS(pairs[users[num]].phoneNumber, "You're partner left the chat, please hang tight while we find someone new to chat with.")
			delete(pairs, pairs[users[num]])
			delete(pairs, users[num])
			MatchLobbiedUsers()
		} else if isInLobby {
			sendSMS(num, "Please wait! We're still trying to find someone for you to chat with...")
		}
		if isRegistered {
			sendSMS(num, "Hang tight, we're trying to connect you...")
			lobby[users[num]] = true
			MatchLobbiedUsers()
		}
		return
	case "BLOCK":
		if isPaired {
			//TODO: actually block the other person
			sendSMS(num, "you've blocked the other user and been added to the lobby")
			sendSMS(pairs[users[num]].phoneNumber, "You're partner left the chat, please hang tight while we find someone new to chat with.")

			// Add paired number to user's block list
			users[num].blocked = append(users[num].blocked, pairs[users[num]])
			// Put the user back into the lobby
			lobby[pairs[users[num]]] = true
			lobby[users[num]] = true
			// Unpair them
			delete(pairs, pairs[users[num]])
			delete(pairs, users[num])
			MatchLobbiedUsers()
		} else {
			sendSMS(num, "You're not currently chatting with anyone.")
		}
		return
	default:
		if isPaired {
			sendSMS(pairs[users[num]].phoneNumber, msg)
			// Send the msg to the paired user.
		} else {
			sendInstructions(num)
		}
		return
	}
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

func sendInstructions(phoneNumber string) {
	sendSMS(phoneNumber, "Welcome to TXTRoulette! commands to txt are:\nCONNECT - start chatting\nNEXT - chat with someone new\nDISCONNECT - stop chatting\nBLOCK - block chat partner")
}
