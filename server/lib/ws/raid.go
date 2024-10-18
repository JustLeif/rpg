package ws

import (
	"os"

	"github.com/google/uuid"
	"xes.software/rpg/lib/utils"
)

type Raid struct {
	clients          map[*Client]bool
	Uuid             uuid.UUID
	state            string
	broadcast        chan []byte
	ConnectClient    chan *Client
	DisconnectClient chan *Client
	SetState         chan string
}

/* returns a raid pointer, and start's the raid on it's own goroutine (will clean up itself). */
func CreateRaid(logger *utils.Logger) *Raid {
	raid := &Raid{
		clients:          make(map[*Client]bool),
		Uuid:             uuid.New(),
		state:            "",
		broadcast:        make(chan []byte),
		ConnectClient:    make(chan *Client),
		DisconnectClient: make(chan *Client),
		SetState:         make(chan string),
	}
	go raid.Run(logger)
	logger.DevLog(os.Stdout, "raid %s was created.", raid.Uuid.String())
	return raid
}

func (r *Raid) Run(logger *utils.Logger) {
	defer logger.DevLog(os.Stdout, "raid %s was destroyed and resources freed.", r.Uuid.String())
	for {
		select {
		case broadcast, ok := <-r.broadcast:
			logger.DevLog(os.Stdout, "recieved a broadcast: %s, ok: %t ", string(broadcast), ok)
			if ok == false {
				break
			}

		case client, ok := <-r.ConnectClient:
			logger.DevLog(os.Stdout, "recieved a connect client: %s, ok: %t ", client.Conn.RemoteAddr().String(), ok)
			if ok == false {
				break
			}

		case client, ok := <-r.DisconnectClient:
			logger.DevLog(os.Stdout, "recieved a disconnect client: %s, ok: %t ", client.Conn.RemoteAddr().String(), ok)
			if ok == false {
				break
			}

		case state, ok := <-r.SetState:
			logger.DevLog(os.Stdout, "recieved a set state: %s, ok: %t ", state, ok)
			if ok == false {
				break
			}
		}
	}
}
