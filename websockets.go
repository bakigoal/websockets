package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := wsUpgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		// Read message from browser
		go readMessages(conn)
		publishMessages(conn)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}

func readMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("\n%s sent: %s\n", conn.RemoteAddr(), string(msg))
	}
}

func publishMessages(conn *websocket.Conn) {
	for {
		// Write message back to browser
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		msgBack := []byte(text)
		if err := conn.WriteMessage(1, msgBack); err != nil {
			return
		}
	}
}
