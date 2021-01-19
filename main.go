package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rakyll/portmidi"
)

var deviceIDFlag = flag.Int("deviceid", 3, "MIDI Device ID")

func main() {
	flag.Parse()

	// midi
	deviceID := portmidi.DeviceID(*deviceIDFlag)
	fmt.Println(portmidi.Info(deviceID))

	in, err := portmidi.NewInputStream(deviceID, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	midiEvent := in.Listen()

	clients := map[int64]*websocket.Conn{}

	go func() {
		for {
			select {
			case event := <-midiEvent:
				log.Printf("write midi: %+v\n", event)
				for _, client := range clients {
					err := client.WriteJSON(event)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}()

	// server
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	http.HandleFunc("/midi", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		uuid := time.Now().UnixNano()
		clients[uuid] = conn

		log.Printf("client connected: %d\n", uuid)

		defer func() {
			conn.Close()
			delete(clients, uuid)
			log.Printf("client disconnected: %d\n", uuid)
		}()

		for {
			// err is returned when the client is disconnected
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
