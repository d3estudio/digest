package sources

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"

	gurl "net/url"

	"image/color"
	_ "image/draw" // This is required in order to use image.DecodeConfig
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/victorgama/colorarty"
)

// PoorLink is a last-resort source acquirer that attempts to extract basic
// data from a webpage. Please notice that this source is a catch-all and
// therefore, MUST be placed on the LAST place in the processing list.
type PoorLink struct{}

// CanHandle returns whether this instance can handle a certain URL
func (p PoorLink) CanHandle(url string) bool {
	return true
}

type poorImageData struct {
	URL             string `json:"url"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	Orientation     string `json:"orientation"`
	HasColorData    bool   `json:"has_color_data"`
	BackgroundColor string `json:"background_color"`
	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	DetailColor     string `json:"detail_color"`
}

type poorParserArray []func(*goquery.Document) *string

var poorLogger = log.WithField("processor", "PoorLink")

var processors = map[string]poorParserArray{
	"title": poorParserArray{
		func(doc *goquery.Document) *string {
			item := doc.Find(`meta[property="og:title"]`)
			if item != nil {
				value := item.AttrOr("content", "")
				return &value
			}
			return nil
		},
		func(doc *goquery.Document) *string {
			item := doc.Find(`meta[property="twitter:title"]`)
			if item != nil {
				value := item.AttrOr("content", "")
				return &value
			}
			return nil
		},
		func(doc *goquery.Document) *string {
			item := doc.Find(`meta[itemprop="name"]`)
			if item != nil {
				value := item.AttrOr("content", "")
				return &value
			}
			queries := []*goquery.Selection{
				doc.Find(`title`).First(),
				doc.Find(`h1`).First(),
				doc.Find(`h2`).First(),
				doc.Find(`h3`).First(),
			}
			for _, sel := range queries {
				if sel == nil {
					continue
				}
				value := sel.Text()
				return &value
			}
			return nil
		},
	},
	"summary": poorParserArray{
		func(doc *goquery.Document) *string {
			metas := []string{
				`meta[property="og:description"]`,
				`meta[name="twitter:description"]`,
				`meta[itemprop="description"]`,
				`meta[name="description"]`,
			}
			for _, sel := range metas {
				res := doc.Find(sel).First()
				if res == nil {
					continue
				}
				value := res.AttrOr("content", "")
				if value == "" {
					continue
				}
				return &value
			}
			return nil
		},
	},
	"image": poorParserArray{
		func(doc *goquery.Document) *string {
			metas := []string{
				`meta[property="og:image"]'`,
				`meta[name="twitter:image"]`,
				`meta[itemprop="image"]`,
			}
			for _, sel := range metas {
				res := doc.Find(sel).First()
				if res == nil {
					continue
				}
				value := res.AttrOr("content", "")
				if value != "" {
					return &value
				}
			}
			sel := doc.Find("div img").First()
			if sel == nil {
				return nil
			}
			value := sel.AttrOr("src", "")
			if value == "" {
				return nil
			}
			return &value
		},
	},
}

// Process attempts to extract data from a given URL
func (p PoorLink) Process(url string) *SourceResult {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		poorLogger.Error(err)
		return nil
	}
	results := make(map[string]string)
	for k, v := range processors {
		var result *string
		for _, f := range v {
			result = f(doc)
			if result != nil {
				break
			}
		}
		if result != nil {
			results[k] = *result
		}
	}
	title, hasTitle := results["title"]
	summary, hasSummary := results["summary"]
	image, hasImage := results["image"]
	poorLogger.Debug(results)
	if hasTitle && hasSummary && hasImage {
		baseURL, _ := gurl.Parse(url)
		imgURL, _ := gurl.Parse(image)
		imgURL = baseURL.ResolveReference(imgURL)
		imageData := processImageSize(imgURL.String())
		data := SourceData{
			"title":   title,
			"summary": summary,
			"url":     url,
		}
		if imageData != nil {
			data["imageData"] = imageData
		}
		return &SourceResult{
			Type: "rich-link",
			Data: data,
		}
	} else if hasTitle && (!hasSummary || !hasImage) {
		return &SourceResult{
			Type: "poor-link",
			Data: SourceData{
				"url":   url,
				"title": title,
			},
		}
	}
	return nil
}

func processImageSize(url string) *poorImageData {
	resp, err := http.Get(url)
	if err != nil {
		poorLogger.Error(err)
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		poorLogger.WithField("step", "read_all").Error(err)
		return nil
	}
	conf, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		poorLogger.WithField("step", "decode_config").Error(err)
		return nil
	}

	result := poorImageData{
		Height:      conf.Height,
		Width:       conf.Width,
		Orientation: "vertical",
	}

	if result.Width > result.Height {
		result.Orientation = "horizontal"
	}

	image, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		poorLogger.WithField("step", "full_decode").Error(err)
		return &result
	}

	result.HasColorData = true
	artyResult := colorarty.Analyse(image)
	result.BackgroundColor = goColorToCSS(artyResult.BackgroundColor)
	result.PrimaryColor = goColorToCSS(artyResult.PrimaryColor)
	result.SecondaryColor = goColorToCSS(artyResult.SecondaryColor)
	result.DetailColor = goColorToCSS(artyResult.DetailColor)

	return &result
}

func goColorToCSS(c *color.Color) string {
	cr, cg, cb, _ := (*c).RGBA()
	r := float64(cr)
	g := float64(cg)
	b := float64(cb)
	r /= 0x101
	g /= 0x101
	b /= 0x101
	return fmt.Sprintf("rgba(%.0f, %.0f, %.0f, 1)", r, g, b)
}
