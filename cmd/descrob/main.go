package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Morley93/descrob"
	dhttp "github.com/Morley93/descrob/http"
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

	scrobbleExplorer := descrob.NewScrobbleExplorer(username, &dhttp.HTTPScrobbleRetriever{
		Client: *http.DefaultClient,
		APIKey: apiKey,
	})

	app := newTUIApp(webClient, scrobbleExplorer, username, apiKey)
	initialScrobbles, err := scrobbleExplorer.FirstPage()
	if err != nil {
		log.Fatalf("Failed to get initial page of recent tracks: %v", err)
	}
	app.renderScrobbles(initialScrobbles)

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
