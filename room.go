package main

import "log"

type room struct {
	name    string
	forward chan msg
	join    chan *client
	leave   chan *client
	members map[*client]bool
}

type msg struct {
	from string
	data []byte
}

func newRoom(name string) *room {
	r := &room{
		name:    name,
		forward: make(chan msg),
		join:    make(chan *client),
		leave:   make(chan *client),
		members: make(map[*client]bool),
	}

	go r.start()

	return r
}

func (r *room) start() {
	log.Printf("room go routine started for %v", r.name)
	for {
		select {
		case client := <-r.join:
			r.members[client] = true
		case client := <-r.leave:
			delete(r.members, client)
			close(client.send)
			client.exit <- struct{}{}
		case msg := <-r.forward:
			for m := range r.members {
				m.send <- msg
			}
		}
	}
}
