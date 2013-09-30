// Package IrcChat abstracts the creation, initialization and use of
// an irc client that connects to a server, joins a definite channel
// and waits for messages to be sent or received from it
package ircchat

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/uovobw/gotapiri/common"
	"html"
)

const ircChannel = "##tapiri"
const configurationFilename = "config.json"

var config common.Config

// FromIrcMessage used to publish incoming messages (from irc)
// that are seen from the client
var FromIrcMessage = make(chan common.Message, 10)
var ircClient *irc.Conn

// Log temporary function that abstracts the logging
// that needs to be around in the code
func Log(msg string) {
	fmt.Printf("IRC: %s\n", msg)
}

// function Init must be called first on a non-initialized irc
// client. It will connect with SSL to a given irc server using a given
// username, join a channel and wait for messages
func Init() (err error) {
	config, err = common.ReadConfigFrom(configurationFilename)
	if err != nil {
		return err
	}

	// IRC INIT
	Log("Create irc client")
	ircClient = irc.SimpleClient(config.Get("ajaxchat", "user"))
	ircClient.EnableStateTracking()
	ircClient.SSL = true

	ircClient.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) { conn.Join(ircChannel) })
	ircClient.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { Connect() })
	ircClient.AddHandler("privmsg", func(conn *irc.Conn, line *irc.Line) { FromIrcMessage <- createMessageFromIrc(line) })

	return nil
}

func createMessageFromIrc(l *irc.Line) common.Message {
	return common.Message{
		"",
		"",
		"",
		"",
		"",
		l.Nick,
		html.UnescapeString(l.Args[1]),
	}

}

// Function Connect connect the client to a server that is to be specified
// in the configuration file
func Connect() (err error) {
	// IRC connect
	Log("Connect irc client")
	if err := ircClient.Connect(config.Get("ajaxchat", "server")); err != nil {
		return err
	}
	return nil
}

// Function SendToIrc sends the message passed to it
// to the currently connected irc server in the currently
// joined channel as the configured user
func SendToIrc(m common.Message) {
	go func(m common.Message) {
		msg := fmt.Sprintf("%s: %s", m.Username, m.Text)
		Log(fmt.Sprintf("sending: %s", msg))
		ircClient.Privmsg(ircChannel, html.UnescapeString(msg))
	}(m)
}
