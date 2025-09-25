package main

import (
	"log"
	"net"
)

func runServer() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		errExit(err)
	}

	log.Println("Server started on localhost:42069")
	srv := newServer()

	r := newRoom("Default")
	srv.addRoom(r)

	for {
		conn, err := l.Accept()
		if err != nil {
			errExit(err)
		}

		client := newClient(conn, r, srv)

		r.join <- client

		go client.read()
		go client.write()

	}
}
