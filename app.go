package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", websocketHandler)
	fmt.Println("Server started on localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()
	log.Println("Client connected")

	for {
		// Read message from the client
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Println("WebSocket connection closed by client")
				break
			}
			log.Println("WebSocket read error:", err)
			break
		}

		// Reverse the text
		text := string(msgBytes)
		reversedText := reverseString(text)

		// Write reversed text back to the client
		err = conn.WriteMessage(websocket.TextMessage, []byte(reversedText))
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

func reverseString(s string) string {
	runes := []int32(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
