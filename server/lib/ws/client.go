package ws

import (
	"encoding/json"
	"os"

	"github.com/gorilla/websocket"
	"xes.software/rpg/lib/utils"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	raid *Raid
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

func CreateClient(conn *websocket.Conn, logger *utils.Logger) *Client {
	client := &Client{
		Conn: conn,
		Send: make(chan []byte),
	}
	logger.DevLog(os.Stdout, "creating a new client: %v", client)
	return client
}

func (c *Client) readPump(logger *utils.Logger, rm *RaidsMap) {
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
			// check if raid currently exists
			raid, ok := rm.Load(joinRaidData.Uuid)
			if !ok {
				raid = CreateRaid(logger, rm)
				logger.DevLog(os.Stdout, "created a raid with id: %s", raid.Uuid)
			}
			raid.ConnectClient <- c
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
		}
	}
}

/* Spawns 2 goroutines, one for handling reads on a connection, and one for writes. */
func (c *Client) Run(logger *utils.Logger, rm *RaidsMap) {
	go c.writePump(logger)
	go c.readPump(logger, rm)
}
