package main

import (
	"sync"
)

// Connections are our connections
type Connections struct {
	Clients    map[*Client]bool
	Inbound    chan []byte
	Register   chan *Client
	Unregister chan *Client
	Channels   []*Channel
	chanMutex  sync.Mutex
}

// NewConnections inits the connections handler
func NewConnections() *Connections {
	return &Connections{
		Clients:    make(map[*Client]bool),
		Inbound:    make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run listens for new websocket connections
func (conns *Connections) Run() {
	for {
		select {
		case client := <-conns.Register:
			conns.Clients[client] = true
		case client := <-conns.Unregister:
			if _, ok := conns.Clients[client]; ok {
				client.leaveChannel()
				delete(conns.Clients, client)
				close(client.Send)
			}
		}
	}
}

type wsError string
type wsType string

const (
	errServer       wsError = "ERR_SERVER"
	errBadMessage   wsError = "ERR_BADMESSAGE"
	errBadUsername  wsError = "ERR_BADUSERNAME"
	errUnauthorized wsError = "ERR_UNAUTHORIZED"
	errNoChannels   wsError = "ERR_NOCHANNELS"

	typeResponse      wsType = "RESPONSE"
	typeSetUsername   wsType = "SET_USERNAME"
	typeCreateChannel wsType = "CREATE_CHANNEL"
	typeJoinChannel   wsType = "JOIN_CHANNEL"
	typeLeaveChannel  wsType = "LEAVE_CHANNEL"
	typeSendMessage   wsType = "SEND_MESSAGE"
	typeNewMessage    wsType = "NEW_MESSAGE"
	typeMemberJoin    wsType = "MEMBER_JOIN"
	typeMemberLeave   wsType = "MEMBER_LEAVE"
)

// MessageData is the message data
type MessageData struct {
	Topic    string `json:"topic,omitempty"`
	Message  string `json:"message,omitempty"`
	Username string `json:"username,omitempty"`
}

// ReceiveMessage is the message from clients
type ReceiveMessage struct {
	Type wsType       `json:"type"`
	Data *MessageData `json:"data,omitempty"`
}

// OutgoingMessage is the message to the clients
type OutgoingMessage struct {
	Type  wsType       `json:"type"`
	Error wsError      `json:"error,omitempty"`
	Data  *MessageData `json:"data,omitempty"`
}
