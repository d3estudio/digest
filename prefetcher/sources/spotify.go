package sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

// Spotify transforms Spotify URLs into embed data
type Spotify struct{}

var spotifyRegs = []*regexp.Regexp{
	regexp.MustCompile(`(?:https?:\/\/)?(open|play)\.spotify\.com\/(album|track|user\/[^\/]+\/playlist)\/([a-zA-Z0-9]+)`),
	regexp.MustCompile(`^spotify:(album|track|user:[^:]+:playlist):([a-zA-Z0-9]+)$`),
}

// CanHandle returns whether this instance can handle a certain URL
func (s Spotify) CanHandle(url string) bool {
	for _, r := range spotifyRegs {
		if r.MatchString(url) {
			return true
		}
	}
	return false
}

// Process attempts to extract data from a given URL
func (s Spotify) Process(url string) *SourceResult {
	logger := log.WithField("processor", "Spotify")
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://embed.spotify.com/oembed/?url=%s", url), nil)
	if err != nil {
		logger.Error(err)
		return nil
	}
	req.Header.Set("User-Agent", "request")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}
	var httpResult map[string]interface{}
	err = json.Unmarshal(body, &httpResult)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return &SourceResult{
		Type: "spotify",
		Data: SourceData{
			"html": httpResult["html"],
		},
	}
}
