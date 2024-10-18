package ws

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"xes.software/rpg/lib/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWs(envVars *utils.EnvVars, logger *utils.Logger, hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.DevLog(os.Stdout, "upgrader.Upgrade() failed: %s", err.Error())
			return
		}
		hub.NewConnection <- conn
	}
}
