package main

import (
	"bufio"
	"flag"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/uovobw/gotapiri/ajaxchat"
	"os"
)

const ircChannel = "##tapiri"

var incoming = make(chan string, 10)

func Log(msg string) {
	fmt.Printf("MAIN: %s\n", msg)
}

func ReadInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		incoming <- scanner.Text()
	}
}

func SendToIrc(d *ajaxchat.XmlData) {
	go func() {
		for _, msg := range d.Messages {
			incoming <- fmt.Sprintf("%s: %s", msg.Username, msg.Text)
		}
	}()
}

func main() {
	flag.Parse() // parses the logging flags.

	// IRC INIT
	Log("Create irc client")
	c := irc.SimpleClient("gotapirc")
	c.EnableStateTracking()
	c.SSL = true

	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) { conn.Join(ircChannel) })
	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { quit <- true })
	c.AddHandler("privmsg", func(conn *irc.Conn, line *irc.Line) { fmt.Println(line.Nick + ":" + line.Args[1]) })

	// WEBCHAT INIT
	Log("Create webchat client")
	err := ajaxchat.Init()
	if err != nil {
		fmt.Printf("Cannot create webclient: %s\n", err)
		os.Exit(1)
	}

	// IRC connect
	Log("Connect irc client")
	if err := c.Connect("irc.freenode.net"); err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}

	// WEBCHAT update
	Log("Connect webchat client")
	go ajaxchat.UpdateLoop()

	// start processing input
	//go ReadInput()

	// Wait for disconnect
	Log("In incoming message loop")
	for {
		select {
		case <-quit:
			os.Exit(0)
		case msg := <-incoming:
			c.Privmsg(ircChannel, msg)
		case xmlData := <-ajaxchat.FromAjaxResult:
			SendToIrc(xmlData)
		}
	}
}
