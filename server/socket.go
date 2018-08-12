package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxMessage = 1028
)

// Client is our client
type Client struct {
	Connections *Connections          `json:"-"`
	Conn        *websocket.Conn       `json:"-"`
	Send        chan *OutgoingMessage `json:"-"`
	Channel     *Channel              `json:"-"`

	Username string `json:"username"`
	Vote     int    `json:"vote"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

func serveWS(conns *Connections, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade:", err)
		return
	}

	client := &Client{
		Connections: conns,
		Conn:        conn,
		Send:        make(chan *OutgoingMessage),
	}
	client.Connections.Register <- client

	go client.read()
	go client.write()
}

func (c *Client) read() {
	defer func() {
		c.Connections.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessage)

	for {
		newMsg := ReceiveMessage{}
		err := c.Conn.ReadJSON(&newMsg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
				websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Client Read Error:", err)
				return
			}
			c.Send <- &OutgoingMessage{Type: typeResponse, Error: errBadMessage}
		}

		switch newMsg.Type {
		case typeSetUsername:
			go func(c *Client, newMsg ReceiveMessage) {
				reg := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
				matches := reg.Find([]byte(newMsg.Data.Username))
				if newMsg.Data.Username == "" || len(matches) == 0 {
					c.Send <- &OutgoingMessage{Type: typeSetUsername, Error: errBadUsername}
					return
				}
				c.Username = newMsg.Data.Username
				c.Send <- &OutgoingMessage{Type: typeSetUsername, Data: &MessageData{Message: "Username set to " + c.Username, Timestamp: time.Now()}}
			}(c, newMsg)
		case typeVote:
			go func(c *Client, newMsg ReceiveMessage) {
				if c.Channel == nil {
					c.Send <- &OutgoingMessage{Type: typeVote, Error: errUnauthorized}
					return
				}

				vote := -1
				if newMsg.Data.Vote {
					vote = 1
				}
				c.Vote = vote

				for _, member := range c.Channel.Members {
					member.Send <- &OutgoingMessage{Type: typeVote, Data: &MessageData{Sender: c, Timestamp: time.Now()}}
				}
			}(c, newMsg)
		case typeCreateChannel:
			go c.createChannel(newMsg)
		case typeJoinChannel:
			go c.joinChannel()
		case typeLeaveChannel:
			go c.leaveChannel()
		case typeSendMessage:
			go c.sendMessage(newMsg)
		}
	}
}

func (c *Client) write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteJSON(msg)
		}
	}
}
