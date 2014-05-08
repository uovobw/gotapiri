// An IRC-TO-WEBCHAT partial, buggy and horrribly-written transport. Relays messages
// from a BlueImp's AjaxChat installation to any IRC channel.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/uovobw/gotapiri/ajaxchat"
	"github.com/uovobw/gotapiri/common"
	"github.com/uovobw/gotapiri/ircchat"
	"github.com/uovobw/gotapiri/tumblr"
	"github.com/uovobw/gotapiri/twitter"
)

var incoming = make(chan string, 10)
var seenMessages = make(map[string]bool)
var toTwitter = make(chan string, 10)

//Log for main package
func Log(msg string) {
	fmt.Printf("MAIN: %s\n", msg)
}

//ReadInput and send to the incoming messages
func ReadInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		incoming <- scanner.Text()
	}
}

func main() {
	flag.Parse() // parses the logging flags.

	// IRCCHAT connect
	if err := ircchat.Connect(); err != nil {
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
					go gPostTweet(msg)
					go gPostImage(msg)
					ircchat.SendToIrc(msg)
					delete(seenMessages, seenmsg)
				}
			}
		case ircMessage := <-ircchat.FromIrcMessage:
			go ajaxchat.SendToAjaxchat(ircMessage)
			go gPostTweet(ircMessage)
			go gPostImage(ircMessage)
			seenMessages[clean(ircMessage.Text)] = true
		}
	}
}

func gPostTweet(msg common.Message) {
	// if posting the tweet fails we silently ignore
	err := twitter.PostTweet(msg)
	if err != nil {
		Log(fmt.Sprintf("Error posting tweet! %s", err))
	}
}

func gPostImage(msg common.Message) {
	// if posting fails we silently ignore
	err := tumblr.PostImage(msg)
	if err != nil {
		Log(fmt.Sprintf("Error posting image! %s", err))
	}
}

func clean(s string) (r string) {
	return html.UnescapeString(s)
}
