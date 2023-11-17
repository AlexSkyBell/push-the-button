package server

import "sync"

type rooms struct {
	mu    sync.Mutex
	rooms map[int]*room

	clients map[int]*room
}

func NewRooms() *rooms {
	return &rooms{}
}

func (r *rooms) getRoom(roomID int) *room {
	r.mu.Lock()
	defer r.mu.Unlock()
	rom, ok := r.rooms[roomID]
	if !ok {
		rom = NewRoom(roomID)
	}
	return rom
}

func (r *rooms) AddClient(client *client, roomID int) {
	room := r.getRoom(roomID)
	room.Add(client)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client.id] = room
}

func (r *rooms) RemoveClient(client *client) {
	if room, ok := r.clients[client.id]; ok {
		r.mu.Lock()
		defer r.mu.Unlock()
		room.Remove(client)
		delete(r.clients, client.id)
	}
}

type room struct {
	id      int
	clients sync.Map
}

func NewRoom(roomID int) *room {
	return &room{id: roomID}
}

func (r *room) Add(client *client) {
	r.clients.Store(&client.id, client)
}

func (r *room) Remove(client *client) {
	r.clients.Delete(&client.id)
}
