package tumblr

import (
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	"github.com/uovobw/gotapiri/common"
)

var (
	client = new(gotumblr.TumblrRestClient)
	config common.Config
)

const configurationFilename = "config.json"

// Log Messages for Tumblr pkg
func Log(msg string) {
	fmt.Printf("TUMBLR: %s\n", msg)
}

//Init for Tumblr pkg
func Init() (err error) {
	config, err = common.ReadConfigFrom(configurationFilename)
	if err != nil {
		return err
	}

	Log("Create tumblr client")

	appKey := config.Get("tumblr", "app_key")
	appSecret := config.Get("tumblr", "app_secret")
	oauthToken := config.Get("tumblr", "oauth_token")
	oauthTokenSecret := config.Get("tumblr", "oauth_token_secret")
	callbackURL := config.Get("tumblr", "url")

	client = gotumblr.NewTumblrRestClient(appKey, appSecret, oauthToken, oauthTokenSecret, callbackURL, "http://api.tumblr.com")

	return
}
