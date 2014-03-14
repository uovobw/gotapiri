package tumblr

import (
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	"github.com/uovobw/gotapiri/common"
	"regexp"
	"strings"
)

var (
	client = new(gotumblr.TumblrRestClient)
	config common.Config
	//lastMsg string
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

	return nil
}

//PostImage post an image on tumblr
func PostImage(status common.Message) (err error) {
	msg := status.Text
	imgRegexp := regexp.MustCompile(config.Get("ajaxchat", "img_regex"))
	imagesUrls := imgRegexp.FindAllString(msg, -1)
	if len(imagesUrls) == 0 {
		return
	}
	tagsRe := regexp.MustCompile("\\[(\\S+?)\\]")
	tags := tagsRe.FindAllString(msg, -1)
	for i, tag := range tags {
		tags[i] = strings.Replace(strings.Trim(tag, "[] "), " ", "_", -1)
	}
	tagList := strings.Join(tags, ", ")
	for _, image := range imagesUrls {
		options := map[string]string{
			"tags":   tagList,
			"source": image,
		}
		Log(fmt.Sprintf("Posting on tumblr: %s with tags %s", image, tags))
		err = client.CreatePhoto(config.Get("tumblr", "url"), options)
		if err != nil {
			Log(fmt.Sprintf("Error,failed to post image %s", err))
			return err
		}
	}
	return
}
