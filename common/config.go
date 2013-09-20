package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config map[string]ConfigItem
type ConfigItem map[string]string

func ReadConfigFrom(filename string) (c Config, e error) {
	c = Config{}
	pwd, _ := os.Getwd()
	b, e := ioutil.ReadFile(pwd + "/" + filename)
	if e != nil {
		return nil, e
	}
	e = json.Unmarshal(b, &c)
	if e != nil {
		return nil, e
	}
	return c, nil
}

func (c Config) Get(section, key string) (value string) {
	sectionVal, ok := c[section]
	if !ok {
		panic(fmt.Sprintf("no such config section %s", section))
	}
	value, ok = sectionVal[key]
	if !ok {
		panic(fmt.Sprintf("no key %s in section %s", key, section))
	}
	return value

}

/*
func writeConfigTo(filename string) {
	cfg := Config{
		General: map[string]string{"basedir": "/some/path"},
		Tumblr: map[string]string{
			"app_key":            "ddsfsfdsdfsfd",
			"app_secret":         "ddsfsfdsdfsfd",
			"oauth_token":        "ddsfsfdsdfsfd",
			"oauth_token_secret": "ddsfsfdsdfsfd",
			"url":                "tapiri.tumblr.com",
			"seenurls":           "seenTumblrUrls",
		},
		Ajaxchat: map[string]string{
			"user":        "asdfasdf",
			"channel":     "#asdfasdf",
			"ajaxuser":    "asdfasdfa",
			"ajaxpass":    "asdfasdf",
			"httpuser":    "asdfasdf",
			"httppass":    "asdfasdf",
			"msg_url":     "https://asdfasdfasdf.org/?ajax=true",
			"ajaxchannel": "Tapiri",
			"img_regex":   "(https?://[a-zA-Z0-9-.]+.[a-zA-Z]{2,3}(?:/S*)?(?:[a-zA-Z0-9_])+.(?:jpg|jpeg|gif|png))",
		},
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filename, b, 0700)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
*/
