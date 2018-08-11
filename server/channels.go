package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message is structure for the message history
type Message struct {
	Sender    *Client `json:"member"`
	Timestamp string  `json:"timestamp"`
	Message   string  `json:"message"`
}

// Channel is the chatroom
type Channel struct {
	Topic    string        `json:"topic"`
	Duration time.Duration `json:"duration"`
	Members  []*Client     `json:"members"`
	History  []*Message    `json:"history"`
	Slot     int           `json:"slots"`
	mutex    sync.Mutex
}

func (c *Client) createChannel(msg ReceiveMessage) {
	if c.Channel != nil || c.Username == "" {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errUnauthorized}
		return
	}

	channel := &Channel{
		Topic:    msg.Data.Topic,
		Duration: time.Nanosecond,
		Members:  []*Client{c},
		Slot:     2,
	}

	c.Connections.chanMutex.Lock()
	c.Connections.Channels = append(c.Connections.Channels, channel)
	c.Connections.chanMutex.Unlock()
	c.Channel = channel

	fmt.Println(c.Connections.Channels, len(c.Connections.Channels))

	c.Send <- &OutgoingMessage{Type: typeJoinChannel}
}

func (c *Client) joinChannel() {
	if c.Username == "" {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errUnauthorized}
		return
	}

	var openChannels []*Channel

	for _, channel := range c.Connections.Channels {
		if channel.Slot == 0 {
			continue
		}

		openChannels = append(openChannels, channel)
	}

	if len(openChannels) == 0 {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errNoChannels}
		return
	}

	rand.Seed(time.Now().UnixNano())
	var join = openChannels[rand.Intn(len(openChannels))]

	join.mutex.Lock()
	if join.Slot == 0 {
		c.joinChannel()
		return
	}

	c.Channel = join
	join.Members = append(join.Members, c)
	join.Slot--
	join.mutex.Unlock()

	c.Send <- &OutgoingMessage{Type: typeJoinChannel, Data: &MessageData{Topic: join.Topic}}
}

func (c *Client) leaveChannel() {
	if c.Channel == nil {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errNoChannels}
		return
	}

	c.Channel.mutex.Lock()
	var tmpMembers []*Client
	var members = c.Channel.Members
	for i := range members {
		members[len(members)-1], members[i] = members[i], members[len(members)-1]
		tmpMembers = members[:len(members)-1]
	}
	c.Channel.Members = tmpMembers
	c.Channel.Slot++

	if len(c.Channel.Members) == 0 && c.Channel.Slot == 3 {
		c.Channel.Slot = 0
		go c.Connections.close(c.Channel)
	}
	c.Channel.mutex.Unlock()
	c.Channel = nil

	c.Send <- &OutgoingMessage{Type: typeLeaveChannel}
}

func (conns *Connections) close(channel *Channel) {
	conns.chanMutex.Lock()
	var tmpChannels []*Channel
	var channels = conns.Channels
	for i := range channels {
		channels[len(channels)-1], channels[i] = channels[i], channels[len(channels)-1]
		tmpChannels = channels[:len(channels)-1]
	}
	conns.Channels = tmpChannels
	conns.chanMutex.Unlock()
}

func (c *Client) sendMessage(msg ReceiveMessage) {
	if c.Channel == nil || c.Username == "" {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errUnauthorized}
		return
	}

	if msg.Data.Message == "" {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errBadMessage}
		return
	}

	for _, member := range c.Channel.Members {
		member.Send <- &OutgoingMessage{Type: typeNewMessage, Data: &MessageData{Message: msg.Data.Message, Username: c.Username}}
	}
}