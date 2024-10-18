package ws

import (
	"os"

	"github.com/gorilla/websocket"
	"xes.software/rpg/lib/utils"
)

type Hub struct {
	/* map[RaidUuid]*Raid */
	Raids            map[string]*Raid
	Connections      map[*websocket.Conn]*Client
	NewConnection    chan *websocket.Conn
	DeleteConnection chan *websocket.Conn
	EndRaid          chan string
	StartRaid        chan string
}

/* Creates a new hub on a goroutine, cleans itself up automatically. */
func CreateHub(logger *utils.Logger) *Hub {
	hub := &Hub{
		Raids:         make(map[string]*Raid),
		Connections:   make(map[*websocket.Conn]*Client),
		StartRaid:     make(chan string),
		EndRaid:       make(chan string),
		NewConnection: make(chan *websocket.Conn),
	}
	logger.DevLog(os.Stdout, "creating a hub.")
	go hub.Run(logger)
	return hub
}

func (h *Hub) Run(logger *utils.Logger) {
	for {
		select {
		case conn, ok := <-h.NewConnection:
			logger.DevLog(os.Stdout, "recieved a NewConnection request: %s, ok: %t ", conn.RemoteAddr(), ok)
			if ok == false {
				break
			}
			client := CreateClient(conn, logger, h)
			h.Connections[conn] = client

		case conn, ok := <-h.DeleteConnection:
			logger.DevLog(os.Stdout, "recieved a DeleteConnection request: %s, ok: %t ", conn.RemoteAddr(), ok)
			if ok == false {
				break
			}
			delete(h.Connections, conn)

		case raidUuid, ok := <-h.StartRaid:
			logger.DevLog(os.Stdout, "recieved a StartRaid request: %s, ok: %t ", raidUuid, ok)
			if ok == false {
				break
			}
			raid := CreateRaid(logger)
			h.Raids[raid.Uuid.String()] = raid

		case raidUuid, ok := <-h.EndRaid:
			logger.DevLog(os.Stdout, "recieved a DeleteRaid request: %s, ok: %t ", raidUuid, ok)
			if ok == false {
				break
			}
			delete(h.Raids, raidUuid)
		}
	}
}
