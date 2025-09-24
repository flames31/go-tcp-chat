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

	defaultRoom := newRoom("Default")
	go defaultRoom.start()

	for {
		conn, err := l.Accept()
		if err != nil {
			errExit(err)
		}

		client := newClient(conn, defaultRoom)

		defaultRoom.join <- client

		go client.read()
		go client.write()

	}
}
