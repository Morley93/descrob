package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Morley93/descrob"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	username := os.Getenv("LASTFM_USERNAME")
	password := os.Getenv("LASTFM_PASSWORD")
	apiKey := os.Getenv("LASTFM_API_KEY")

	fmt.Println("Starting a web session...")
	webClient, err := descrob.NewLastFMWebClient(username, password)
	if err != nil {
		log.Fatalf("Failed to create a LastFM web client: %v", err)
	}

	var recents []descrob.Track
	list := tview.NewList()
	populateList := func() {
		recents = descrob.GetRecentTracks(username, apiKey)
		for i, track := range recents[:9] {
			list.AddItem(track.Name, track.Artist.Name, rune(i+0x31), nil)
		}
	}
	populateList()
	list.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Name() == "Backspace2" || e.Name() == "Delete" {
			webClient.DeleteTrack(recents[list.GetCurrentItem()])
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
