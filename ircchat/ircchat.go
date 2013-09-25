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
var FromIrcMessage = make(chan common.Message, 10)
var ircClient *irc.Conn

func Log(msg string) {
	fmt.Printf("IRC: %s\n", msg)
}

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
	ircClient.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { conn.Join(ircChannel) })
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
		l.Args[1],
	}

}

func Connect() (err error) {
	// IRC connect
	Log("Connect irc client")
	if err := ircClient.Connect(config.Get("ajaxchat", "server")); err != nil {
		return err
	}
	return nil
}

func SendToIrc(m common.Message) {
	go func(m common.Message) {
		msg := fmt.Sprintf("%s: %s", m.Username, m.Text)
		Log(fmt.Sprintf("sending: %s", msg))
		ircClient.Privmsg(ircChannel, html.UnescapeString(msg))
	}(m)
}
