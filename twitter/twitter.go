package twitter

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/uovobw/gotapiri/common"
	"strings"
)

var api anaconda.TwitterApi
var config common.Config
var twitter_tag string

const configurationFilename = "config.json"

func Log(msg string) {
	fmt.Printf("TWITTER: %s\n", msg)
}

func Init() (err error) {
	config, err = common.ReadConfigFrom(configurationFilename)
	if err != nil {
		return err
	}

	twitter_tag = config.Get("twitter", "twitter_tag")

	Log("Create twitter client")

	appKey := config.Get("twitter", "app_key")
	appSecret := config.Get("twitter", "app_secret")
	oauthToken := config.Get("twitter", "oauth_token")
	oauthTokenSecret := config.Get("twitter", "oauth_token_secret")

	anaconda.SetConsumerKey(appKey)
	anaconda.SetConsumerSecret(appSecret)

	api = anaconda.NewTwitterApi(oauthToken, oauthTokenSecret)

	return nil
}

func PostTweet(status string) (err error) {
	if strings.Contains(status, twitter_tag) {
		if len(status) > 140 {
			Log("Trimming tweet")
			status = status[:140]
		}
		Log(fmt.Sprintf("Posting tweet: %s", status))
		_, err = api.PostTweet(status, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
