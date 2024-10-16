package ws

import (
	"os"
	"sync"

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
func CreateRaid(logger *utils.Logger, rm *RaidsMap) *Raid {
	raid := &Raid{
		clients:          make(map[*Client]bool),
		Uuid:             uuid.New(),
		state:            "",
		broadcast:        make(chan []byte),
		ConnectClient:    make(chan *Client),
		DisconnectClient: make(chan *Client),
		SetState:         make(chan string),
	}
	rm.Store(raid.Uuid.String(), raid)
	return raid
}

func (r *Raid) run(logger *utils.Logger, rm *RaidsMap) {
	rm.Store(r.Uuid.String(), r)
	defer rm.Delete(r.Uuid.String())
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

/* thread safe map of [raid.Uuid.ToString()] to [*Raid] */
type RaidsMap struct {
	internal sync.Map
}

func (m *RaidsMap) Store(key string, v *Raid) {
	m.internal.Store(key, v)
}

func (m *RaidsMap) Load(key string) (*Raid, bool) {
	if value, ok := m.internal.Load(key); ok {
		if myStruct, ok := value.(*Raid); ok {
			return myStruct, true
		}
	}
	return nil, false
}

func (m *RaidsMap) Delete(key string) {
	m.internal.Delete(key)
}
