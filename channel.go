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
	IncomingMessageHandler func(PhoenixMessage)
}

type ControlChanMessage struct {
	Type    string
	Message string
}

type MessagePayload struct {
	Body string `json:"body"`
}

type PhoenixMessage struct {
	JoinRef string
	Ref     string
	Topic   string
	Event   string
	Payload json.RawMessage
}

func (channel *Channel) Start() {
	go readMessages(channel)
	go heartbeat(channel)
}

func readMessages(channel *Channel) {
	for {
		var phoenixResponce []json.RawMessage
		var joinRef, ref, topic, event string
		_, p, _ := channel.Socket.Conn.ReadMessage()
		json.Unmarshal(p, &phoenixResponce)
		json.Unmarshal(phoenixResponce[0], &joinRef)
		json.Unmarshal(phoenixResponce[1], &ref)
		json.Unmarshal(phoenixResponce[2], &topic)
		json.Unmarshal(phoenixResponce[3], &event)
		if !((ref == "1") && (event == "phx_reply")) {
			channel.IncomingMessageHandler(
				PhoenixMessage{
					JoinRef: joinRef,
					Ref:     ref,
					Topic:   topic,
					Event:   event,
					Payload: phoenixResponce[4],
				},
			)
		}
	}
}

func heartbeat(channel *Channel) {
	for {
		strD := [5]string{"1", "1", "phoenix", "heartbeat", "{}"}
		strB, _ := json.Marshal(strD)
		channel.Socket.Conn.WriteMessage(websocket.TextMessage, strB)
		time.Sleep(channel.Socket.HeartbeatInterval)
	}
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
