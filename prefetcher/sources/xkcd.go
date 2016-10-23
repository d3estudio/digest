package sources

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

// XKCD transforms XKCD URLs into embed data
type XKCD struct{}

var xkcdReg = regexp.MustCompile(`^(https?:\/\/)?(www\.)?xkcd\.com\/(\d+)\/?$`)

// CanHandle returns whether this instance can handle a certain URL
func (x XKCD) CanHandle(url string) bool {
	return xkcdReg.MatchString(url)
}

// Process attempts to extract data from a given URL
func (x XKCD) Process(url string) *SourceResult {
	logger := log.WithField("processor", "XKCD")
	doc, err := goquery.NewDocument(url)
	if err != nil {
		logger.Error(err)
		return nil
	}
	item := doc.Find("#comic > img").First()
	if item == nil {
		logger.Error(err)
		return nil
	}
	return &SourceResult{
		Type: "xkcd",
		Data: SourceData{
			"img":     item.AttrOr("src", ""),
			"title":   item.AttrOr("title", ""),
			"explain": item.AttrOr("title", ""),
			"link":    url,
		},
	}
}
