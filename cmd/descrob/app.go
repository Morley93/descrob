package main

import (
	"math"
	"strings"
	"sync"

	"github.com/Morley93/descrob"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type app struct {
	*tview.Application
	listCtrl  *tview.List
	webClient *descrob.LastFMWebClient
	expl      *descrob.ScrobbleExplorer
	tracks    []descrob.Scrobble
	pageSize  int
	mut       sync.Mutex
}

func newTUIApp(webClient *descrob.LastFMWebClient, expl *descrob.ScrobbleExplorer, pageSize int, user, apiKey string) *app {
	tuiApp, listCtrl := createTviewApp()
	app := app{
		Application: tuiApp,
		listCtrl:    listCtrl,
		webClient:   webClient,
		expl:        expl,
		pageSize:    pageSize,
	}
	app.installKeyHandlers()
	return &app
}

func (a *app) installKeyHandlers() {
	a.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Key() {
		case tcell.KeyBackspace2:
			scrobbleToDelete := a.expl.CurrentPage()[a.listCtrl.GetCurrentItem()]
			// TODO: Handle errors
			a.webClient.DeleteTrack(scrobbleToDelete)
			a.expl.Uncache(scrobbleToDelete)
			a.renderScrobbles(a.expl.CurrentPage())
		}
		switch e.Rune() {
		case 'j':
			a.listCtrl.SetCurrentItem(a.listCtrl.GetCurrentItem() + 1)
		case 'k':
			a.listCtrl.SetCurrentItem(a.listCtrl.GetCurrentItem() - 1)
		case 'n':
			a.mut.Lock()
			defer a.mut.Unlock()
			// TODO: Handle error
			scrobs, _ := a.expl.NextPage()
			go func() {
				a.mut.Lock()
				defer a.mut.Unlock()
				if a.expl.PreBufferedWindows() < 3 {
					a.expl.BufferWindows(3)
				}
			}()
			a.renderScrobbles(scrobs)
		case 'p':
			a.mut.Lock()
			defer a.mut.Unlock()
			scrobs := a.expl.PrevPage()
			a.renderScrobbles(scrobs)
		}
		return e
	})
}

func (a *app) renderScrobbles(scrobbles []descrob.Scrobble) {
	a.listCtrl.Clear()
	lastScrobIndex := int(math.Min(float64(len(scrobbles)), float64(a.pageSize)))
	for i, scrobble := range scrobbles[:lastScrobIndex] {
		runeOffset := (i + 1) % a.pageSize
		a.listCtrl.AddItem(scrobble.Name, scrobble.Artist, rune('0'+runeOffset), nil)
	}
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
		"[j/↓] Next scrobble",
		"[k/↑] Previous scrobble",
		"[n] Next page",
		"[p] Previous page",
		"[Bckspc] Unscrobble",
	}
	return tview.NewTextView().
		SetText(strings.Join(keybinds, " | "))
}
