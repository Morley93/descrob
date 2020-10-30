package descrob

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type LastFMWebClient struct {
	client   http.Client
	username string
	password string
	loginURL string
}

func NewLastFMWebClient(username, password string) (*LastFMWebClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating cookiejar for web client HTTP client: %w", err)
	}
	webClient := LastFMWebClient{
		client:   http.Client{Jar: jar},
		username: username,
		password: password,
		loginURL: "https://secure.last.fm/login",
	}
	err = webClient.startSession()
	if err != nil {
		return nil, fmt.Errorf("Error starting web session: %w", err)
	}
	return &webClient, nil
}

func (c LastFMWebClient) startSession() error {
	_, err := c.client.Get(c.loginURL)
	if err != nil {
		return fmt.Errorf("Failed to begin a LastFM web session: %w", err)
	}

	loginReq, err := c.buildLoginRequest()
	if err != nil {
		return fmt.Errorf("Failed to build login request: %w", err)
	}
	resp, err := c.client.Do(loginReq)
	if err != nil {
		return fmt.Errorf("Error performing LastFM web login: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Got unsuccessful status code (%d) on LastFM login", resp.StatusCode)
	}
	loginCsrfCookie := findCookieForDomain(c.client.Jar, "https://www.last.fm", "csrftoken")
	if loginCsrfCookie == nil {
		return errors.New("Login failed, check your credentials")
	}
	return nil
}

func (c LastFMWebClient) buildLoginRequest() (*http.Request, error) {
	loginForm, err := c.buildLoginForm()
	if err != nil {
		return nil, fmt.Errorf("Failed to create the form required to do login: %w", err)
	}
	formContent := strings.NewReader(loginForm.Encode())
	req, err := http.NewRequest(http.MethodPost, c.loginURL, formContent)
	if err != nil {
		return nil, fmt.Errorf("Error creating a login request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", c.loginURL)
	return req, nil
}

func (c LastFMWebClient) buildLoginForm() (url.Values, error) {
	form := url.Values{}

	csrfCookie := findCookieForDomain(c.client.Jar, "https://secure.last.fm", "csrftoken")
	if csrfCookie == nil {
		return form, errors.New("Failed to find csrf cookie to continue login")
	}

	form.Add("username_or_email", c.username)
	form.Add("password", c.password)
	form.Add("csrfmiddlewaretoken", csrfCookie.Value)
	return form, nil
}

func findCookieForDomain(cookieJar http.CookieJar, domain, name string) *http.Cookie {
	domURL, err := url.Parse(domain)
	if err != nil {
		return nil
	}
	for _, c := range cookieJar.Cookies(domURL) {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (c LastFMWebClient) DeleteTrack(track Track) error {
	deleteURL := fmt.Sprintf("https://www.last.fm/user/%s/library/delete", c.username)

	csrfCookie := findCookieForDomain(c.client.Jar, "https://www.last.fm", "csrftoken")
	if csrfCookie == nil {
		return errors.New("Unable to find CSRF token for delete from established session")
	}

	form := url.Values{}
	form.Add("artist_name", track.Artist.Name)
	form.Add("track_name", track.Name)
	form.Add("timestamp", track.Date.Timestamp)
	form.Add("csrfmiddlewaretoken", csrfCookie.Value)
	form.Add("ajax", "1")

	req, err := http.NewRequest(http.MethodPost, deleteURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("Error creating deletion request: %w", err)
	}
	req.Header.Add("Referer", "https://www.last.fm")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("Request failed: %w", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading delete response body: %w", err)
	}
	ajaxResp := struct {
		Result bool `json:"result"`
	}{}
	err = json.Unmarshal(b, &ajaxResp)
	if err != nil {
		return fmt.Errorf("Failed to ready delete response JSON to check result: %w", err)
	}
	if !ajaxResp.Result {
		return errors.New("Delete response indicate failure")
	}
	return nil
}
