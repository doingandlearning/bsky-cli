package main

import (
	"fmt"
	"time"
)

func streamLatestPosts(authToken, did string, interval time.Duration) {
	var lastSeenURI string
	displayedReposts := make(map[string]bool) // Track reposts that have been displayed

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
			postData, ok := post.(map[string]interface{})
			if !ok {
				continue
			}

			uri, _ := postData["post"].(map[string]interface{})["uri"].(string)

			// Stop printing if we reach the last seen post
			if uri == lastSeenURI {
				break
			}

			// Check if this is a repost by looking for the "reason" field
			if _, isRepost := postData["reason"]; isRepost {
				// Skip if this repost has already been displayed
				if displayedReposts[uri] {
					continue
				}
				// Mark this repost URI as displayed
				displayedReposts[uri] = true
			}

			// Display the post
			PrintPost(postData)
		}

		// Update the last seen URI to the latest post's URI
		if len(posts) > 0 {
			firstPost := posts[0].(map[string]interface{})
			lastSeenURI, _ = firstPost["post"].(map[string]interface{})["uri"].(string)
		}

		// Wait for the next fetch
		time.Sleep(interval)
	}
}
