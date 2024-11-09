package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func getAuthToken(username, appPassword string) (string, string, error) {
	loginURL := "https://bsky.social/xrpc/com.atproto.server.createSession"

	loginPayload := map[string]string{
		"identifier": username,
		"password":   appPassword,
	}

	payloadBytes, err := json.Marshal(loginPayload)
	if err != nil {
		return "", "", fmt.Errorf("failed to encode login payload: %v", err)
	}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", "", fmt.Errorf("failed to create login request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to login, status code: %d", resp.StatusCode)
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", "", fmt.Errorf("failed to decode login response: %v", err)
	}

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
