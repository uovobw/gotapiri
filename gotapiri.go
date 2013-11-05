// An IRC-TO-WEBCHAT partial, buggy and horrribly-written transport. Relays messages
// from a BlueImp's AjaxChat installation to any IRC channel.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/uovobw/gotapiri/ajaxchat"
	"github.com/uovobw/gotapiri/ircchat"
	"github.com/uovobw/gotapiri/twitter"
	"html"
	"os"
	"strings"
)

var incoming = make(chan string, 10)
var seenMessages = make(map[string]bool)
var toTwitter = make(chan string, 2)

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

	// TWITTER INIT
	if err = twitter.Init(); err != nil {
		Log(fmt.Sprintf("Cannot create twitter client: %s", err))
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

	g_PostTweet()

	// Wait for disconnect
	Log("In message loop")
	for {
		select {
		case xmlData := <-ajaxchat.FromAjaxMessage:
			for _, msg := range xmlData.Messages {
				seen := false
				seenmsg := ""
				for el := range seenMessages {
					if strings.HasSuffix(clean(msg.Text), clean(el)) {
						seen = true
						seenmsg = clean(el)
						break
					}
				}
				if !seen {
					toTwitter <- msg.Text
					ircchat.SendToIrc(msg)
					delete(seenMessages, seenmsg)
				}
			}
		case ircMessage := <-ircchat.FromIrcMessage:
			go ajaxchat.SendToAjaxchat(ircMessage)
			seenMessages[clean(ircMessage.Text)] = true
		}
	}
}

func g_PostTweet() {
	// if posting the tweet fails we silently ignore
	Log("Running twitter goroutine")
	go func() {
		t := <-toTwitter
		twitter.PostTweet(t)
	}()
}

func clean(s string) (r string) {
	return html.UnescapeString(s)
}
