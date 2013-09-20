package common

import (
	"testing"
)

func TestReadConfigFrom(t *testing.T) {
	filename := "test_data/config_test.json"
	config, err := ReadConfigFrom(filename)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}
	if len(config) != 3 {
		t.Fatalf("expected 3 elements, got %s", len(config))
	}
	if len(config["ajaxchat"]) != 4 {
		t.Fatalf("expected 4 elements, got %s", len(config["ajaxchat"]))
	}
	if len(config["general"]) != 1 {
		t.Fatalf("expected 1 elements, got %s", len(config["general"]))
	}
	if len(config["tumbl"]) != 4 {
		t.Fatalf("expected 4 elements, got %s", len(config["tumbl"]))
	}
}

func TestGet(t *testing.T) {
	c := Config{
		"S1": map[string]string{
			"k11": "v11",
		},
		"S2": map[string]string{
			"k21": "v21",
		},
	}
	if c.Get("S1", "k11") != "v11" {
		t.Fatalf("got the wrong key")
	}
	if c.Get("S2", "k21") != "v21" {
		t.Fatalf("got the wrong key")
	}
	defer func() {
		if r := recover(); r != nil {
		} else {
			t.Fatalf("did not raise!")
		}
	}()
	c.Get("S2", "k11")
}
