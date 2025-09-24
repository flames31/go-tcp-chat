package main

import (
	"fmt"
	"log"
	"net"
)

type client struct {
	conn net.Conn
	name string
	room *room
	send chan msg
	exit chan struct{}
}

func newClient(conn net.Conn, room *room) *client {
	return &client{
		conn: conn,
		name: "anonymous",
		room: room,
		send: make(chan msg),
		exit: make(chan struct{}),
	}
}

func (c *client) read() {
	defer c.conn.Close()
	buf := make([]byte, 256)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			log.Printf("ERR : %v", err)
			c.room.leave <- c
			defer c.conn.Close()
			break
		}

		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			if string(data) == "/quit\n" {
				c.room.leave <- c
				break
			} else {
				m := append([]byte(fmt.Sprintf("%v: ", c.name)), data...)
				c.room.forward <- msg{
					from: c.conn.RemoteAddr().String(),
					data: m,
				}
			}
		}
	}
	log.Printf("%v has left the chat", c)
}

func (c *client) write() {
	defer c.conn.Close()

	for msg := range c.send {
		if c.conn.RemoteAddr().String() != msg.from {
			c.conn.Write(msg.data)
		}
	}
}
