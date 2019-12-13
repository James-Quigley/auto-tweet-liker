package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

type Effector func(client *twitter.Client, tweet *twitter.Tweet) error

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(client *twitter.Client, tweet *twitter.Tweet) error {
		for r := 0; ; r++ {
			err := effector(client, tweet)
			if err == nil || r >= retries {
				return err
			}
			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)
			select {
			case <-time.After(delay):
			}
		}
	}
}

func likeTweetOnDelay(client *twitter.Client, tweet *twitter.Tweet) error {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(600)
	log.Printf("Waiting %d seconds to like tweet %d", n, tweet.ID)
	time.Sleep(time.Duration(n) * time.Second)

	favParams := &twitter.FavoriteCreateParams{
		ID: tweet.ID,
	}
	_, _, err := client.Favorites.Create(favParams)
	if err != nil {
		log.Fatalf("Failed to like tweet: %d", tweet.ID)
		return err
	}
	log.Printf("Liked tweet: %d", tweet.ID)
	return nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Successfully loaded environment variables")

	TwitterConsumerAPIKey := os.Getenv("TWITTER_CONSUMER_API_KEY")
	TwitterConsumerASecret := os.Getenv("TWITTER_CONSUMER_API_SECRET")
	TwitterAccessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	TwitterAccessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	accountID := os.Getenv("TWITTER_ACCOUNT_ID")

	config := oauth1.NewConfig(TwitterConsumerAPIKey, TwitterConsumerASecret)
	token := oauth1.NewToken(TwitterAccessToken, TwitterAccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)
	log.Println("Twitter client created")

	params := &twitter.StreamFilterParams{
		Follow:        []string{accountID},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		go Retry(likeTweetOnDelay, 3, 2*time.Second)
	}
	log.Printf("Listening to stream for userID %s", accountID)
	for message := range stream.Messages {
		demux.Handle(message)
	}
}
