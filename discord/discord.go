package discord

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
)

var client = &http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}}

type MessageReference struct {
	ID string `json:"message_id,omitempty"`
}

type Message struct {
	Author *struct {
		Username      string `json:"username,omitempty"`
		Discriminator string `json:"discriminator,omitempty"`
	} `json:"author,omitempty"`

	ID        string `json:"id,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`

	MessageReference  *MessageReference `json:"message_reference,omitempty"`
	ReferencedMessage *Message          `json:"referenced_message"`

	Content                  string `json:"content"`
	ContentFormatted         string `json:"-"`
	ContentFormattedUsername string `json:"-"`
}

func FetchMessages(token, channel, after string) ([]*Message, error) {
	url := "https://discord.com/api/v9/channels/" + channel + "/messages?limit=10"
	if after != "" {
		url += "&after=" + after
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var msgs []*Message
	err = json.NewDecoder(res.Body).Decode(&msgs)
	if err != nil {
		return msgs, err
	}

	for _, msg := range msgs {
		msg.ContentFormatted = FormatMsg(msg, false)
		msg.ContentFormattedUsername = FormatMsg(msg, true)
	}

	return msgs, nil
}

func SendMessage(token, channel string, msg Message) (Message, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return Message{}, err
	}

	req, err := http.NewRequest("POST", "https://discord.com/api/v9/channels/"+channel+"/messages", bytes.NewReader(b))
	if err != nil {
		return Message{}, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return Message{}, err
	}

	var discordMsg Message
	err = json.NewDecoder(res.Body).Decode(&discordMsg)

	return discordMsg, err
}
