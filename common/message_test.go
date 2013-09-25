package common

import (
	"fmt"
	"os"
	"testing"
)

func TestParseFromXml(t *testing.T) {
	filename := "test_data/login_fragment.xml"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	xml, err := ParseFromXml(file)
	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}
	if len(xml.Infos) != 5 || len(xml.Users) != 3 || len(xml.Messages) != 2 {
		t.Fatalf("wrong number of pieces in demarshaled structure")
	}
	pass := false
	for _, tstInfo := range xml.Infos {
		if tstInfo.Type == "channelID" &&
			tstInfo.Value == "0" {
			pass = true
		}
	}
	if !pass {
		t.Fatalf("loaded info with wrong parameters")
	}
	for _, tstUser := range xml.Users {
		if tstUser.Nick != "Newfag_tapirc" ||
			tstUser.UserID != "471536092" ||
			tstUser.UserRole != "0" ||
			tstUser.ChannelID != "0" {
			pass = true
		}
	}
	if !pass {
		t.Fatalf("loaded user with wrong parameters")
	}

	for _, tstMessage := range xml.Messages {
		if tstMessage.Id != "212117" ||
			tstMessage.DateTime != "Tue, 03 Sep 2013 01:02:51 +0200" ||
			tstMessage.UserID != "471536092" ||
			tstMessage.UserRole != "0" ||
			tstMessage.ChannelID != "0" {
			pass = true
		}
	}
	if !pass {
		t.Fatalf("loaded message with wrong parameters")
	}
}

func TestHash(t *testing.T) {
	m1 := Message{
		"1",
		"datetime",
		"userid",
		"userrole",
		"channelid",
		"username",
		"text",
	}
	m2 := Message{
		"1",
		"datetime",
		"different_userid",
		"userrole",
		"channelid",
		"username",
		"text",
	}
	m3 := Message{
		"1",
		"datetime",
		"userid",
		"userrole",
		"channelid",
		"different_username",
		"text",
	}
	m4 := Message{
		"1",
		"datetime",
		"userid",
		"userrole",
		"channelid",
		"username",
		"different_text",
	}
	hash1 := m1.Hash()
	hash2 := m2.Hash()
	hash3 := m3.Hash()
	hash4 := m4.Hash()
	if (hash1 != hash2) ||
		(hash2 == hash3) ||
		(hash1 == hash3) ||
		(hash1 == hash4) {
		t.Fatalf("hashing function not working!")
	}

}
