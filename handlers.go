package main

import (
	"./discord"
	"./pop3"
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/quotedprintable"
	"strings"
	"time"
)

func newMsgs(msgs []*discord.Message) {
	AddPop3Messages(msgs)
	AddTelnetMessages(msgs)
}

func toQuotedPrintable(s string) (string, error) {
	var ac bytes.Buffer
	w := quotedprintable.NewWriter(&ac)
	_, err := w.Write([]byte(s))
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	return ac.String(), nil
}

func AddTelnetMessages(msgs []*discord.Message) {
	for _, msg := range msgs {
		telnets.Broadcast(msg.ContentFormattedUsername)
	}
}

func AddPop3Messages(msgs []*discord.Message) {
	pop3s.EmailsMutex.Lock()
	for _, msg := range msgs {
		t, err := time.Parse(time.RFC3339, msg.Timestamp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		message, err := toQuotedPrintable(msg.ContentFormatted)
		if err != nil {
			fmt.Println(err)
			continue
		}

		pop3s.Emails = append(pop3s.Emails, &pop3.Email{
			Headers: pop3.Headers{
				ContentType:             `text/plain; format=flowed; charset="utf-8"`,
				ContentTransferEncoding: "quoted-printable",
				Date:    t.Format("Mon Jan 02 15:04:05 -0700 2006"),
				From:    msg.Author.Username + "#" + msg.Author.Discriminator,
				ReplyTo: fmt.Sprintf("%s+%s@discord.com", msg.ID, msg.ChannelID),
				To:      "luna@localhost",

				Subject: "New Message",
			},

			Message: message,
		})
	}
	pop3s.EmailsMutex.Unlock()
}

func HandleSMTP() {
	for {
		email := smtps.Accept()

		var body string

		if email.Header.Get("Content-Transfer-Encoding") == "quoted-printable" {
			b, err := ioutil.ReadAll(quotedprintable.NewReader(email.Body))
			if err != nil {
				fmt.Println(err)
				continue
			}

			body = string(b)
		} else {
			b, err := ioutil.ReadAll(email.Body)
			if err != nil {
				fmt.Println(err)
				continue
			}

			body = string(b)

		}

		if strings.Contains(body, "----- Original Message -----") {
			split := strings.Split(body, "----- Original Message -----")
			if len(split) < 2 {
				continue
			}

			body = split[0]
		}

		to := email.To[0]
		if strings.HasSuffix(to, "@discord.com") {
			split := strings.Split(to, "@discord.com")
			if len(split) < 2 {
				continue
			}

			var chanid string
			var replyid string

			datasplit := strings.Split(split[0], "+")
			if len(datasplit) < 2 {
				chanid = split[0]
			} else {
				chanid = datasplit[len(datasplit)-1]
				replyid = datasplit[len(datasplit)-2]
			}

			send(body, chanid, replyid)
		}

	}
}

func HandleTelnet() {
	for {
		line := <-telnets.MsgCh
		send(line, CHANNEL, "")
	}
}
