package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

type RecentTrackResp struct {
	RecentTracks struct {
		Tracks []Track `json:"track"`
	} `json:"recenttracks"`
}

func main() {
	username := os.Getenv("LASTFM_USERNAME")
	password := os.Getenv("LASTFM_PASSWORD")
	apiKey := os.Getenv("LASTFM_API_KEY")

	fmt.Println("Getting your tracks...")
	recents := getRecentTracks(username, apiKey)

	fmt.Println("Starting a web session...")
	webClient, err := NewLastFMWebClient(username, password)
	if err != nil {
		log.Fatalf("Failed to create a LastFM web client: %v", err)
	}

	list := tview.NewList()
	populateList := func() {
		recents = getRecentTracks(username, apiKey)
		for i, track := range recents.RecentTracks.Tracks[:9] {
			list.AddItem(track.Name, track.Artist.Name, rune(i+49), nil)
		}
	}
	populateList()
	list.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Name() == "Backspace2" || e.Name() == "Delete" {
			webClient.DeleteTrack(recents.RecentTracks.Tracks[list.GetCurrentItem()])
			list.Clear()
			populateList()
		}
		return e
	})
	app := tview.NewApplication()
	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}
}

func getRecentTracks(username, apiKey string) RecentTrackResp {
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
	decoded := RecentTrackResp{}
	err = json.Unmarshal(b, &decoded)
	if err != nil {
		panic(err)
	}
	return decoded
}
