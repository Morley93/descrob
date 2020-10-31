package main

import (
	"fmt"
	"strings"

	"github.com/Morley93/descrob"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type app struct {
	*tview.Application
	listCtrl    *tview.List
	webClient   *descrob.LastFMWebClient
	currentPage int
	tracks      []descrob.Track
	user        string
	apiKey      string
}

func newTUIApp(webClient *descrob.LastFMWebClient, user, apiKey string) *app {
	tuiApp, listCtrl := createTviewApp()
	app := app{
		Application: tuiApp,
		listCtrl:    listCtrl,
		webClient:   webClient,
		user:        user,
		apiKey:      apiKey,
	}
	app.installKeyHandlers()
	return &app
}

func (a *app) installKeyHandlers() {
	a.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Key() {
		case tcell.KeyBackspace2:
			// TODO: Handle error
			a.webClient.DeleteTrack(a.tracks[a.listCtrl.GetCurrentItem()])
			a.renderTracks()
		case tcell.KeyCtrlN:
			// TODO: Handle error
			a.nextPage()
		case tcell.KeyCtrlP:
			// TODO Handle error
			a.prevPage()
		}
		return e
	})
}

func (a *app) renderTracks() {
	a.listCtrl.Clear()
	for i, track := range a.tracks[:9] {
		a.listCtrl.AddItem(track.Name, track.Artist.Name, rune(i+0x31), nil)
	}
}

func (a *app) nextPage() error {
	a.currentPage++
	tracks, err := descrob.GetRecentTracks(a.user, a.apiKey, a.currentPage)
	if err != nil {
		return fmt.Errorf("Error getting track page: %w", err)
	}
	a.tracks = tracks
	a.renderTracks()
	return nil
}

func (a *app) prevPage() error {
	if a.currentPage == 1 {
		return nil
	}
	a.currentPage--
	tracks, err := descrob.GetRecentTracks(a.user, a.apiKey, a.currentPage)
	if err != nil {
		return fmt.Errorf("Error getting track page: %w", err)
	}
	a.tracks = tracks
	a.renderTracks()
	return nil
}

func createTviewApp() (*tview.Application, *tview.List) {
	app := tview.NewApplication()
	list := tview.NewList()

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(createKeybindView(), 1, 0, false)
	app.SetRoot(flex, true)
	return app, list
}

func createKeybindView() *tview.TextView {
	keybinds := []string{
		"[Ctrl+n] Next page",
		"[Ctrl+p] Previous page",
		"[Bckspc] Unscrobble",
	}
	return tview.NewTextView().
		SetText(strings.Join(keybinds, " | "))
}
