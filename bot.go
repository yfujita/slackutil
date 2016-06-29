package slackutil

import (
	"net/http"
	"encoding/json"
	"bytes"
	"time"
	"io/ioutil"
	"errors"
)

type Bot struct {
	url      string
	channel  string
	botName  string
	faceIcon string
}

func NewBot(url, channel, botName, faceIcon string) *Bot {
	bot := new(Bot)
	bot.url = url
	bot.channel = channel
	bot.botName = botName
	bot.faceIcon = faceIcon
	return bot
}

func (bot *Bot) Message(title, text string) error {
	return bot.MessageWithAttachments("", title + "\n" + text, nil)
}

func (bot *Bot) MessageWithAttachments(title, text string, attachments []map[string]string) error {
	textMap := make(map[string]interface{})
	textMap["channel"] = bot.channel
	textMap["username"] = bot.botName
	textMap["icon_emoji"] = bot.faceIcon
	if len(title) > 0 {
		textMap["title"] = title
	}
	if len(text) > 0 {
		textMap["text"] = text
	}
	textMap["link_names"] = 1

	if attachments != nil {
		textMap["attachments"] = attachments
	}

	message, _ := json.Marshal(textMap)

	req, err := http.NewRequest(
		"POST",
		bot.url,
		bytes.NewBuffer([]byte(message)),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(15 * time.Second) }
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(b))
	}
	return nil
}


