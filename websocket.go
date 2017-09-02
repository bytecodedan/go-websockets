package main

import (
	"net/http"

	"github.com/ChimeraCoder/anaconda"
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
	Data string `json:"data,omitempty"`
}

var HandleWSConnections = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	go handleConn(conn)
})

var twitter Twitter

func InitWS() {
	twitter = Twitter{}
	twitter.Init()
}

func handleConn(conn *websocket.Conn) {
	// channel to be used by one Go routine to
	// signal the other to stop processing and
	// exit.
	var exit = make(chan bool)

	go func() {
		for {
			msg := wsEvent{}
			err := conn.ReadJSON(&msg)

			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				// check for client closed connections
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					log.Info("Websocket closed by client.\n")
					// stop current tweet stream
					twitter.Stop()
					// tell bottom Go routine to stop processing current
					// tweet stream and exit.
					exit <- true
					return
				}
			}

			if ok := twitter.SetupStream(msg.Data); !ok {
				log.Error("Unexpected error creating tweet stream.")
				return
			}

			go func() {
				for {
					select {
					case v := <-twitter.Stream.C:
						tweet, ok := v.(anaconda.Tweet)
						if !ok {
							log.Warningf("Received unexpected value type of %T\n", v)
							return
						}
						log.Infof("%v\n", tweet.Text)
						if err := conn.WriteJSON(tweet.Text); err != nil {
							log.Error(err)
							return
						}
					case <-exit:
						return
					}
				}
			}()
		}
	}()
}
