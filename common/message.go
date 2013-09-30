package common

import (
	"crypto/sha256"
	"encoding/xml"
	"io"
	"io/ioutil"
)

// Type XmlData is used to unmarshal ajaxchat responses from xml
type XmlData struct {
	Infos    []Info    `xml:"infos>info"`
	Messages []Message `xml:"messages>message"`
	Users    []User    `xml:"users>user"`
}

type Info struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type Message struct {
	Id        string `xml:"id,attr"`
	DateTime  string `xml:"dateTime,attr"`
	UserID    string `xml:"userID,attr"`
	UserRole  string `xml:"userRole,attr"`
	ChannelID string `xml:"channelID,attr"`
	Username  string `xml:"username"`
	Text      string `xml:"text"`
}

type User struct {
	Nick      string `xml:",chardata"`
	UserID    string `xml:"userID,attr"`
	UserRole  string `xml:"userRole,attr"`
	ChannelID string `xml:"channelID,attr"`
}

// Function ToString returns a human readable representation
// of a given common.Message struct in the form
// username: text of the message
func (m Message) ToString() (s string) {
	return m.Username + ": " + m.Text
}

// Function Hash implements the Hashable interface, used in
// object sets
func (m *Message) Hash() (hash string) {
	h := sha256.New()
	io.WriteString(h, m.Id+m.DateTime+m.Username+m.Text)
	return string(h.Sum(nil))
}

// Function ParseFromXml receives a binary ReadCloser and
// expects to unmarshal all its contents in an XmlData structure,
// returning error on failures
func ParseFromXml(source io.ReadCloser) (v *XmlData, e error) {
	data, err := ioutil.ReadAll(source)
	if err != nil {
		return nil, e
	}
	defer source.Close()
	v = &XmlData{}
	err = xml.Unmarshal([]byte(data), v)
	if err != nil {
		return nil, e
	}
	//	fmt.Printf("%+v\n", v)
	return v, nil
}
