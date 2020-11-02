package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Morley93/descrob"
)

type HTTPScrobbleRetriever struct {
	Client http.Client
	APIKey string
}

type trackResponse struct {
	Name   string `json:"name"`
	Artist struct {
		Name string `json:"#text"`
	} `json:"artist"`
	Date struct {
		Timestamp string `json:"uts"`
	} `json:"date"`
}

func (sr *HTTPScrobbleRetriever) FetchScrobblePage(username string, page int) ([]descrob.Scrobble, error) {
	req, err := sr.buildRecentTrackRequest(username, page)
	if err != nil {
		return nil, fmt.Errorf("Error creating recent track request: %w", err)
	}

	resp, err := sr.Client.Do(req)
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
			Tracks []trackResponse `json:"track"`
		} `json:"recenttracks"`
	}{}
	err = json.Unmarshal(b, &respPayload)
	if err != nil {
		return nil, fmt.Errorf("Unexpected response: %v", string(b))
	}
	return mapTracksResponse(respPayload.RecentTracks.Tracks), nil
}

func (sr *HTTPScrobbleRetriever) buildRecentTrackRequest(username string, page int) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "https://ws.audioscrobbler.com/2.0", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("method", "user.getrecenttracks")
	q.Add("user", username)
	q.Add("api_key", sr.APIKey)
	q.Add("format", "json")
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", "10")
	req.URL.RawQuery = q.Encode()

	return req, err
}

func mapTracksResponse(respElems []trackResponse) []descrob.Scrobble {
	tracks := []descrob.Scrobble{}
	for _, respTrack := range respElems {
		scrobbleTimestamp, err := strconv.Atoi(respTrack.Date.Timestamp)
		if err != nil {
			scrobbleTimestamp = 0
		}
		track := descrob.Scrobble{
			Name:     respTrack.Name,
			Artist:   respTrack.Artist.Name,
			Datetime: time.Unix(int64(scrobbleTimestamp), 0),
		}
		tracks = append(tracks, track)
	}
	return tracks
}
