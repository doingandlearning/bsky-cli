package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

	for i, post := range posts {
		postData, ok := post.(map[string]interface{})
		if !ok {
			continue
		}

		authorData, ok := postData["post"].(map[string]interface{})["author"].(map[string]interface{})
		if !ok {
			fmt.Printf("%d: (No author field)\n", i+1)
			continue
		}
		displayName, _ := authorData["displayName"].(string)

		record, ok := postData["post"].(map[string]interface{})["record"].(map[string]interface{})
		if !ok {
			fmt.Printf("%d: (No record field)\n", i+1)
			continue
		}

		text, _ := record["text"].(string)
		uri, _ := postData["post"].(map[string]interface{})["uri"].(string)
		url := transformUriToUrl(uri)

		fmt.Printf("%s: %s (%s)\n", displayName, text, url)
	}

	return nil
}

func transformUriToUrl(uri string) string {
	parts := strings.Split(uri, "/")
	if len(parts) < 5 {
		return uri
	}

	did := parts[2]
	postId := parts[len(parts)-1]
	return fmt.Sprintf("https://bsky.app/profile/%s/post/%s", did, postId)
}
