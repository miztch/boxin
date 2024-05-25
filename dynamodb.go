package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	partitionKey = "tweetDate"
	sortKey      = "tweetId"
)

// DynamoDBClient is a wrapper for the DynamoDB client
type DynamoDBClient struct {
	client    *dynamodb.Client
	tableName string
}

func newDynamoDBClient(ctx context.Context, tableName string) (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("[error] failed to load configuration, %w", err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBClient{client: client, tableName: tableName}, nil
}

func (d *DynamoDBClient) checkItemExists(ctx context.Context, currentDate string) (bool, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{
				Value: currentDate,
			},
		},
	}

	output, err := d.client.GetItem(ctx, input)
	if err != nil {
		return false, fmt.Errorf("[error] failed to get item, %w", err)
	}

	return len(output.Item) > 0, nil
}

func (d *DynamoDBClient) saveSentRecord(ctx context.Context, jstDate string, tweetId string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{
				Value: jstDate,
			},
			sortKey: &types.AttributeValueMemberS{
				Value: tweetId,
			},
		},
	}

	_, err := d.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("[error] failed to put item, %w", err)
	}

	return nil
}
