package sources

import (
	gurl "net/url"
	"regexp"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/Sirupsen/logrus"
	"github.com/d3estudio/digest/shared"
)

// Twitter transforms Twitter URLs into embed data
type Twitter struct{}

var twitterReg = regexp.MustCompile(`^https?:\/\/twitter\.com\/(?:#!\/)?(\w+)\/status(es)?\/(\d+)$`)
var twitterEnabled = false

func init() {
	s := shared.Settings
	if s.TwitterKey != "" && s.TwitterSecret != "" {
		twitterEnabled = true
		anaconda.SetConsumerKey(s.TwitterKey)
		anaconda.SetConsumerSecret(s.TwitterSecret)
	}
}

// CanHandle returns whether this instance can handle a certain URL
func (v Twitter) CanHandle(url string) bool {
	return twitterReg.MatchString(url) && twitterEnabled
}

// Process attempts to extract data from a given URL
func (v Twitter) Process(url string) *SourceResult {
	logger := log.WithField("processor", "Twitter")
	api := anaconda.NewTwitterApi("", "")
	o, err := api.GetOEmbed(gurl.Values{
		"url": []string{url},
	})
	if err != nil {
		logger.Error(err)
		return nil
	}
	return &SourceResult{
		Type: "twitter",
		Data: SourceData{
			"html": o.Html,
		},
	}
}
