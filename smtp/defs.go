package smtp

import (
	"net"
	"net/mail"
	"regexp"
)

var mailFromRE = regexp.MustCompile(`[Ff][Rr][Oo][Mm]:\s?<(.*)>`)
var rcptToRE = regexp.MustCompile(`[Tt][Oo]:\s?<(.+)>`)

type Email struct {
	*mail.Message
	From string
	To   []string
}

type SMTPServer struct {
	Listener net.Listener
	emailch  chan Email
}
