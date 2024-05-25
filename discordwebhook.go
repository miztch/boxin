package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	twitterDomain     = "x.com"
	twitterBrandColor = 1942002
)

// Payload is a payload for a Discord webhook
type Payload struct {
	Embeds    []Embed `json:"embeds"`
	Username  *string `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url"`
}

// Embed is an embed for a Discord webhook
type Embed struct {
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       int    `json:"color"`
	Timestamp   string `json:"timestamp"`
	Author      Author `json:"author"`
	Footer      Footer `json:"footer"`
}

// Author is an author for a Discord webhook
type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

// Footer is a footer for a Discord webhook
type Footer struct {
	Text string `json:"text"`
}

// buildEmbed builds an embed for a Discord webhook
func buildEmbed(tweet Tweet) Embed {
	author := Author{
		Name:    fmt.Sprintf("%s (%s)", tweet.authorName, tweet.authorId),
		URL:     tweet.url,
		IconURL: tweet.icon,
	}

	footer := Footer{
		Text: twitterDomain,
	}

	return Embed{
		Description: tweet.body,
		URL:         tweet.url,
		Color:       twitterBrandColor,
		Timestamp:   tweet.time,
		Author:      author,
		Footer:      footer,
	}
}

// buildPayload builds a payload for a Discord webhook
func buildPayload(embed Embed, webhookConfig webhookConfig) *Payload {
	return &Payload{
		Embeds:    []Embed{embed},
		Username:  webhookConfig.Username,
		AvatarURL: webhookConfig.AvatarURL,
	}
}

// sendRequest sends a POST request to a URL with a payload
func sendRequest(url string, payloadBytes []byte) (*http.Response, error) {
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("[error] failed to send request, %w", err)
	}
	return res, nil
}

// DoToWebhook sends a tweet to a Discord webhook
func doWebhook(webhookConfig webhookConfig, tweet Tweet) (int, error) {
	embed := buildEmbed(tweet)
	payload := buildPayload(embed, webhookConfig)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("[error] failed to marshal payload, %w", err)
	}

	res, err := sendRequest(webhookConfig.URL, payloadBytes)
	if err != nil {
		return 0, fmt.Errorf("[error] failed to send request, %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return res.StatusCode, fmt.Errorf("[error] failed to send webhook. statusCode: %v", res.StatusCode)
	}
	return res.StatusCode, err
}
