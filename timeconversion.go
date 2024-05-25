package main

import (
	"time"
)

// getCurrentDateJST returns the current date in JST
func getCurrentDateJST() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	currentJSTTime := time.Now().In(jst)

	currentJSTDate := currentJSTTime.Format("2006-01-02")

	return currentJSTDate
}

func UTCToJST(timeUTC time.Time) time.Time {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	timeJST := timeUTC.In(jst)
	return timeJST
}

func JSTToUTC(timeJST time.Time) time.Time {
	timeUTC := timeJST.UTC()
	return timeUTC
}
