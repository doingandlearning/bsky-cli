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

func main() {
	username := os.Getenv("USERNAME")
	appPassword := os.Getenv("APP_PASSWORD")

	if username == "" || appPassword == "" {
			log.Fatalf("Error: Username and App Password must be set in .env file")
	}

	postContent := flag.String("content", "", "Content of the post")
	flag.Parse()

	if *postContent == "" {
			log.Fatalf("Error: Please provide post content using the -content flag.")
	}
	// Authenticate and retrieve token and DID
	authToken, did, err := getAuthToken(username, appPassword)
	if err != nil {
		log.Fatalf("Error during login: %v", err)
	}


	// Create the post
	if err := createPost(authToken, did, *postContent); err != nil {
		log.Fatalf("Error creating post: %v", err)
	}
}
