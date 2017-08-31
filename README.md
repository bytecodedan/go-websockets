# Websockets with Go!
This is an example project showing us how to use WS (Websockets) with Go (golang). The WS takes a Twitter hastag and subscripts to live tweets. Tweets are sent back to the client through the WS. At any point in time, you can subscribe to new tweets via new hashtag. The current stream will be closed and new one created.

## Twitter API Keys and Token
You'll need the following environment variables defined. 

* TWITTER_CONSUMER_KEY
* TWITTER_CONSUMER_SECRET
* TWITTER_ACCESS_TOKEN
* TWITTER_ACCESS_TOKEN_SECRET

If you don't have these Twitter credentials yet, simply create a developer account, create an app, and follow the instructions to have them generated. See [Twitter Developer Docs](https://dev.twitter.com/rest/public).

## Thanks to ...
[Francesc Campoy](https://github.com/campoy) and his [example](https://github.com/campoy/justforfunc/tree/master/14-twitterbot) of using the [Anaconda](https://github.com/ChimeraCoder/anaconda) client library for Twitter.