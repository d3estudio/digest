package sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

// YouTube transforms YouTube URLs into embed data
type YouTube struct{}

var youtubeReg = regexp.MustCompile(`^(?:https?:\/\/)?(?:(?:www|m)\.)?(?:youtu\.be\/|youtube(?:-nocookie)?\.com\/(?:embed\/|v\/|watch\?v=|watch\?.+&v=))((\w|-){11})(?:\S+)?$`)

// CanHandle returns whether this instance can handle a certain URL
func (y YouTube) CanHandle(url string) bool {
	return youtubeReg.MatchString(url)
}

// Process attempts to extract data from a given URL
func (y YouTube) Process(url string) *SourceResult {
	logger := log.WithField("processor", "YouTube")
	fields := []string{"title", "html", "thumbnail_height", "thumbnail_width", "thumbnail_url"}
	resp, err := http.Get(fmt.Sprintf("http://www.youtube.com/oembed?url=%s", url))
	if err != nil {
		logger.Error(err)
		return nil
	}
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
	result := SourceResult{Type: "youtube", Data: make(SourceData)}
	for _, f := range fields {
		result.Data[f] = httpResult[f]
	}
	return &result
}
