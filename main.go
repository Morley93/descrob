package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

type Track struct {
	Name   string `json:"name"`
	Artist struct {
		Name string `json:"#text"`
	} `json:"artist"`
	Date struct {
		Timestamp string `json:"uts"`
	} `json:"date"`
}

type RecentTrackResp struct {
	RecentTracks struct {
		Tracks []Track `json:"track"`
	} `json:"recenttracks"`
}

var client *http.Client

func main() {
	username := os.Getenv("LASTFM_USERNAME")
	password := os.Getenv("LASTFM_PASSWORD")
	apiKey := os.Getenv("LASTFM_API_KEY")

	kl, err := os.OpenFile("keys.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	_ = &http.Transport{
		TLSClientConfig: &tls.Config{KeyLogWriter: kl, InsecureSkipVerify: true},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client = &http.Client{Jar: jar}

	recents := getRecentTracks(username, apiKey)

	startSession(username, password)

	doDelete(recents.RecentTracks.Tracks[0], username)
}

func getRecentTracks(username, apiKey string) RecentTrackResp {
	req, err := http.NewRequest(http.MethodGet, "https://ws.audioscrobbler.com/2.0", nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	q.Add("method", "user.getrecenttracks")
	q.Add("user", username)
	q.Add("api_key", apiKey)
	q.Add("format", "json")

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	decoded := RecentTrackResp{}
	err = json.Unmarshal(b, &decoded)
	if err != nil {
		panic(err)
	}
	return decoded
}

func startSession(username, password string) {
	loginURL := "https://secure.last.fm/login"
	_, err := client.Get(loginURL)

	csrfCookie := findCookieForDomain(client.Jar, "https://secure.last.fm", "csrftoken")

	form := url.Values{}
	form.Add("username_or_email", username)
	form.Add("password", password)
	form.Add("csrfmiddlewaretoken", csrfCookie.Value)
	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", loginURL)
	resp, err := client.Do(req)
	fmt.Println(resp)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic("login status " + string(resp.StatusCode))
	}
}

func doDelete(track Track, username string) {
	deleteURL := fmt.Sprintf("https://www.last.fm/user/%s/library/delete", username)

	csrfCookie := findCookieForDomain(client.Jar, "https://www.last.fm", "csrftoken")

	form := url.Values{}
	form.Add("artist_name", track.Artist.Name)
	form.Add("track_name", track.Name)
	form.Add("timestamp", track.Date.Timestamp)
	form.Add("csrfmiddlewaretoken", csrfCookie.Value)
	form.Add("ajax", "1")

	req, err := http.NewRequest(http.MethodPost, deleteURL, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Referer", "https://www.last.fm")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func findCookieForDomain(cookieJar http.CookieJar, domain, name string) *http.Cookie {
	domURL, err := url.Parse(domain)
	if err != nil {
		panic("Domain wouldn't parse")
	}
	for _, c := range cookieJar.Cookies(domURL) {
		if c.Name == name {
			return c
		}
	}
	panic("Couldn't find cookie " + name)
}
