package smtp

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/mail"
	"strings"
)

func New(network, address string) (*SMTPServer, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	return &SMTPServer{
		Listener: l,
		emailch:  make(chan Email),
	}, nil
}

func (smtps *SMTPServer) Start() error {
	for {
		conn, err := smtps.Listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			if err := smtps.handleConn(conn); err != nil {
				fmt.Println("smtps.handleConn(conn) err:", err)
			}
		}()
	}
}

func (smtps *SMTPServer) Accept() Email {
	return <-smtps.emailch
}

func (smtps *SMTPServer) readData(reader *bufio.Reader) ([]byte, error) {
	var data []byte

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		if bytes.Equal(line, []byte(".\r\n")) {
			break
		}

		if line[0] == '.' {
			line = line[1:]
		}

		data = append(data, line...)
	}

	return data, nil
}

func (smtps *SMTPServer) handleConn(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	fmt.Fprintf(conn, "220 localhost LunasGoSMTPServer ESMTP Service ready\r\n")

	var from string
	var to []string

	for {
		raw_line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		fmt.Printf("[DEBUG] smtpd: %s", raw_line)

		line := strings.Trim(raw_line, "\r \n")
		split := strings.Split(line, " ")

		cmd := split[0]
		args := strings.Join(split[1:], " ")

		switch cmd {
		case "HELO":
			fmt.Fprintf(conn, "250 localhost greets %s\r\n", args)
		case "EHLO":
			fmt.Fprintf(conn, "250-localhost greets %s\r\n", args)
			fmt.Fprintf(conn, "250-SIZE 0\r\n")
			fmt.Fprintf(conn, "250-AUTH PLAIN\r\n") // fuck u outlook
			fmt.Fprintf(conn, "250 ENHANCEDSTATUSCODES\r\n")
		case "AUTH":
			if args != "" { // "Read" username
				fmt.Fprintf(conn, "334 \r\n")
				_, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
			}

			fmt.Fprintf(conn, "235 2.7.0 Authentication successful\r\n")
		case "MAIL":
			match := mailFromRE.FindStringSubmatch(args)
			if match == nil {
				fmt.Fprintf(conn, "501 5.5.4 Syntax error in parameters or arguments (invalid FROM parameter)\r\n")
				continue
			}

			from = match[1]
			fmt.Fprintf(conn, "250 2.1.0 Ok\r\n")
		case "RCPT":
			if from == "" {
				fmt.Fprintf(conn, "503 5.5.1 Bad sequence of commands (MAIL required before RCPT)\r\n")
				continue
			}

			match := rcptToRE.FindStringSubmatch(args)
			if match == nil {
				fmt.Fprintf(conn, "501 5.5.4 Syntax error in parameters or arguments (invalid TO parameter)\r\n")
				continue
			}

			if len(to) > 100 {
				fmt.Fprintf(conn, "452 4.5.3 Too many recipients\r\n")
				continue
			}

			to = append(to, match[1])
			fmt.Fprintf(conn, "250 2.1.5 Ok\r\n")
		case "DATA":
			if from == "" || len(to) == 0 {
				fmt.Fprintf(conn, "503 5.5.1 Bad sequence of commands (MAIL & RCPT required before DATA)\r\n")
				continue
			}

			fmt.Fprintf(conn, "354 Start mail input, end with <CR><LF>.<CR><LF>\r\n")

			data, err := smtps.readData(reader)
			if err != nil {
				return err
			}

			msg, err := mail.ReadMessage(bytes.NewReader(data))
			if err != nil {
				return err
			}

			smtps.emailch <- Email{
				Message: msg,
				From:    from,
				To:      to,
			}

			fmt.Fprintf(conn, "250 2.0.0 Ok: queued\r\n")

		case "QUIT":
			fmt.Fprintf(conn, "221 2.0.0 localhost LunasGoSMTPServer ESMTP Service closing transmission channel\r\n")
			return nil
		default:
			fmt.Fprintf(conn, "500 5.5.2 Syntax error, command unrecognized\r\n")
		}
	}

	return nil
}
