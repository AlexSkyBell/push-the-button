package server

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

var clientId int

type client struct {
	id        int
	name      string
	ws        *websocket.Conn
	ch        chan *Message
	doneRead  chan struct{}
	doneWrite chan struct{}
	errCh     chan error
	serverCh  chan *client
}

func NewClient(ws *websocket.Conn, errCh chan error, serverCh chan *client) *client {
	clientId++
	return &client{
		id:        clientId,
		name:      "Client",
		ws:        ws,
		ch:        make(chan *Message),
		doneRead:  make(chan struct{}),
		doneWrite: make(chan struct{}),
		errCh:     errCh,
		serverCh:  serverCh,
	}
}

func (c *client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *client) write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		fmt.Println("do nothing")
		c.Close()
	}
}

func (c *client) listenWrite() {
	for {
		select {
		case msg := <-c.ch:
			websocket.JSON.Send(c.ws, msg)
		case <-c.doneWrite:
			fmt.Printf("client write done channel %d\n", c.id)
			return
		}
	}
}

func (c *client) Close() {
	close(c.ch)
	close(c.doneRead)
	close(c.doneWrite)
	c.serverCh <- c
}

type Message struct {
	Command string `json:"command"`
	Payload string `json:"payload"`
}

func (c *client) listenRead() {
	for {
		select {
		case <-c.doneRead:
			fmt.Printf("client read done channel %d\n", c.id)
			return

		default:
			fmt.Printf("client default listen read %d\n", c.id)
			msg := &Message{}
			err := websocket.JSON.Receive(c.ws, msg)
			if err == io.EOF {
				fmt.Printf("client wants to disconnect %d\n", c.id)
				c.Close()
			} else if err != nil {
				fmt.Printf("error occurred %d %s\n", c.id, err)
				c.errCh <- err
			} else {
				fmt.Printf("process message %d\n", c.id)
				c.processMessage(msg)
			}
		}
	}
}

func (c *client) processMessage(msg *Message) {
	c.write(msg)
}
