package main

import (
	"./pop3"
	"./smtp"
	"./telnet"
)

var pop3s *pop3.Pop3Server
var smtps *smtp.SMTPServer
var telnets *telnet.TelnetServer

func main() {
	var err error
	pop3s, err = pop3.New("tcp4", "0.0.0.0:110")
	if err != nil {
		panic(err)
	}

	smtps, err = smtp.New("tcp4", "0.0.0.0:25")
	if err != nil {
		panic(err)
	}

	telnets, err = telnet.New("tcp4", "0.0.0.0:23")
	if err != nil {
		panic(err)
	}

	go func() {
		panic(pop3s.Start())
	}()

	go func() {
		panic(smtps.Start())
	}()

	go func() {
		panic(telnets.Start())
	}()

	go HandleSMTP()
	go HandleTelnet()
	discordRecvLoop()
}
