package descrob

import (
	"fmt"
	"time"
)

type Scrobble struct {
	Name     string
	Artist   string
	Datetime time.Time
}

type ScrobbleExplorer struct {
	pageCache [][]Scrobble
	sr        ScrobbleRetriever
	pageIdx   int
	username  string
}

type ScrobbleRetriever interface {
	FetchScrobblePage(username string, page int) ([]Scrobble, error)
}

func NewScrobbleExplorer(username string, sr ScrobbleRetriever) *ScrobbleExplorer {
	return &ScrobbleExplorer{
		pageCache: [][]Scrobble{},
		sr:        sr,
		username:  username,
	}
}

func (se *ScrobbleExplorer) CurrentPage() []Scrobble {
	if len(se.pageCache) == 0 {
		return []Scrobble{}
	}
	return se.pageCache[se.pageIdx]
}

func (se *ScrobbleExplorer) FirstPage() ([]Scrobble, error) {
	se.pageIdx = 0
	if len(se.pageCache) == 0 {
		scrobbles, err := se.sr.FetchScrobblePage(se.username, se.pageIdx+1)
		if err != nil {
			return []Scrobble{}, fmt.Errorf("Failed to fetch page 1: %w", err)
		}
		se.pageCache = append(se.pageCache, scrobbles)
	}
	return se.pageCache[0], nil
}

func (se *ScrobbleExplorer) NextPage() ([]Scrobble, error) {
	se.pageIdx++
	if len(se.pageCache) < se.pageIdx+1 {
		scrobbles, err := se.sr.FetchScrobblePage(se.username, se.pageIdx)
		if err != nil {
			return []Scrobble{}, fmt.Errorf("Failed to fetch page %d: %w", se.pageIdx+1, err)
		}
		se.pageCache = append(se.pageCache, scrobbles)
	}
	return se.pageCache[se.pageIdx], nil
}

func (se *ScrobbleExplorer) PrevPage() ([]Scrobble, error) {
	if se.pageIdx == 0 {
		return se.pageCache[0], nil
	}
	se.pageIdx--
	return se.pageCache[se.pageIdx], nil
}

func (se *ScrobbleExplorer) RefreshPage() ([]Scrobble, error) {
	scrobbles, err := se.sr.FetchScrobblePage(se.username, se.pageIdx+1)
	if err != nil {
		return []Scrobble{}, fmt.Errorf("Failed to fetch page %d: %w", se.pageIdx+1, err)
	}
	se.pageCache[se.pageIdx] = scrobbles
	return scrobbles, nil
}
