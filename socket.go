package phochan

import (
	"github.com/gorilla/websocket"
	"time"
)

type Socket struct {
	Address           string
	HeartbeatInterval time.Duration // in milliseconds
	Channels          []*Channel
	Conn              *websocket.Conn
}

func NewSocket(Address string) *Socket {
	socket := new(Socket)
	socket.Address = Address
	socket.HeartbeatInterval = time.Second * 30 // @warning HARDCODE!!! Default heartbeat value in seconds
	return socket
}

//Channel creates and adds new Channel to socket.Channels
func (socket *Socket) Channel(room string, incomingMessageHandler func(PhoenixMessage)) *Channel {
	channel := Channel{Socket: socket, Room: room, IncomingMessageHandler: incomingMessageHandler}
	socket.Channels = append(socket.Channels, &channel)
	return &channel
}

//Connect sets gorilla/websocket Conn for socket
func (socket *Socket) Connect() error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(socket.Address, nil)
	socket.Conn = conn
	return err
}
