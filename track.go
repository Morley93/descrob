package descrob

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func GetRecentTracks(username, apiKey string) []Track {
	req, err := http.NewRequest(http.MethodGet, "https://ws.audioscrobbler.com/2.0", nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	q.Add("method", "user.getrecenttracks")
	q.Add("user", username)
	q.Add("api_key", apiKey)
	q.Add("format", "json")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	decoded := struct {
		RecentTracks struct {
			Tracks []Track `json:"track"`
		} `json:"recenttracks"`
	}{}
	err = json.Unmarshal(b, &decoded)
	if err != nil {
		panic(err)
	}
	return decoded.RecentTracks.Tracks
}
