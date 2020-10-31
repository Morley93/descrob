package descrob

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Track struct {
	Name   string `json:"name"`
	Artist struct {
		Name string `json:"#text"`
	} `json:"artist"`
	Date struct {
		Timestamp string `json:"uts"`
	} `json:"date"`
}

func GetRecentTracks(username, apiKey string, page int) ([]Track, error) {
	req, err := buildRecentTrackRequest(username, apiKey, page)
	if err != nil {
		return nil, fmt.Errorf("Error creating recent track request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch recent tracks: %w", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %w", err)
	}
	respPayload := struct {
		RecentTracks struct {
			Tracks []Track `json:"track"`
		} `json:"recenttracks"`
	}{}
	err = json.Unmarshal(b, &respPayload)
	if err != nil {
		return nil, fmt.Errorf("Unexpected response: %v", string(b))
	}
	return respPayload.RecentTracks.Tracks, nil
}

func buildRecentTrackRequest(username, apiKey string, page int) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "https://ws.audioscrobbler.com/2.0", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("method", "user.getrecenttracks")
	q.Add("user", username)
	q.Add("api_key", apiKey)
	q.Add("format", "json")
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", "10")
	req.URL.RawQuery = q.Encode()

	return req, err
}
