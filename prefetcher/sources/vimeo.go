package sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

// Vimeo transforms Vimeo URLs into embed data
type Vimeo struct{}

var vimeoReg = regexp.MustCompile(`^https?:\/\/vimeo\.com\/.*$`)

// CanHandle returns whether this instance can handle a certain URL
func (v Vimeo) CanHandle(url string) bool {
	return vimeoReg.MatchString(url)
}

// Process attempts to extract data from a given URL
func (v Vimeo) Process(url string) *SourceResult {
	logger := log.WithField("processor", "Vimeo")
	fields := []string{"description", "title", "html", "thumbnail_height", "thumbnail_width", "thumbnail_url"}
	resp, err := http.Get(fmt.Sprintf("https://vimeo.com/api/oembed.json?url=%s", url))
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
	result := SourceResult{Type: "vimeo", Data: make(SourceData)}
	for _, f := range fields {
		result.Data[f] = httpResult[f]
	}
	return &result
}
