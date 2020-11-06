package descrob

import (
	"fmt"
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
	FetchScrobbles(username string, before *time.Time) ([]Scrobble, error)
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
	windowEnd := int(math.Min(float64(len(se.cache)), float64(se.windowStart+se.windowSize)))
	return se.cache[se.windowStart:windowEnd]
}

func (se *ScrobbleExplorer) FirstPage() ([]Scrobble, error) {
	if len(se.cache) == 0 {
		err := se.BufferNextWindow()
		if err != nil {
			return nil, fmt.Errorf("Error populating cache for first page: %w", err)
		}
	}
	windowEnd := int(math.Min(float64(se.windowSize), float64(len(se.cache))))
	return se.cache[0:windowEnd], nil
}

func (se *ScrobbleExplorer) NextPage() ([]Scrobble, error) {
	if len(se.cache) < se.windowStart+se.windowSize {
		err := se.BufferNextWindow()
		if err != nil {
			return nil, fmt.Errorf("Error populating cache for next page: %w", err)
		}
	}
	se.windowStart += se.windowSize
	return se.CurrentPage(), nil
}

func (se *ScrobbleExplorer) PrevPage() []Scrobble {
	if se.windowStart == 0 {
		return se.cache[0:se.windowSize]
	}
	se.windowStart -= se.windowSize
	return se.CurrentPage()
}

func (se *ScrobbleExplorer) BufferNextWindow() error {
	return se.BufferWindows(1)
}

func (se *ScrobbleExplorer) BufferWindows(windows int) error {
	var scrobblesFetched int
	targetScrobblesFetched := windows * se.windowSize
	for {
		var before *time.Time
		if len(se.cache) > 0 {
			before = &se.cache[len(se.cache)-1].Datetime
		}

		newScrobbles, err := se.sr.FetchScrobbles(se.username, before)
		if err != nil {
			return fmt.Errorf("Error fetching scrobbles: %w", err)
		}
		scrobblesFetched += len(newScrobbles)
		for _, s := range newScrobbles {
			se.cache = append(se.cache, s)
		}
		if scrobblesFetched >= targetScrobblesFetched || len(newScrobbles) < se.sr.PageSize() {
			break
		}
	}
	return nil
}

func (se *ScrobbleExplorer) PreBufferedWindows() int {
	currWindowEnd := se.windowStart + se.windowSize
	cachedAfterWindow := len(se.cache) - currWindowEnd
	return int(math.Max(float64(0), float64(cachedAfterWindow/se.windowSize)))
}

func (se *ScrobbleExplorer) Uncache(scrobble Scrobble) {
	for indInPage, cacheItem := range se.CurrentPage() {
		if cacheItem == scrobble {
			cacheInd := se.windowStart + indInPage
			se.cache = append(se.cache[:cacheInd], se.cache[cacheInd+1:]...)
		}
	}
}
