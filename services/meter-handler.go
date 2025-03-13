package services

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vanspaul/SmartMeterSystem/utils"
)

// Upgrader defines the parameters for upgrading an HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (customize as needed)
	},
}

func MeterHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	utils.Logger.Sugar().Debugf("Client connected")

	// Create a channel to signal client disconnection
	done := make(chan struct{})

	// Goroutine to read messages from client
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Println("Client disconnected normally")
				} else {
					log.Println("Read error:", err)
				}
				close(done) // Signal disconnection
				return      // Exit the goroutine
			}
			log.Println("Received:", string(message))
		}
	}()

	// Send periodic messages to client
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			message := []byte("Hello from Go Server - " + time.Now().Format(time.RFC3339))
			log.Println("Sending:", string(message))
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Write error:", err)
				return
			}
		case <-done:
			log.Println("Stopping sending due to client disconnection")
			return
		}
	}
}
