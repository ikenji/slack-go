package slacklog

import (
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type SlackApi struct {
	Ok       string
	Oldest   string
	Messages SlackLog
}

type SlackLog []struct {
	Type        string `json:"type"`
	Subtype     string `json:"subtype"`
	Text        string `json:"text"`
	Ts          string `json:"ts"`
	Username    string `json:"username"`
	BotID       string `json:"bot_id"`
	Attachments []struct {
		Fallback string `json:"fallback"`
		Text     string `json:"text"`
		Title    string `json:"title"`
		ID       int    `json:"id"`
		Color    string `json:"color"`
	} `json:"attachments"`
}

type Raw struct {
	Time      time.Time
	Store     string
	Kind      string
	Mail      string
	Ip        string
	Useragent string
}

type Raws []Raw

func (r Raws) Less(i, j int) bool {
	return r[i].Time.Before(r[j].Time)
}

func Format(slackLog SlackLog) Raws {
	var responses Raws

	for _, slog := range slackLog {
		if slog.Attachments == nil {
			continue
		}

		fallback := slog.Attachments[0].Fallback
		if !strings.Contains(fallback, "Google reCaptcha") {
			continue
		}

		timestamp := slog.Ts
		responses = append(responses, CustomMessages(slog.Attachments[0].Text, timestamp))
	}

	sort.Slice(responses, func(i, j int) bool {
		return responses.Less(i, j)
	})

	return responses
}

func CustomMessages(messages string, timestamp string) Raw {
	mailStringStart := strings.Index(messages, "mailto:") + utf8.RuneCountInString("mailto:")
	mailStringEnd := strings.Index(messages, "|")

	kindStringStart := strings.Index(messages, "kind:") + utf8.RuneCountInString("kind:")
	kindStringEnd := strings.Index(messages, "\nemail:")

	storeStringStart := strings.Index(messages, "store:") + utf8.RuneCountInString("store:")
	storeStringEnd := strings.Index(messages, "\nkind:")

	ipStringStart := strings.Index(messages, "ip:") + utf8.RuneCountInString("ip:")
	ipStringEnd := strings.Index(messages, "\nuser_agent:")

	c := strings.Index(timestamp, ".")
	f, _ := strconv.ParseInt(timestamp[:c], 10, 64)
	time := time.Unix(f, 0)

	var raw Raw
	raw.Time = time
	raw.Store = messages[storeStringStart:storeStringEnd]
	raw.Kind = messages[kindStringStart:kindStringEnd]
	raw.Mail = messages[mailStringStart:mailStringEnd]
	raw.Ip = messages[ipStringStart:ipStringEnd]

	return raw
}
