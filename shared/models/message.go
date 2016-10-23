package models

import (
	"time"

	"github.com/d3estudio/digest/shared/slack"
)

// DigestMessage represents an incoming message from a channel the bot is
// member of.
type DigestMessage struct {
	Channel   slack.RTMChannel `json:"channel"`
	User      DigestUser       `json:"user"`
	Text      string           `json:"text"`
	Timestamp string           `json:"ts"`
}

func (m *DigestMessage) Digest() *DigestedMessage {
	return &DigestedMessage{
		Timestamp:   m.Timestamp,
		Text:        m.Text,
		ChannelName: m.Channel.Name,
		Date:        time.Now(),
		User:        m.User.Username,
		Reactions:   make(map[string]int),
	}
}

// DigestMessageDeleted represents an event indicating that a message from
// a channel the bot belongs to has been deleted.
type DigestMessageDeleted struct {
	Timestamp string `json:"ts"`
}

// DigestMessageReaction represents an event indicating that a message from
// a channel the bot belongs to has received or lost a reaction status.
type DigestMessageReaction struct {
	Reaction  string `json:"reaction"`
	Timestamp string `json:"ts"`
	Added     bool   `json:"add"`
}

// DigestMessageChanged represents an event indicating that a message from
// a channel the bot belongs to has changed.s
type DigestMessageChanged struct {
	DigestMessage
}

type DigestedMessage struct {
	Timestamp      string         `json:"ts"`
	Text           string         `json:"text"`
	ChannelName    string         `json:"channel"`
	Reactions      map[string]int `json:"reactions"`
	Date           time.Time      `json:"date"`
	User           string         `json:"user"`
	EmbededContent string         `json:"embed"`
	DetectedURL    string         `json:"url"`
}
