package sources

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"net/http"

	_ "image/draw" // This is required in order to use image.DecodeConfig
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/victorgama/colorarty"
)

// SourceData represents a set of data returned by a Source after processing a
// given URL.
type SourceData map[string]interface{}

// SourceResult represents a SourceResult structure that is stored into the
// database
type SourceResult struct {
	Type string     `json:"type"`
	Data SourceData `json:"data"`
}

// Source defines a set of methods that every Source processor must implement
// in order to be used by the Prefetcher
type Source interface {
	CanHandle(url string) bool
	Process(url string) *SourceResult
}

type imageData struct {
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

func processImageData(url string) *imageData {
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

	result := imageData{
		URL:         url,
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
