package main

import (
	"fmt"
	"strconv"
	"time"
)

const (
	twitterSnowflakeEpoch = int64(1288834974657)
)

// Tweet is a tweet from Yahoo! Realtime Search
type Tweet struct {
	id         string
	url        string
	time       string
	icon       string
	body       string
	authorId   string
	authorName string
}

// getTimestamp converts a tweet ID to an ISO 8601 timestamp
func (t *Tweet) getTimestamp() (string, error) {
	tweetIdInt64, err := strconv.ParseInt(t.id, 10, 64)
	if err != nil {
		return "", fmt.Errorf("[error] failed to parse tweet ID, %w", err)
	}

	tweetTime := (tweetIdInt64 >> 22) + twitterSnowflakeEpoch
	tweetUnixTime := time.Unix(tweetTime/1000, (tweetTime%1000)*int64(time.Millisecond))
	iso8601Timestamp := tweetUnixTime.Format(time.RFC3339)

	return iso8601Timestamp, nil
}

// isOnToday returns if the tweet is on today in JST or not
func (t *Tweet) isOnToday() (bool, error) {
	tweetTimeUTC, err := time.Parse(time.RFC3339, t.time)
	if err != nil {
		return false, fmt.Errorf("[error] failed to parse tweet timestamp, %w", err)
	}

	tweetTimeJST := UTCToJST(tweetTimeUTC)
	currentDateJST := getCurrentDateJST()
	tweetDateJST := tweetTimeJST.Format("2006-01-02")

	return tweetDateJST == currentDateJST, nil
}
