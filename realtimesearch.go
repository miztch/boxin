package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	realtimeSearchEndpoint = "https://search.yahoo.co.jp/realtime/search"
	realtimeSearchDomain   = "search.yahoo.co.jp"
	crawlerUserAgent       = "Googlebot/2.1 (+http://www.google.com/bot.html)"
	keywordQueryParam      = "?p="
	authorIdPrefix         = "+id:"
)

// buildRequestURL builds a request URL for Yahoo! Realtime Search
func buildRequestURL(scrapeConfig scrapeConfig) (string, error) {
	if scrapeConfig.keyword == "" {
		return "", fmt.Errorf("[error] please set keyword")
	}

	queryParam := keywordQueryParam + url.QueryEscape(scrapeConfig.keyword)
	if len(scrapeConfig.authorId) > 0 {
		queryParam += authorIdPrefix + url.QueryEscape(scrapeConfig.authorId)
	}

	requestURL := realtimeSearchEndpoint + queryParam

	return requestURL, nil
}

// setupColly sets up a new colly collector
func setupColly() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(realtimeSearchDomain),
		colly.UserAgent(crawlerUserAgent),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       time.Second,
		RandomDelay: time.Second,
	},
	)

	c.OnError(func(_ *colly.Response, err error) {
		panic(err)
	})

	return c
}

// parseTweet parses a tweet from a Yahoo! Realtime Search result
func parseTweet(e *colly.HTMLElement) (Tweet, error) {
	tweet := Tweet{
		icon:       e.ChildAttr("img", "src"),
		body:       e.ChildText("div[class^='Tweet_bodyContainer'] div[class^='Tweet_body']"),
		authorName: e.ChildText("span[class^='Tweet_authorName__']"),
		authorId:   e.ChildText("a[class^='Tweet_authorID__']"),
	}

	tweetURLWithQuery := e.ChildAttr("time[class^='Tweet_time'] a", "href")
	tweet.url = strings.Split(tweetURLWithQuery, "?")[0]
	if tweet.url == "" {
		return Tweet{}, fmt.Errorf("[error] failed to get tweet URL")
	}

	tweetIDRegex := regexp.MustCompile(`/status/(\d+)`)
	tweet.id = tweetIDRegex.FindStringSubmatch(tweet.url)[1]
	if tweet.id == "" {
		return Tweet{}, fmt.Errorf("[error] failed to get tweet ID")
	}

	tweetTime, err := tweet.getTimestamp()
	if err != nil {
		return Tweet{}, fmt.Errorf("[error] failed to get tweet timestamp, %w", err)
	}
	tweet.time = tweetTime

	return tweet, nil
}

// realtimeSearch scrapes Yahoo! Realtime Search and returns tweets
func realtimeSearch(scrapeConfig scrapeConfig) ([]Tweet, error) {
	requestURL, err := buildRequestURL(scrapeConfig)
	if err != nil {
		return nil, fmt.Errorf("[error] failed to build request URL, %w", err)
	}

	c := setupColly()

	var tweets []Tweet
	c.OnHTML("div[class^='Tweet_TweetContainer']", func(e *colly.HTMLElement) {
		tweet, err := parseTweet(e)
		if err != nil {
			log.Println("[error] failed to parse a tweet, %w", err)
			return
		}
		tweets = append(tweets, tweet)
	})

	err = c.Visit(requestURL)
	if err != nil {
		return nil, fmt.Errorf("[error] failed to visit Yahoo! Realtime Search, %w", err)
	}

	return tweets, nil
}
