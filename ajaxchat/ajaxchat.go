// Package ajaxchat abstracts the details for connecting to blueimp's AjaxChat
// (homepage: http://frug.github.io/AJAX-Chat/ )
package ajaxchat

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"github.com/uovobw/gotapiri/common"
)

const configurationFilename = "config.json"
const sleeptime = 5

var ajaxClient *http.Client
var lastID = "0"
var config common.Config

// Fromajaxmessage returns the messages that are returned
// from the webchat.
var FromAjaxMessage = make(chan *common.XmlData, 10)

// Log temporary function to abstract away the log that needs
// to be around the code
func Log(msg string) {
	fmt.Printf("AC: %s\n", msg)
}

// printBody internal debugging function used to print
// the content of an http.Response to screen for testing
func printBody(r *http.Response) {
	stuff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Body.Close()
	fmt.Printf("%s\n", stuff)
}

// UpdateLoop periodically polls the webchat and fetches the new
// messages, unmarshals them from xml to a common.XmlData object and
// publishes them via the Fromajaxmessage channel
func UpdateLoop() {
	Log("Running update loop")
	for {
		time.Sleep(sleeptime * time.Second)
		resp, err := ajaxClient.Get(config.Get("ajaxchat", "login_url") + "?" + "ajax=true&lastID=" + lastID)
		if err != nil {
			fmt.Printf("error getting update from chat: %s\n", err)
			continue
		}
		xmlData, err := common.ParseFromXml(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("error in parsing data: %s\n", err)
			continue
		}
		FromAjaxMessage <- xmlData
		for _, msg := range xmlData.Messages {
			lastID = msg.Id
		}
		//printBody(resp)
	}
}

// createClient initializes the http.Client used throughout the
// module by setting the username/password BasicAuthentication
// for each request, ignoring self-signed SSL certificates and
// using the default in-memory cookiejar
func createClient(user, pass string) (err error) {
	Log("Creating client")
	j, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

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

	return nil
}

func readConfiguration(configfilename string) (c common.Config, err error) {
	Log("Reading configuraton")
	cfg, err := common.ReadConfigFrom(configfilename)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Init performs the operations required to authenticate to the
// webchat as a registered client and returns an error should any
// http errors occur. MUST be called first on a non initialized client
func init() {
	Log("In init")
	var err error
	config, err = readConfiguration(configurationFilename)
	if err != nil {
		panic(fmt.Errorf("coulf not read configuration: %s\n", err))
	}

	user := config.Get("ajaxchat", "httpuser")
	pass := config.Get("ajaxchat", "httppass")

	err = createClient(user, pass)
	if err != nil {
		panic(fmt.Errorf("could not create web client: %s\n", err))
	}
	// first get to init the state on the remote end
	Log("Login (1/2)")
	_, err = ajaxClient.Get(config.Get("ajaxchat", "login_url"))
	if err != nil {
		panic(fmt.Errorf("could not reach login page: %s\n", err))
	}

	loginData := url.Values{
		"login":       {"login"},
		"redirect":    {""},
		"userName":    {config.Get("ajaxchat", "ajaxuser")},
		"password":    {config.Get("ajaxchat", "ajaxpass")},
		"channelName": {config.Get("ajaxchat", "ajaxchannel")},
		"lang":        {"en"},
		"submit":      {"Login"},
	}

	Log("Login (2/2)")
	_, err = ajaxClient.PostForm(config.Get("ajaxchat", "login_url"), loginData)
	if err != nil {
		panic(fmt.Errorf("could not finalize login: %s\n", err))
	}
}

// SendToAjaxchat sends a message to the webchat, must be called only
// on an already Init-ed client of the call will fail
func SendToAjaxchat(msg common.Message) (err error) {
	postData := url.Values{
		"ajax":   {"true"},
		"text":   {msg.Username + ": " + msg.Text},
		"lastID": {lastID},
	}
	_, err = ajaxClient.PostForm(config.Get("ajaxchat", "msg_url"), postData)
	if err != nil {
		return fmt.Errorf("could not post message: %s\n", err)
	}
	Log(fmt.Sprintf("sending: %s", postData["text"]))
	return nil
}
