package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"flag"
	"os"
	"strings"

	"github.com/joho/godotenv"
)



func init() {
	// Load .env file at the start of the application
	err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error loading .env file")
	}
	
}

func getAuthToken(username, appPassword string) (string, string, error) {
	loginURL := "https://bsky.social/xrpc/com.atproto.server.createSession"

	// Login payload
	loginPayload := map[string]string{
		"identifier": username,
		"password":   appPassword,
	}

	// Encode login payload as JSON
	payloadBytes, err := json.Marshal(loginPayload)
	if err != nil {
		return "", "", fmt.Errorf("failed to encode login payload: %v", err)
	}

	// Make the login request
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", "", fmt.Errorf("failed to create login request: %v", err)
	}

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send login request: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to login, status code: %d", resp.StatusCode)
	}

	// Parse the response
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", "", fmt.Errorf("failed to decode login response: %v", err)
	}

	// Extract the auth token and DID
	authToken, ok := responseData["accessJwt"].(string)
	if !ok {
		return "", "", fmt.Errorf("failed to retrieve auth token")
	}
	did, ok := responseData["did"].(string)
	if !ok {
		return "", "", fmt.Errorf("failed to retrieve DID")
	}

	return authToken, did, nil
}

func createPost(authToken, did, postContent string) error {
	postURL := "https://bsky.social/xrpc/com.atproto.repo.createRecord"

	// Post payload
	postData := map[string]interface{}{
		"repo":       did,
		"collection": "app.bsky.feed.post",
		"record": map[string]interface{}{
			"$type":     "app.bsky.feed.post",
			"text":      postContent,
			"createdAt": time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Encode post payload as JSON
	postDataJSON, err := json.Marshal(postData)
	if err != nil {
		return fmt.Errorf("failed to encode post data: %v", err)
	}

	// Make the post request
	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(postDataJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	if err != nil {
		return fmt.Errorf("failed to create post request: %v", err)
	}

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send post request: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create post, status code: %d", resp.StatusCode)
	}

	fmt.Println("Post created successfully!")
	return nil
}

func getLatestPosts(authToken, did string) error {
	feedURL := "https://bsky.social/xrpc/app.bsky.feed.getTimeline?limit=10"

	req, err := http.NewRequest("GET", feedURL, nil)
	req.Header.Set("Authorization", "Bearer "+authToken)
	if err != nil {
			return fmt.Errorf("failed to create feed request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
			return fmt.Errorf("failed to send feed request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to retrieve feed, status code: %d", resp.StatusCode)
	}

	var feedResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&feedResponse); err != nil {
			return fmt.Errorf("failed to decode feed response: %v", err)
	}

	posts, ok := feedResponse["feed"].([]interface{})
	if !ok {
			return fmt.Errorf("failed to parse feed items")
	}

	fmt.Println("Latest 10 Posts:")
	for i, post := range posts {
			postData, ok := post.(map[string]interface{})
			if !ok {
					continue
			}

			// Get the author's display name
			authorData, ok := postData["post"].(map[string]interface{})["author"].(map[string]interface{})
			if !ok {
					fmt.Printf("%d: (No author field)\n", i+1)
					continue
			}
			displayName, _ := authorData["displayName"].(string)

			// Navigate to `record` for post text
			record, ok := postData["post"].(map[string]interface{})["record"].(map[string]interface{})
			if !ok {
					fmt.Printf("%d: (No record field)\n", i+1)
					continue
			}

			// Extract the `text` field
			text, _ := record["text"].(string)

			// Extract the `uri` field and transform it into a URL
			uri, _ := postData["post"].(map[string]interface{})["uri"].(string)
			url := transformUriToUrl(uri)

			// Print in `user: message (URL)` format
			fmt.Printf("%d: %s: %s (%s)\n", i+1, displayName, text, url)
	}

	return nil
}

// Helper function to transform the Bluesky URI into a URL
func transformUriToUrl(uri string) string {
	// Example uri: "at://did:plc:xyz123/app.bsky.feed.post/abcdefg"
	parts := strings.Split(uri, "/")
	if len(parts) < 5 {
			return uri // return original uri if format is unexpected
	}
	
	did := parts[2]            // Extract DID part
	postId := parts[len(parts)-1] // Extract Post ID part
	return fmt.Sprintf("https://bsky.app/profile/%s/post/%s", did, postId)
}


func main() {
	username := os.Getenv("USERNAME")
	appPassword := os.Getenv("APP_PASSWORD")

	if username == "" || appPassword == "" {
			log.Fatalf("Error: Username and App Password must be set in .env file")
	}

	authToken, did, err := getAuthToken(username, appPassword)
	if err != nil {
			log.Fatalf("Error during login: %v", err)
	}

	// Define and parse flags
	postContent := flag.String("content", "", "Content of the post")
	fetchFeed := flag.Bool("fetch", false, "Fetch the latest 10 posts from the feed")
	flag.Parse()

	if *fetchFeed {
			// Fetch the latest posts if -fetch flag is provided
			if err := getLatestPosts(authToken, did); err != nil {
					log.Fatalf("Error fetching latest posts: %v", err)
			}
	} else if *postContent != "" {
			// Create a post if -content flag is provided
			if err := createPost(authToken, did, *postContent); err != nil {
					log.Fatalf("Error creating post: %v", err)
			}
	} else {
			log.Fatalf("Error: Provide either -content to post or -fetch to retrieve feed")
	}
}