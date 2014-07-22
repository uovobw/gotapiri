// Package common contains all utility functions and
// common data structures used in the project
package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config maps the structure of the json configuration
// file in the form section: key/value as follows:
//          "item1" : { ... },
//          "item2" : { ... }
type Config map[string]ConfigItem

// ConfigItem contains a single section key/value map
type ConfigItem map[string]string

// ReadConfigFrom loads the configuration from a
// json file passed as the parameter and returns the object
// that it unmashalled
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

// Get requires a section and a key to retrieve a configuration value,
// so for the example:
//      {
//          section1 : { key1 : value1 , key2 : value2 },
//          section2 : { key3 : value3 , key4 : value4 }
//      }
// the call to get value3 would be Get("section2", "key3")
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
