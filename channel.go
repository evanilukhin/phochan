package phochan

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type Channel struct {
	Room                   string
	Socket                 *Socket
	ControlChan            chan ControlChanMessage
	Joined                 bool
	IncomingMessageHandler func(string)
}

type ControlChanMessage struct {
	Type    string
	Message string
}

type MessagePayload struct {
	Body string `json:"body"`
}

func (channel *Channel) Start() {
	go func() {
		for {
			_, p, _ := channel.Socket.Conn.ReadMessage()
			channel.IncomingMessageHandler(string(p[:]))
		}
	}()
	go func() {
		for {
			strD := [5]string{"1", "1", "phoenix", "heartbeat", "{}"}
			strB, _ := json.Marshal(strD)
			channel.Socket.Conn.WriteMessage(websocket.TextMessage, strB)
			time.Sleep(channel.Socket.HeartbeatInterval)
		}
	}()
}

func (channel *Channel) Join() string {
	strD := [5]string{"1", "1", channel.Room, "phx_join", "{}"}
	strB, _ := json.Marshal(strD)
	channel.Socket.Conn.WriteMessage(websocket.TextMessage, strB)
	return "stub"
}

func (channel *Channel) Push(action, message string) string {
	strD := []interface{}{"1", "2", channel.Room, action, MessagePayload{message}}
	strB, _ := json.Marshal(strD)
	channel.Socket.Conn.WriteMessage(websocket.TextMessage, strB)
	return "stub"
}
