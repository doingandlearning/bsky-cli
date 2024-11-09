package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("USERNAME")
	appPassword := os.Getenv("APP_PASSWORD")
	if username == "" || appPassword == "" {
		log.Fatalf("Error: Username and App Password must be set in .env file")
	}

	authToken, did, err := getAuthToken(username, appPassword)
	if err != nil {
		log.Fatalf("Error during login: %v", err)
	}

	postContent := flag.String("content", "", "Content of the post")
	fetchFeed := flag.Bool("fetch", false, "Fetch the latest 10 posts from the feed")
	flag.Parse()

	switch {
	case *postContent != "":
		err := createPost(authToken, did, *postContent)
		if err != nil {
			log.Fatalf("Error creating post: %v", err)
		}
	case *fetchFeed:
		err := getLatestPosts(authToken, did)
		if err != nil {
			log.Fatalf("Error fetching latest posts: %v", err)
		}
	default:
		log.Fatalf("Error: Provide -content to post, -fetch to retrieve feed, or -stream to start streaming")
	}
}
