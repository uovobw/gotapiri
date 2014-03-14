package twitter

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/uovobw/gotapiri/common"
	"strings"
)

var (
	api        = new(anaconda.TwitterApi)
	config     common.Config
	twitterTag string
	lastTweet  string
)

const configurationFilename = "config.json"

// Log Messages for Twitter pkg
func Log(msg string) {
	fmt.Printf("TWITTER: %s\n", msg)
}

func Init() (err error) {
	config, err = common.ReadConfigFrom(configurationFilename)
	if err != nil {
		return err
	}

	twitterTag = config.Get("twitter", "twitter_tag")

	Log("Create twitter client")

	appKey := config.Get("twitter", "app_key")
	appSecret := config.Get("twitter", "app_secret")
	oauthToken := config.Get("twitter", "oauth_token")
	oauthTokenSecret := config.Get("twitter", "oauth_token_secret")

	anaconda.SetConsumerKey(appKey)
	anaconda.SetConsumerSecret(appSecret)

	api = anaconda.NewTwitterApi(oauthToken, oauthTokenSecret)

	return
}

func PostTweet(status common.Message) (err error) {
	msg := status.Text
	user := status.Username
	// TODO: fix deduplication of tweets
	if msg == lastTweet {
		return nil
	}
	lastTweet = msg
	if strings.Contains(msg, twitterTag) {
		msg = strings.Replace(msg, twitterTag, "", -1)
		msg = strings.Replace(msg, user, "", -1)
		if len(msg) > 140 {
			Log("Trimming tweet")
			msg = msg[:140]
		}
		Log(fmt.Sprintf("Posting tweet: %s", msg))
		_, err = api.PostTweet(msg, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
