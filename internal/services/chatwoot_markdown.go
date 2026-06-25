package services

import (
	"regexp"
	"strings"
)

var (
	chatwootBoldAutolinkRE      = regexp.MustCompile(`\*\*<(https?://[^\s<>]+)>\*\*`)
	chatwootMixedBoldAutolinkRE = regexp.MustCompile(`<\*\*((?:https?://)[^\s*<>]+)>\*\*`)
	chatwootAutolinkRE          = regexp.MustCompile(`<((?:https?://)[^\s<>]+)>`)
	chatwootBoldURLRE           = regexp.MustCompile(`\*\*((?:https?://)[^\s*<>]+)\*\*`)
	chatwootDanglingBoldURLRE   = regexp.MustCompile(`((?:https?://)[^\s*<>]+)\*\*`)
)

func normalizeChatwootMarkdownLinksForCodechat(text string) string {
	out := text
	out = chatwootMixedBoldAutolinkRE.ReplaceAllString(out, "$1")
	out = chatwootBoldAutolinkRE.ReplaceAllString(out, "$1")
	out = chatwootAutolinkRE.ReplaceAllString(out, "$1")
	out = chatwootBoldURLRE.ReplaceAllString(out, "$1")
	out = chatwootDanglingBoldURLRE.ReplaceAllString(out, "$1")

	if out != text {
		return strings.TrimSpace(out)
	}
	return out
}
