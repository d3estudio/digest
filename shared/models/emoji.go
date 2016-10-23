package models

// DigestEmojiChanged represents an event indicating that the custom emoji list
// from the team the bot belongs to has changed.
type DigestEmojiChanged struct {
	Timestamp string `json:"ts"`
}

// Emoji represents an unicode or custom emoji registered on the bot's team
type Emoji struct {
	Aliases []string `json:"aliases"`
	URL     string   `json:"url"`
	Unicode string   `json:"emoji"`
}

// SingleEmoji is a representaton of a single emoji, mapped out from unicode,
// custom emoji URL or aliases
type SingleEmoji struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Unicode string `json:"unicode"`
}
