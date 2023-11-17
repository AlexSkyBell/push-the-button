package server

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type server struct {
	rooms   *rooms
	clients sync.Map

	addCh  chan *client
	delCh  chan *client
	doneCh chan bool
	errCh  chan error
}

func NewServer() *server {
	return &server{
		rooms:   &rooms{},
		clients: sync.Map{},
		addCh:   make(chan *client),
		delCh:   make(chan *client),
		doneCh:  make(chan bool),
		errCh:   make(chan error),
	}
}

func (s *server) Listen() {
	log.Println("Starting ws server")

	onConnected := func(ws *websocket.Conn) {
		wsClient := NewClient(ws, s.errCh, s.delCh)
		s.addCh <- wsClient
		wsClient.Listen()

		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()
	}

	basePath := "/"
	http.Handle(basePath, websocket.Handler(onConnected))

	for {
		select {
		case client := <-s.addCh:
			log.Printf("client connected %d\n", client.id)
			s.addClient(client)
		case client := <-s.delCh:
			log.Printf("client disconnected %d\n", client.id)
			s.removeClient(client)
		case err := <-s.errCh:
			log.Printf("error occurred: %s\n", err.Error())
		case <-s.doneCh:
			return
		}
	}
}

func (s *server) addClient(client *client) {
	s.clients.Store(&client.id, client)
}

func (s *server) removeClient(client *client) {
	s.clients.Delete(&client.id)
}
