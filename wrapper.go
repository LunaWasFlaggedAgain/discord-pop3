package main

import (
	"./discord"
	"fmt"
	"time"
)

var after string

func discordRecvLoop() {
	for {
		time.Sleep(1 * time.Second)

		msgs, err := discord.FetchMessages(TOKEN, CHANNEL, after)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(msgs) == 0 {
			continue
		}

		after = msgs[0].ID

		newMsgs(msgs)
	}
}

func send(content, chanid, replyid string) {
	msg := discord.Message{Content: content}
	if replyid != "" {
		msg.MessageReference = &discord.MessageReference{
			ID: replyid,
		}
	}
	_, err := discord.SendMessage(TOKEN, chanid, msg)

	if err != nil {
		fmt.Println(err)
	}
}
