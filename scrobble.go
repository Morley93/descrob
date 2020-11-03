package descrob

import (
	"math"
	"time"
)

type Scrobble struct {
	Name     string
	Artist   string
	Datetime time.Time
}

type ScrobbleExplorer struct {
	cache       []Scrobble
	sr          ScrobbleRetriever
	username    string
	windowSize  int
	windowStart int
}

type ScrobbleRetriever interface {
	FetchScrobblePage(username string, page int) ([]Scrobble, error)
	PageSize() int
}

func NewScrobbleExplorer(username string, sr ScrobbleRetriever, windowSize int) *ScrobbleExplorer {
	return &ScrobbleExplorer{
		cache:      []Scrobble{},
		sr:         sr,
		username:   username,
		windowSize: windowSize,
	}
}

func (se *ScrobbleExplorer) CurrentPage() []Scrobble {
	windowStart := int(math.Max(0, float64(se.windowStart)))
	windowEnd := int(math.Min(float64(len(se.cache)), float64(se.windowStart)))
	return se.cache[windowStart:windowEnd]
}

func (se *ScrobbleExplorer) FirstPage() ([]Scrobble, error) {
	if len(se.cache) == 0 {
		scrobbles, err := se.sr.FetchScrobblePage(se.username, 0)
		if err != nil {
			return nil, err
		}
		for _, scrob := range scrobbles {
			se.cache = append(se.cache, scrob)
		}
	}
	windowEnd := int(math.Min(float64(se.windowStart+se.windowSize), float64(len(se.cache))))
	return se.cache[se.windowStart:windowEnd], nil
}

func (se *ScrobbleExplorer) NextPage() ([]Scrobble, error) {
	if len(se.cache) < se.windowStart+se.windowSize {
		se.BufferNextWindow()
	}
	se.windowStart += se.windowSize
	return se.cache[se.windowStart : se.windowStart+se.windowSize], nil
}

func (se *ScrobbleExplorer) PrevPage() ([]Scrobble, error) {
	if se.windowStart == 0 {
		return se.cache[0:se.windowSize], nil
	}
	se.windowStart -= se.windowSize
	return se.cache[se.windowStart : se.windowStart+se.windowSize], nil
}

func (se *ScrobbleExplorer) RefreshPage() ([]Scrobble, error) {
	return nil, nil
}

func (se *ScrobbleExplorer) BufferWindows(windows int) {
	var scrobblesFetched int
	targetScrobblesFetched := windows * se.windowSize
	lastPageFetched := len(se.cache) / se.sr.PageSize()
	if len(se.cache) == 0 {
		lastPageFetched = -1
	}
	for {
		lastPageFetched++
		newScrobbles, err := se.sr.FetchScrobblePage(se.username, lastPageFetched)
		if err != nil {
			panic(err)
		}
		scrobblesFetched += len(newScrobbles)
		for _, s := range newScrobbles {
			se.cache = append(se.cache, s)
		}
		if scrobblesFetched >= targetScrobblesFetched || len(newScrobbles) < se.sr.PageSize() {
			break
		}
	}
}

func (se *ScrobbleExplorer) BufferNextWindow() {
	se.BufferWindows(1)
}

func (se *ScrobbleExplorer) BufferedWindows() int {
	return (len(se.cache) - (se.windowStart + se.windowSize)) / se.windowSize
}
