package main

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("Missing required environment variable " + name)
	}
	return v
}

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }

var (
	api      *anaconda.TwitterApi
	stream   *anaconda.Stream
	tweetReq = make(chan string)
	tweetRes = make(chan string)
)

func initTwitter() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	log = &logger{logrus.New()}
	api.SetLogger(log)

	go handleRequests()
}

func handleRequests() {
	go func() {
		for {
			q := <-tweetReq
			if stream != nil {
				stream.Stop()
			}
			stream = api.PublicStreamFilter(url.Values{
				"track": []string{q},
			})
			go handleStream()

		}
	}()
}

func handleStream() {
	defer stream.Stop()
	// send stream of tweets to the response channel
	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			log.Warningf("Received unexpected value type of %T\n", v)
			continue
		}
		log.Infof("%v\n", t.Text)
		tweetRes <- t.Text
	}
}
