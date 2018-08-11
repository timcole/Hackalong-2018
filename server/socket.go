package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

const (
	timeout    = 5 * time.Minute
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
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(timeout))
		return nil
	})

	for {
		newMsg := ReceiveMessage{}
		err := c.Conn.ReadJSON(&newMsg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Client Read Error:", err)
				break
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
				c.Send <- &OutgoingMessage{Type: typeSetUsername, Data: &MessageData{Message: "Username set to " + c.Username}}
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
	keepAlive := time.NewTicker(timeout)
	defer func() {
		keepAlive.Stop()
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
		case <-keepAlive.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
