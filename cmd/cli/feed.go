package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchLatestPosts(authToken, did string) ([]interface{}, error) {
	feedURL := "https://bsky.social/xrpc/app.bsky.feed.getTimeline?limit=10"
	req, err := http.NewRequest("GET", feedURL, nil)
	req.Header.Set("Authorization", "Bearer "+authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create feed request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send feed request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve feed, status code: %d", resp.StatusCode)
	}

	var feedResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&feedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode feed response: %v", err)
	}

	posts, ok := feedResponse["feed"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse feed items")
	}

	return posts, nil
}

func getLatestPosts(authToken, did string) error {
	posts, err := FetchLatestPosts(authToken, did)
	if err != nil {
		return err
	}

	for _, post := range posts {
		PrintPost(post)
	}

	return nil
}
