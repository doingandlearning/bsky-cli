package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SearchPosts(authToken, searchTerm string) ([]interface{}, error) {
	searchURL := fmt.Sprintf("https://bsky.social/xrpc/app.bsky.feed.searchPosts?q=%s", searchTerm)
	req, err := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("Authorization", "Bearer "+authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send search request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve search results, status code: %d", resp.StatusCode)
	}

	var searchResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
	}

	// Verify structure and retrieve posts
	posts, ok := searchResponse["posts"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected structure for posts in search response")
	}

	return posts, nil
}
