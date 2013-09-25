package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/uovobw/gotapiri/ajaxchat"
	"github.com/uovobw/gotapiri/ircchat"
	"os"
)

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

func main() {
	flag.Parse() // parses the logging flags.

	// WEBCHAT INIT
	Log("Create webchat client")
	err := ajaxchat.Init()
	if err != nil {
		Log(fmt.Sprintf("Cannot create webclient: %s", err))
		os.Exit(1)
	}

	// IRCCHAT INIT
	if err = ircchat.Init(); err != nil {
		Log(fmt.Sprintf("Cannot create irc client: %s", err))
		os.Exit(1)
	}
	// IRCCHAT connect
	if err = ircchat.Connect(); err != nil {
		Log(fmt.Sprintf("Cannot connect irc client: %s", err))
		os.Exit(1)
	}

	// WEBCHAT update
	Log("Connect webchat client")
	go ajaxchat.UpdateLoop()

	// start processing input
	//go ReadInput()

	// Wait for disconnect
	Log("In message loop")
	for {
		select {
		case xmlData := <-ajaxchat.FromAjaxMessage:
			for _, msg := range xmlData.Messages {
				ircchat.SendToIrc(msg)
			}
		case ircMessage := <-ircchat.FromIrcMessage:
			go ajaxchat.SendToAjaxchat(ircMessage)
		}
	}
}
