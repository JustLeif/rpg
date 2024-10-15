package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"xes.software/2d-game/lib/utils"
)

// Define an upgrader for WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		return true
	},
}

// Thread-safe map to store connections
var connections sync.Map

// Handle incoming WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Store the connection in the map
	connections.Store(conn.RemoteAddr(), conn)
	fmt.Println("New connection from:", conn.RemoteAddr())
	defer connections.Delete(conn.RemoteAddr())

	// Read messages from WebSocket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		fmt.Printf("Received message from %v: %s\n", conn.RemoteAddr(), string(msg))
		// Broadcast the message to all connections
		broadcastMessage(msg)
	}
}

// Broadcast a message to all WebSocket connections
func broadcastMessage(msg []byte) {
	connections.Range(func(key, value interface{}) bool {
		conn, ok := value.(*websocket.Conn)
		if !ok {
			return true
		}

		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			conn.Close()
			connections.Delete(key)
		}
		return true
	})
}

func main() {

	vars := utils.LoadEnvVars()
	logger := utils.NewLogger(vars.Environment)
	logger.DevLog(os.Stdout, "Hello!\n")

	http.HandleFunc("/ws", handleWebSocket)

	// Start WebSocket server
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
