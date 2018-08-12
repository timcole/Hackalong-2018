package main

import (
	"math/rand"
	"sync"
	"time"
)

// Channel is the chatroom
type Channel struct {
	Topic        string         `json:"topic"`
	Duration     time.Duration  `json:"duration"`
	Members      []*Client      `json:"members"`
	History      []*MessageData `json:"history"`
	historyMutex sync.Mutex
	Slot         int `json:"slots"`
	mutex        sync.Mutex
}

func (c *Client) createChannel(msg ReceiveMessage) {
	if c.Channel != nil || c.Username == "" {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errUnauthorized, Data: &MessageData{Timestamp: time.Now()}}
		return
	}

	channel := &Channel{
		Topic:    msg.Data.Topic,
		Duration: time.Second,
		Members:  []*Client{c},
		Slot:     2,
	}

	c.Connections.chanMutex.Lock()
	c.Connections.Channels = append(c.Connections.Channels, channel)
	c.Connections.chanMutex.Unlock()
	c.Channel = channel

	go func() {
		var dur = time.NewTicker(time.Second)
		for range dur.C {
			channel.Duration = channel.Duration + time.Second
		}
	}()

	c.Vote = 1

	c.Send <- &OutgoingMessage{Type: typeJoinChannel, Data: &MessageData{Topic: channel.Topic, Members: channel.Members, Timestamp: time.Now()}}
}

func (c *Client) joinChannel() {
	if c.Username == "" {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errUnauthorized, Data: &MessageData{Timestamp: time.Now()}}
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
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errNoChannels, Data: &MessageData{Timestamp: time.Now()}}
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

	c.Vote = 0
	c.Send <- &OutgoingMessage{Type: typeJoinChannel, Data: &MessageData{Topic: join.Topic, Members: join.Members, Timestamp: time.Now()}}

	for _, msg := range c.Channel.History {
		c.Send <- &OutgoingMessage{Type: typeNewMessage, Data: msg}
	}

	for _, member := range c.Channel.Members {
		member.Send <- &OutgoingMessage{Type: typeMemberJoin, Data: &MessageData{Username: c.Username, Timestamp: time.Now()}}
	}
}

func (c *Client) leaveChannel() {
	if c.Channel == nil {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errNoChannels, Data: &MessageData{Timestamp: time.Now()}}
		return
	}

	c.Channel.mutex.Lock()
	var tmpMembers []*Client
	var members = c.Channel.Members
	for i := range members {
		if members[i] == c {
			tmpMembers = append(members[:i], members[i+1:]...)
			break
		}
	}
	c.Channel.Members = tmpMembers
	c.Channel.Slot++

	for _, member := range c.Channel.Members {
		member.Send <- &OutgoingMessage{Type: typeMemberLeave, Data: &MessageData{Username: c.Username, Timestamp: time.Now()}}
	}

	if len(c.Channel.Members) == 0 && c.Channel.Slot == 3 {
		c.Channel.Slot = 0
		go c.Connections.close(c.Channel)
	}

	c.Channel.mutex.Unlock()
	c.Channel = nil

	c.Send <- &OutgoingMessage{Type: typeLeaveChannel, Data: &MessageData{Timestamp: time.Now()}}
}

func (conns *Connections) close(channel *Channel) {
	conns.chanMutex.Lock()
	var tmpChannels []*Channel
	var channels = conns.Channels
	for i := range channels {
		if channels[i] == channel {
			tmpChannels = append(channels[:i], channels[i+1:]...)
			break
		}
	}
	conns.Channels = tmpChannels
	conns.chanMutex.Unlock()
}

func (c *Client) sendMessage(msg ReceiveMessage) {
	if c.Channel == nil || c.Username == "" {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errUnauthorized, Data: &MessageData{Timestamp: time.Now()}}
		return
	}

	if msg.Data.Message == "" {
		c.Send <- &OutgoingMessage{Type: typeResponse, Error: errBadMessage, Data: &MessageData{Timestamp: time.Now()}}
		return
	}

	var storageMsg = &MessageData{Message: msg.Data.Message, Timestamp: time.Now(), Sender: c}
	c.Channel.historyMutex.Lock()
	c.Channel.History = append(c.Channel.History, storageMsg)
	c.Channel.historyMutex.Unlock()

	if len(c.Channel.History)%20 == 0 && c.Channel.Slot == 0 {
		c.Channel.mutex.Lock()
		c.Channel.Slot = c.Channel.Slot + 2
		c.Channel.mutex.Unlock()
	}

	for _, member := range c.Channel.Members {
		member.Send <- &OutgoingMessage{Type: typeNewMessage, Data: storageMsg}
	}
}
