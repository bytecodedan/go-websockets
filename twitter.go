package main

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

var (
	// api       *anaconda.TwitterApi
	// stream    *anaconda.Stream
	// tweetReq  chan string
	// tweetRes  chan string
	// tweetStop chan bool
	// Twitter API credentials
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

type Twitter struct {
	Api    *anaconda.TwitterApi
	Stream *anaconda.Stream
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("Missing required environment variable " + name)
	}
	return v
}

func (t *Twitter) Init() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	t.Api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	t.Api.SetLogger(log)
}

func (t *Twitter) SetupStream(s string) bool {
	if t.Stream != nil {
		t.Stop()
	}
	t.Stream = t.Api.PublicStreamFilter(url.Values{
		"track": []string{s},
	})
	return t.Stream != nil
}

func (t *Twitter) Stop() {
	if t.Stream != nil {
		t.Stream.Stop()
	}
}
