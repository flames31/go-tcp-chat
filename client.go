package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	conn net.Conn
	name string
	room *room
	srv  *server
	send chan msg
	exit chan struct{}
}

func newClient(conn net.Conn, room *room, srv *server) *client {
	return &client{
		conn: conn,
		name: "anonymous",
		room: room,
		srv:  srv,
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

			c.parseMsg(data)
		}
	}
	log.Printf("%v has left the chat.\n", c)
}

func (c *client) write() {
	defer c.conn.Close()

	for msg := range c.send {
		if c.conn.RemoteAddr().String() != msg.from {
			c.conn.Write(msg.data)
		}
	}
}

func (c *client) sendTowrite(data []byte) {
	c.room.forward <- msg{
		from: c.conn.RemoteAddr().String(),
		data: data,
	}
}

func (c *client) sendToUser(data []byte) {
	c.conn.Write(data)
}

func (c *client) parseMsg(m []byte) {
	msgText, err := bufio.NewReader(bytes.NewReader(m)).ReadString('\n')
	if err != nil {
		return
	}

	msgText = strings.Trim(msgText, "\r\n")

	args := strings.Split(msgText, " ")
	first := strings.TrimSpace(args[0])

	if first[0] == '/' {
		c.execCmd(first, args[1:])
	} else {
		m = append([]byte(fmt.Sprintf("%v: ", c.name)), m...)
		c.sendTowrite(m)
	}
}

func (c *client) execCmd(cmd string, args []string) {
	switch cmd {
	case "/name":
		b := []byte{}
		if len(args) != 1 {
			b = []byte("ERR : /name should have exactly one argument\n")
		} else {
			c.name = args[0]
			b = []byte("Username changed to " + c.name + "\n")
		}
		c.sendToUser(b)
	case "/quit":
		c.sendTowrite([]byte(fmt.Sprintf("%v has left the chat.\n", c.name)))
		c.room.leave <- c
	case "/users":
		b := []byte{}
		for m := range c.room.members {
			b = append(b, []byte(m.name)...)
			b = append(b, []byte("\n")...)
		}
		c.sendToUser(b)
	case "/join":
		c.sendToUser([]byte("To be implemented...\n"))
	case "/rooms":
		b := []byte{}
		for r := range c.srv.rooms {
			b = append(b, []byte(r.name)...)
			b = append(b, []byte("\n")...)
		}
		c.sendToUser(b)
	default:
		c.sendToUser([]byte("ERR : unknown command! /cmd to list all available commands\n"))
	}
}
