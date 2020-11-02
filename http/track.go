package http

import (
	"strconv"
	"time"

	"github.com/Morley93/descrob"
)

type recentTracksResponse struct {
	RecentTracks struct {
		Tracks []trackResponse `json:"track,omitempty"`
	} `json:"recenttracks"`
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

func (rtt recentTracksResponse) mapToScrobbles() []descrob.Scrobble {
	scrobbles := []descrob.Scrobble{}

	if rtt.RecentTracks.Tracks == nil {
		return scrobbles
	}
	for _, trackResp := range rtt.RecentTracks.Tracks {
		scrobbleTimestamp, err := strconv.Atoi(trackResp.Date.Timestamp)
		if err != nil {
			scrobbleTimestamp = 0
		}
		track := descrob.Scrobble{
			Name:     trackResp.Name,
			Artist:   trackResp.Artist.Name,
			Datetime: time.Unix(int64(scrobbleTimestamp), 0),
		}
		scrobbles = append(scrobbles, track)
	}
	return scrobbles
}
