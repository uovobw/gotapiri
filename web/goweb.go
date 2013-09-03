package main

import (
	"crypto/tls"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

const config = "config.json"
const sleeptime = 5

var cfg *Config
var ajaxClient *http.Client
var fromAjaxchat = make(chan string)
var lastID = "0"

func printBody(r *http.Response) {
	stuff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Body.Close()
	fmt.Printf("%s\n", stuff)
}

func main() {

	cfg, err := ReadConfigFrom(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	j, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	user := cfg.Get("ajaxchat", "httpuser")
	pass := cfg.Get("ajaxchat", "httppass")

	ajaxClient = &http.Client{
		Jar: j,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy: func(req *http.Request) (*url.URL, error) {
				req.SetBasicAuth(user, pass)
				return nil, nil
			},
		},
	}

	resp, err := ajaxClient.Get(cfg.Get("ajaxchat", "login_url"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loginData := url.Values{
		"login":       {"login"},
		"redirect":    {""},
		"username":    {cfg.Get("ajaxchat", "ajaxuser")},
		"password":    {""},
		"channelName": {cfg.Get("ajaxchat", "ajaxchannel")},
		"lang":        {"en"},
		"submit":      {"Login"},
	}
	resp, err = ajaxClient.PostForm(cfg.Get("ajaxchat", "login_url"), loginData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		for {
			resp, err = ajaxClient.Get(cfg.Get("ajaxchat", "login_url") + "?" + "ajax=true&lastID=" + lastID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			stuff, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			xmlData, err := ParseData(stuff)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, msg := range xmlData.Messages {
				fmt.Printf("%s: %s\n", msg.Username, html.UnescapeString(msg.Text))
				lastID = msg.Id
			}
			//printBody(resp)
			//parseBody(resp)
			time.Sleep(sleeptime * time.Second)
		}
	}()

	<-fromAjaxchat
}
