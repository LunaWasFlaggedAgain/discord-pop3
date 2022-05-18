package discord

import "fmt"

func FormatMsg(msg *Message, usernamePrefix bool) string {
	content := ""

	if msg.ReferencedMessage != nil {
		content += fmt.Sprintf("┌ %s#%s: %s\r\n", msg.ReferencedMessage.Author.Username, msg.ReferencedMessage.Author.Discriminator, msg.ReferencedMessage.Content)
	}

	if usernamePrefix {
		content += fmt.Sprintf("%s#%s: %s", msg.Author.Username, msg.Author.Discriminator, msg.Content)
	} else {
		content += msg.Content
	}

	return content
}
