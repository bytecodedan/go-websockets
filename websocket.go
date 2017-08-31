package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type wsEvent struct {
	// Data interface{} `json:"data,omitempty"`
	Data string `json:"data,omitempty"`
	Type string `json:"type,omitempty"`
}

var HandleWSConnections = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	go handleMessages(conn)
})

func handleMessages(conn *websocket.Conn) {
	go func() {
		for {
			msg := wsEvent{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("Error reading json.", err)
			}
			tweetReq <- msg.Data
		}
	}()

	go func() {
		for {
			tweet := <-tweetRes
			if err := conn.WriteJSON(tweet); err != nil {
				fmt.Println(err)
			}
		}
	}()
}
