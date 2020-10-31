package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

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

	page := 1
	var recents []descrob.Track
	list := tview.NewList()
	populateList := func(page int) {
		list.Clear()
		recents, err = descrob.GetRecentTracks(username, apiKey, page)
		if err != nil {
			//TODO: No way of signalling this back yet
			panic(fmt.Sprintf("A recent track request failed: %v", err))
		}
		for i, track := range recents[:9] {
			list.AddItem(track.Name, track.Artist.Name, rune(i+0x31), nil)
		}
	}
	populateList(page)
	app := tview.NewApplication()
	list.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Key() {
		case tcell.KeyBackspace2:
			webClient.DeleteTrack(recents[list.GetCurrentItem()])
			populateList(page)
			break
		case tcell.KeyCtrlN:
			page++
			populateList(page)
		case tcell.KeyCtrlP:
			page = int(math.Max(1.0, float64(page-1)))
			populateList(page)
		}
		return e
	})

	keybinds := []string{
		"[Ctrl+n] Next page",
		"[Ctrl+p] Previous page",
		"[Bckspc] Unscrobble",
	}
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(tview.NewTextView().
			SetRegions(true).
			SetText(strings.Join(keybinds, " | ")), 1, 0, false)
	app.SetRoot(flex, true)

	log.Fatal(app.Run())
}
