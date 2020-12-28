package main

import (
	"strings"

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
	user      string
	apiKey    string
}

func newTUIApp(webClient *descrob.LastFMWebClient, expl *descrob.ScrobbleExplorer, user, apiKey string) *app {
	tuiApp, listCtrl := createTviewApp()
	app := app{
		Application: tuiApp,
		listCtrl:    listCtrl,
		webClient:   webClient,
		expl:        expl,
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
			scrobbleToDelete := a.expl.CurrentPage()[a.listCtrl.GetCurrentItem()]
			// TODO: Handle errors
			a.webClient.DeleteTrack(scrobbleToDelete)
			a.expl.Uncache(scrobbleToDelete)
			a.renderScrobbles(a.expl.CurrentPage())
		case tcell.KeyCtrlN:
			// TODO: Handle error
			scrobs, _ := a.expl.NextPage()
			if a.expl.PreBufferedWindows() < 3 {
				// TODO: There's no locking around this this
				go a.expl.BufferWindows(3)
			}
			a.renderScrobbles(scrobs)
		case tcell.KeyCtrlP:
			scrobs := a.expl.PrevPage()
			a.renderScrobbles(scrobs)
		}
		switch e.Rune() {
		case 'j':
			a.listCtrl.SetCurrentItem(a.listCtrl.GetCurrentItem() + 1)
		case 'k':
			a.listCtrl.SetCurrentItem(a.listCtrl.GetCurrentItem() - 1)
		}
		return e
	})
}

func (a *app) renderScrobbles(scrobbles []descrob.Scrobble) {
	a.listCtrl.Clear()
	for i, scrobble := range scrobbles[:9] {
		a.listCtrl.AddItem(scrobble.Name, scrobble.Artist, rune(i+'0'), nil)
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
		"[Ctrl+n] Next page",
		"[Ctrl+p] Previous page",
		"[Bckspc] Unscrobble",
	}
	return tview.NewTextView().
		SetText(strings.Join(keybinds, " | "))
}
