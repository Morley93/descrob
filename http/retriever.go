package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Morley93/descrob"
)

type HTTPScrobbleRetriever struct {
	client   *http.Client
	apiKey   string
	pageSize int
}

func NewHTTPScrobbleRetriever(client *http.Client, apiKey string, pageSize int) *HTTPScrobbleRetriever {
	return &HTTPScrobbleRetriever{
		client:   client,
		apiKey:   apiKey,
		pageSize: pageSize,
	}
}

func (sr *HTTPScrobbleRetriever) FetchScrobblePage(username string, page int) ([]descrob.Scrobble, error) {
	req, err := sr.buildRecentTrackRequest(username, page)
	if err != nil {
		return nil, fmt.Errorf("Error creating recent track request: %w", err)
	}

	resp, err := sr.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch recent tracks: %w", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %w", err)
	}
	respPayload := recentTracksResponse{}
	err = json.Unmarshal(b, &respPayload)
	if err != nil {
		return nil, fmt.Errorf("Unexpected response: %v", string(b))
	}
	return respPayload.mapToScrobbles(), nil
}

func (sr *HTTPScrobbleRetriever) buildRecentTrackRequest(username string, page int) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "https://ws.audioscrobbler.com/2.0", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("method", "user.getrecenttracks")
	q.Add("user", username)
	q.Add("api_key", sr.apiKey)
	q.Add("format", "json")
	q.Add("page", strconv.Itoa(page+1))
	q.Add("limit", strconv.Itoa(sr.pageSize))
	req.URL.RawQuery = q.Encode()

	return req, err
}

func (sr *HTTPScrobbleRetriever) PageSize() int {
	return sr.pageSize
}
