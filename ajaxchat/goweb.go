package ajaxchat

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/uovobw/gotapiri/common"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

const configurationFilename = "config.json"
const sleeptime = 5

var ajaxClient *http.Client
var lastID = "0"
var config common.Config

var FromAjaxResult = make(chan *XmlData, 10)

func Log(msg string) {
	fmt.Printf("AC: %s\n", msg)
}

func printBody(r *http.Response) {
	stuff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Body.Close()
	fmt.Printf("%s\n", stuff)
}

func UpdateLoop() {
	Log("Running update loop")
	for {
		resp, err := ajaxClient.Get(config.Get("ajaxchat", "login_url") + "?" + "ajax=true&lastID=" + lastID)
		if err != nil {
			fmt.Printf("error getting update from chat: %s\n", err)
		}
		xmlData, err := ParseFromXml(resp.Body)
		if err != nil {
			fmt.Printf("error in parsing data: %s\n", err)
		}
		fmt.Printf("got xmldata as: %+v\n", xmlData)
		FromAjaxResult <- xmlData
		//printBody(resp)
		time.Sleep(sleeptime * time.Second)
	}
}

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

func Init() (err error) {
	Log("In init")

	config, err = readConfiguration(configurationFilename)
	if err != nil {
		return errors.New(fmt.Sprintf("coulf not read configuration: %s\n", err))
	}

	user := config.Get("ajaxchat", "httpuser")
	pass := config.Get("ajaxchat", "httppass")

	err = createClient(user, pass)
	if err != nil {
		return errors.New(fmt.Sprintf("could not create web client: %s\n", err))
	}
	// first get to init the state on the remote end
	Log("Login (1/2)")
	_, err = ajaxClient.Get(config.Get("ajaxchat", "login_url"))
	if err != nil {
		return errors.New(fmt.Sprintf("could not reach login page: %s\n", err))
	}

	loginData := url.Values{
		"login":       {"login"},
		"redirect":    {""},
		"username":    {config.Get("ajaxchat", "ajaxuser")},
		"password":    {""},
		"channelName": {config.Get("ajaxchat", "ajaxchannel")},
		"lang":        {"en"},
		"submit":      {"Login"},
	}

	Log("Login (2/2)")
	_, err = ajaxClient.PostForm(config.Get("ajaxchat", "login_url"), loginData)
	if err != nil {
		return errors.New(fmt.Sprintf("could not finalize login: %s\n", err))
	}
	return nil
}
