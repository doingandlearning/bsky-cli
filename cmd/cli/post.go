package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func createPost(authToken, did, postContent string) error {
	postURL := "https://bsky.social/xrpc/com.atproto.repo.createRecord"

	postData := map[string]interface{}{
		"repo":       did,
		"collection": "app.bsky.feed.post",
		"record": map[string]interface{}{
			"$type":     "app.bsky.feed.post",
			"text":      postContent,
			"createdAt": time.Now().UTC().Format(time.RFC3339),
		},
	}

	postDataJSON, err := json.Marshal(postData)
	if err != nil {
		return fmt.Errorf("failed to encode post data: %v", err)
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(postDataJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	if err != nil {
		return fmt.Errorf("failed to create post request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send post request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create post, status code: %d", resp.StatusCode)
	}

	fmt.Println("Post created successfully!")
	return nil
}

func getColorFromEnv(envVar string, defaultColor *color.Color) *color.Color {
	colorMap := map[string]*color.Color{
		"red":     color.New(color.FgRed),
		"green":   color.New(color.FgGreen),
		"yellow":  color.New(color.FgYellow),
		"blue":    color.New(color.FgBlue),
		"cyan":    color.New(color.FgCyan),
		"magenta": color.New(color.FgMagenta),
		"white":   color.New(color.FgWhite),
	}

	// Check environment variable, fallback to default color if unset or unrecognized
	colorName := os.Getenv(envVar)
	if customColor, exists := colorMap[colorName]; exists {
		return customColor
	}
	return defaultColor
}

func PrintPost(post interface{}) {
	// Load color preferences from environment or use defaults
	displayNameColor := getColorFromEnv("DISPLAY_NAME_COLOR", color.New(color.FgYellow))
	textColor := getColorFromEnv("TEXT_COLOR", color.New(color.FgBlue))
	urlColor := getColorFromEnv("URL_COLOR", color.New(color.FgCyan))

	postData, ok := post.(map[string]interface{})
	if !ok {
		return
	}

	authorData, _ := postData["post"].(map[string]interface{})["author"].(map[string]interface{})
	displayName, _ := authorData["displayName"].(string)
	record, _ := postData["post"].(map[string]interface{})["record"].(map[string]interface{})
	text, _ := record["text"].(string)
	uri, _ := postData["post"].(map[string]interface{})["uri"].(string)
	url := transformUriToUrl(uri)

	// Print using color functions
	fmt.Printf("%s: %s (%s)\n", displayNameColor.Sprint(displayName), textColor.Sprint(text), urlColor.Sprint(url))
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
