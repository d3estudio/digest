package models

import "github.com/d3estudio/digest/shared/slack"

// DigestUser represents an User from the Team account to which the bot belongs
// to. This model is an redux from the original Slack User structure, tailored
// for Digest's needs.
type DigestUser struct {
	RealName string `json:"real_name"`
	Username string `json:"username"`
	Image    string `json:"image"`
	Title    string `json:"title"`
}

// DigestUserFromSlack converts an slack.RTMUser into an DigestUser by
// assigning the RTMUser values into the fields we care about.
func DigestUserFromSlack(u slack.RTMUser) DigestUser {
	return DigestUser{
		RealName: u.Profile.RealName,
		Username: u.Username,
		Image:    u.Profile.ProfilePicture,
		Title:    u.Profile.Title,
	}
}
