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
