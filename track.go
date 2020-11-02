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

func (rte *ScrobbleExplorer) CurrentPage() []Scrobble {
	return rte.pageCache[rte.pageIdx]
}

func (rte *ScrobbleExplorer) FirstPage() ([]Scrobble, error) {
	rte.pageIdx = 0
	if len(rte.pageCache) == 0 {
		tracks, err := rte.sr.FetchScrobblePage(rte.username, rte.pageIdx+1)
		if err != nil {
			return []Scrobble{}, fmt.Errorf("Failed to fetch page 1: %w", err)
		}
		rte.pageCache = append(rte.pageCache, tracks)
	}
	return rte.pageCache[0], nil
}

func (rte *ScrobbleExplorer) NextPage() ([]Scrobble, error) {
	rte.pageIdx++
	if len(rte.pageCache) < rte.pageIdx+1 {
		tracks, err := rte.sr.FetchScrobblePage(rte.username, rte.pageIdx)
		if err != nil {
			return []Scrobble{}, fmt.Errorf("Failed to fetch page %d: %w", rte.pageIdx+1, err)
		}
		rte.pageCache = append(rte.pageCache, tracks)
	}
	return rte.pageCache[rte.pageIdx], nil
}

func (rte *ScrobbleExplorer) PrevPage() ([]Scrobble, error) {
	if rte.pageIdx == 0 {
		return rte.pageCache[0], nil
	}
	rte.pageIdx--
	return rte.pageCache[rte.pageIdx], nil
}

func (rte *ScrobbleExplorer) RefreshPage() ([]Scrobble, error) {
	tracks, err := rte.sr.FetchScrobblePage(rte.username, rte.pageIdx+1)
	if err != nil {
		return []Scrobble{}, fmt.Errorf("Failed to fetch page %d: %w", rte.pageIdx+1, err)
	}
	rte.pageCache[rte.pageIdx] = tracks
	return tracks, nil
}
