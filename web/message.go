package main

import (
	"encoding/xml"
)

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

func ParseData(data []byte) (v *XmlData, e error) {
	v = &XmlData{}
	err := xml.Unmarshal([]byte(data), v)
	if err != nil {
		return nil, e
	}
	//	fmt.Printf("%+v\n", v)
	return v, nil
}
