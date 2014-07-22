package common

import (
	"crypto/sha256"
	"encoding/xml"
	"io"
	"io/ioutil"
)

// XMLData is used to unmarshal ajaxchat responses from xml
type XMLData struct {
	Infos    []Info    `xml:"infos>info"`
	Messages []Message `xml:"messages>message"`
	Users    []User    `xml:"users>user"`
}

// Info represents cdata fields from the ajaxchat response
type Info struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// Message struct unmasharls the ajaxchat response from the chat
type Message struct {
	ID        string `xml:"id,attr"`
	DateTime  string `xml:"dateTime,attr"`
	UserID    string `xml:"userID,attr"`
	UserRole  string `xml:"userRole,attr"`
	ChannelID string `xml:"channelID,attr"`
	Username  string `xml:"username"`
	Text      string `xml:"text"`
}

// User contains only the data that is useful for the bot, it is missing
// some internally meaningful ids for the ajaxchat. Future versions will need
// to use them in order to better simulate the client
type User struct {
	Nick      string `xml:",chardata"`
	UserID    string `xml:"userID,attr"`
	UserRole  string `xml:"userRole,attr"`
	ChannelID string `xml:"channelID,attr"`
}

// ToString returns a human readable representation
// of a given common.Message struct in the form
// username: text of the message
func (m Message) ToString() (s string) {
	return m.Username + ": " + m.Text
}

// Hash implements the Hashable interface, used in
// object sets
func (m *Message) Hash() (hash string) {
	h := sha256.New()
	io.WriteString(h, m.ID+m.DateTime+m.Username+m.Text)
	return string(h.Sum(nil))
}

// ParseFromXML receives a binary ReadCloser and
// expects to unmarshal all its contents in an XMLData structure,
// returning error on failures
func ParseFromXML(source io.ReadCloser) (v *XMLData, e error) {
	data, err := ioutil.ReadAll(source)
	if err != nil {
		return nil, e
	}
	defer source.Close()
	v = &XMLData{}
	err = xml.Unmarshal([]byte(data), v)
	if err != nil {
		return nil, e
	}
	//	fmt.Printf("%+v\n", v)
	return v, nil
}
