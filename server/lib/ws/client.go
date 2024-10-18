package ws

import (
	"encoding/json"
	"os"

	"github.com/gorilla/websocket"
	"xes.software/rpg/lib/utils"
)

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	JoinRaid chan *Raid
	raid     *Raid
	Hub      *Hub
}

type ClientMessage struct {
	Type ClientMessageType `json:"type"`
	Data json.RawMessage   `json:"data"`
}

type ClientMessageJoinRaid struct {
	Uuid string `json:"uuid"`
}
type ClientMessageSendMessage struct {
	Message string `json:"message"`
}

type ClientMessageType string

const (
	JoinRaid    ClientMessageType = "join_raid"
	SendMessage ClientMessageType = "send_message"
)

func CreateClient(conn *websocket.Conn, logger *utils.Logger, h *Hub) *Client {
	client := &Client{
		Conn:     conn,
		Send:     make(chan []byte),
		JoinRaid: make(chan *Raid),
		Hub:      h,
	}
	logger.DevLog(os.Stdout, "creating a new client: %+v", *client)
	go client.Run(logger, h)
	return client
}

func (c *Client) readPump(logger *utils.Logger, hub *Hub) {
	defer c.Conn.Close()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.DevLog(os.Stdout, "error: %v", err)
			}
			return
		}
		// use message
		logger.DevLog(os.Stdout, "recieved a message: %s on conn: %s", string(message), c.Conn.RemoteAddr().String())
		// marshall json
		var clientMsg ClientMessage
		err = json.Unmarshal(message, &clientMsg)
		if err != nil {
			logger.DevLog(os.Stdout, "message recieved could not be parsed: %s", err.Error())
			continue
		}
		if clientMsg.Type == JoinRaid {
			// join raid
			var joinRaidData ClientMessageJoinRaid
			err = json.Unmarshal(clientMsg.Data, &joinRaidData)
			if err != nil {
				logger.DevLog(os.Stdout, "joinRaidData could not be parsed: %s", err.Error())
				continue
			}
			logger.DevLog(os.Stdout, "client requested to join raid with id: %s", joinRaidData.Uuid)
			hub.StartRaid <- joinRaidData.Uuid
		} else if clientMsg.Type == SendMessage {
			// send message
			var sendMessageData ClientMessageSendMessage
			err = json.Unmarshal(clientMsg.Data, &sendMessageData)
			if err != nil {
				logger.DevLog(os.Stdout, "sendMessageData could not be parsed: %s", err.Error())
				continue
			}
			logger.DevLog(os.Stdout, "client sent message: %s", sendMessageData.Message)
		}

	}
}

func (c *Client) writePump(logger *utils.Logger) {
	for {
		select {
		case message, ok := <-c.Send:
			logger.DevLog(os.Stdout, "recieved a message:%s ok:%t", string(message), ok)
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.DevLog(os.Stdout, "error while closing connection: %s", err.Error())
				}
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.DevLog(os.Stdout, "error while writing a message to connection: %s", err.Error())
				return
			}
		case raid, ok := <-c.JoinRaid:
			logger.DevLog(os.Stdout, "recieved a JoinRaid event:%s ok:%t", string(raid.Uuid.String()), ok)
			c.raid = raid
		}
	}
}

/* Spawns 2 goroutines, one for handling reads on a connection, and one for writes. */
func (c *Client) Run(logger *utils.Logger, h *Hub) {
	go c.writePump(logger)
	go c.readPump(logger, h)
}
