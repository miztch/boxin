package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvFromConfig loads environment variables from a .env file
func LoadEnvFromConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	return nil
}

// scrapeConfig is a configuration for scraping Yahoo! Realtime Search
type scrapeConfig struct {
	keyword  string
	authorId string
}

// getScrapeConfig gets a scrape configuration from environment variables
func getScrapeConfig() scrapeConfig {
	return scrapeConfig{
		keyword:  os.Getenv("SEARCH_KEYWORD"),
		authorId: os.Getenv("SEARCH_AUTHOR_ID"),
	}
}

// WebhookConfig is a configuration for a Discord webhook
type webhookConfig struct {
	URL       string
	Username  *string
	AvatarURL string
}

// getWebhookConfig gets a webhook configuration from environment variables
func getWebhookConfig() webhookConfig {
	userName := os.Getenv("WEBHOOK_USERNAME")
	var userNamePointer *string
	if userName != "" {
		userNamePointer = &userName
	}

	return webhookConfig{
		URL:       os.Getenv("WEBHOOK_URL"),
		Username:  userNamePointer,
		AvatarURL: os.Getenv("WEBHOOK_AVATAR_URL"),
	}
}

// dynamoDBConfig is a configuration for DynamoDB
type dynamoDBConfig struct {
	TableName   string
	Region      *string
	EndpointURL *string
}

// getDynamoDBConfig gets a DynamoDB configuration from environment variables
func getDynamoDBConfig() dynamoDBConfig {
	region := os.Getenv("AWS_REGION")
	endpointURL := os.Getenv("DYNAMODB_ENDPOINT_URL")
	return dynamoDBConfig{
		TableName:   os.Getenv("HISTORY_TABLE"),
		Region:      &region,
		EndpointURL: &endpointURL,
	}
}
