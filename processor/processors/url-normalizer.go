package processors

import (
	"regexp"
	"strings"
)

var slackUndupeRegex = regexp.MustCompile(`(?P<beforePipe>https?://[^\|]+)\|(?P<afterPipe>.*)`)

func attemptURLNormalization(input string) string {
	matches := slackUndupeRegex.FindStringSubmatch(input)
	if matches != nil && strings.HasSuffix(matches[1], matches[2]) {
		return matches[1]
	}
	return input
}
