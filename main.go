package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) {
	// Create a DynamoDB client
	dbClient, err := newDynamoDBClient(ctx, getDynamoDBConfig().TableName)
	if err != nil {
		log.Printf("[error] Failed to create DynamoDB client: %v", err)
		return
	}

	// Check if today's tweets have already been sent
	// If it has, exit
	currentDateJST := getCurrentDateJST()
	if exists, err := dbClient.checkItemExists(ctx, currentDateJST); err != nil {
		log.Printf("[error] Failed to check if tweets sent today: %v", err)
		return
	} else if exists {
		log.Println("[info] Today's tweet have already been sent")
		return
	}

	// Search for tweets
	tweets, err := realtimeSearch(getScrapeConfig())
	if err != nil {
		log.Printf("[error] Failed to search latest tweet: %v", err)
		return
	}

	// if no tweets are found, exit
	if len(tweets) == 0 {
		log.Println("[info] No tweets found")
		return
	}

	// Check if the latest tweet is from today in JST
	latestTweet := tweets[0]
	haveUnsentTweet, err := latestTweet.isOnToday()
	if err != nil {
		log.Printf("[error] Failed to check if the latest tweet is from today: %v", err)
		return
	}

	// If it's from today in JST, it means it hasn't been sent yet, so send it to the webhook
	// If it's from yesterday, it means there are no tweets for today yet, so exit
	if !haveUnsentTweet {
		log.Println("[info] Today's tweet have not been found yet")
		return
	}

	// Send the tweet to the webhook
	statusCode, err := doWebhook(getWebhookConfig(), latestTweet)
	if err != nil || statusCode >= 400 {
		log.Printf("[error] Failed to send to webhook. status Code: %v, error: %v", statusCode, err)
		return
	}

	// Save the sent record in DynamoDB
	err = dbClient.saveSentRecord(ctx, currentDateJST, latestTweet.id)
	if err != nil {
		log.Printf("[error] Failed to save sent record to DynamoDB: %v", err)
		return
	}
}

func main() {
	if isRunningOnLambda() {
		lambda.Start(handler)
	} else {
		err := LoadEnvFromConfig()
		if err != nil {
			log.Printf("[error] Failed to load environment variables: %v", err)
			return
		}
		handler(context.Background())
	}
}
