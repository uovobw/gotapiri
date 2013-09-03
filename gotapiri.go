package main

import (
	"bufio"
	"flag"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"os"
)

const ircChannel = "##tapiri"

var incoming = make(chan string)

func ReadInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		incoming <- scanner.Text()
	}
}

func main() {
	flag.Parse() // parses the logging flags.
	c := irc.SimpleClient("gotapirc")
	c.EnableStateTracking()
	// Optionally, enable SSL
	c.SSL = true

	// Add handlers to do things here!
	// e.g. join a channel on connect.
	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) { conn.Join(ircChannel) })
	// And a signal on disconnect
	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { quit <- true })
	c.AddHandler("privmsg", func(conn *irc.Conn, line *irc.Line) { fmt.Println(line.Nick + ":" + line.Args[1]) })

	// Tell client to connect
	if err := c.Connect("irc.freenode.net"); err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}

	// start processing input
	go ReadInput()

	// Wait for disconnect
	for {
		select {
		case <-quit:
			os.Exit(0)
		case msg := <-incoming:
			c.Privmsg(ircChannel, msg)
		}
	}
}
