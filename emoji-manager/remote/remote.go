package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/d3estudio/digest/shared"
	"github.com/d3estudio/digest/shared/models"
)

type emojiList struct {
	Ok    bool              `json:"ok"`
	Emoji map[string]string `json:"emoji"`
	Cache string            `emoji:"cache_ts"`
}

// ErrRemoteNotOkay indicates that the remote server could not process the request
// at the moment. Try again later.
var ErrRemoteNotOkay = errors.New("Remote reported an internal error. Try again later.")

// BuildEmojiDatabase takes an officialData input that should contain official
// unicode mappings and augments it with data obtained from the slack team
func BuildEmojiDatabase(officialData []byte) (arr []models.Emoji, err error) {
	err = json.Unmarshal(officialData, &arr)
	if err != nil {
		return
	}
	url := fmt.Sprintf("https://slack.com/api/emoji.list?token=%s", shared.Settings.Token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result emojiList
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	if !result.Ok {
		err = ErrRemoteNotOkay
		return
	}

	aliases := make(map[string]string)
	for k, v := range result.Emoji {
		if strings.HasPrefix(v, "alias:") {
			aliases[strings.Split(v, ":")[1]] = k
			continue
		}
		arr = append(arr, models.Emoji{
			Aliases: []string{k},
			URL:     v,
		})
	}

	for k, v := range aliases {
	ItemLoop:
		for _, e := range arr {
			for _, a := range e.Aliases {
				if a == k {
					e.Aliases = append(e.Aliases, v)
					break ItemLoop
				}
			}
		}
	}
	return
}
