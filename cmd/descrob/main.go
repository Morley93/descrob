package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Morley93/descrob"
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

	app := newTUIApp(webClient, username, apiKey)
	app.nextPage()

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
