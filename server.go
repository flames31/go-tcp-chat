package main

type server struct {
	rooms map[*room]bool
}

func newServer() *server {
	return &server{
		rooms: map[*room]bool{},
	}
}

func (s *server) addRoom(r *room) {
	s.rooms[r] = true
}
