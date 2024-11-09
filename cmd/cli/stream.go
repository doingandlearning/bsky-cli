package main

import (
	"fmt"
	"time"
)

func streamLatestPosts(authToken, did string, interval time.Duration) {
	var lastSeenURI string

	for {
		// Fetch the latest posts
		posts, err := FetchLatestPosts(authToken, did)
		if err != nil {
			fmt.Printf("Error fetching latest posts: %v\n", err)
			time.Sleep(interval)
			continue
		}

		// Iterate over posts and print only new ones
		for _, post := range posts {
			postData := post.(map[string]interface{})
			uri, _ := postData["post"].(map[string]interface{})["uri"].(string)

			// Stop printing if we reach the last seen post
			if uri == lastSeenURI {
				break
			}

			// Display the post
			PrintPost(postData)
		}

		// Update the last seen URI
		if len(posts) > 0 {
			firstPost := posts[0].(map[string]interface{})
			lastSeenURI, _ = firstPost["post"].(map[string]interface{})["uri"].(string)
		}

		// Wait for the next fetch
		time.Sleep(interval)
	}
}
