package pop3

import (
	"bufio"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func BuildHeaders(email *Email) string {
	headersArr := make([]string, 0)

	if email.RawHeaders != "" {
		headersArr = append(headersArr, email.RawHeaders+"\r\n")
	}

	t := reflect.TypeOf(email.Headers)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("pop3")
		value := reflect.ValueOf(email.Headers).Field(i).String()

		if value == "" || tag == "-" {
			continue
		}

		headersArr = append(headersArr, fmt.Sprintf("%s: %s", tag, value))
	}

	return strings.Join(headersArr, "\r\n")
}

func getLength(email *Email) int {
	return len(BuildHeaders(email)) + len(email.Message)
}

func getLengthAll(emails []*Email) int {
	length := 0

	for _, email := range emails {
		length += getLength(email)
	}

	return length
}

func New(network, address string) (*Pop3Server, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	return &Pop3Server{
		Emails:      make([]*Email, 0),
		Listener:    l,
		EmailsMutex: sync.Mutex{},
	}, nil
}

func (pop3s *Pop3Server) Start() error {
	for {
		conn, err := pop3s.Listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			if err := pop3s.handleConn(conn); err != nil {
				fmt.Println("pop3s.handleConn(conn) err:", err)
			}
		}()
	}
}

func (pop3s *Pop3Server) handleConn(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	shouldDelete := make([]int, 0)

	fmt.Fprintf(conn, "+OK luna's go pop3 server\r\n")

	defer func() {
		pop3s.EmailsMutex.Lock()
		for _, id := range shouldDelete {
			pop3s.Emails[id] = nil
		}

		for i := 0; i < len(pop3s.Emails); {
			if pop3s.Emails[i] != nil {
				i++
				continue
			}

			if i < len(pop3s.Emails)-1 {
				copy(pop3s.Emails[i:], pop3s.Emails[i+1:])
			}

			pop3s.Emails[len(pop3s.Emails)-1] = nil
			pop3s.Emails = pop3s.Emails[:len(pop3s.Emails)-1]
		}
		pop3s.EmailsMutex.Unlock()
	}()

	for {
		raw_line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		fmt.Printf("[DEBUG] pop3d: %s", raw_line)

		line := strings.Trim(raw_line, "\r \n")
		split := strings.Split(line, " ")

		cmd := split[0]
		args := split[1:]

		switch cmd {
		case "CAPA":
			fmt.Fprintf(conn, "+OK Capability list follows\r\n")
			fmt.Fprintf(conn, "TOP\r\n")
			fmt.Fprintf(conn, "USER\r\n")
			fmt.Fprintf(conn, "UIDL\r\n")
			fmt.Fprintf(conn, "IMPLEMENTATION sex\r\n")
			fmt.Fprintf(conn, ".\r\n")
		case "USER":
			fmt.Fprintf(conn, "+OK User accepted\r\n")
		case "PASS":
			fmt.Fprintf(conn, "+OK Pass accepted\r\n")
		case "STAT":
			pop3s.EmailsMutex.Lock()
			fmt.Fprintf(conn, "+OK %d %d\r\n", len(pop3s.Emails), getLengthAll(pop3s.Emails))
			pop3s.EmailsMutex.Unlock()
		case "LIST":
			// Print all messages
			pop3s.EmailsMutex.Lock()
			fmt.Fprintf(conn, "+OK %d messages (%d octets)\r\n", len(pop3s.Emails), getLengthAll(pop3s.Emails))
			for i, email := range pop3s.Emails {
				fmt.Fprintf(conn, "%d %d\r\n", i+1, getLength(email))
			}
			pop3s.EmailsMutex.Unlock()
			// Ending
			fmt.Fprintf(conn, ".\r\n")
		case "UIDL":
			if len(args) < 1 {
				pop3s.EmailsMutex.Lock()
				for i := range pop3s.Emails {
					fmt.Fprintf(conn, "+OK %d %d\r\n", i+1, i+1)
				}
				pop3s.EmailsMutex.Unlock()
				continue
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			pop3s.EmailsMutex.Lock()
			if len(pop3s.Emails) <= id-1 {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}
			pop3s.EmailsMutex.Unlock()

			fmt.Fprintf(conn, "+OK %d %d\r\n", id, id)
		case "TOP":
			if len(args) < 1 {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			pop3s.EmailsMutex.Lock()
			if len(pop3s.Emails) <= id-1 {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			email := pop3s.Emails[id-1]
			pop3s.EmailsMutex.Unlock()

			var lines int = 0
			lines, _ = strconv.Atoi(args[1])

			split := strings.Split(email.Message, "\r\n")

			if lines > len(split) {
				lines = len(split)
			}

			fmt.Fprintf(conn, "+OK Top of message follows\r\n")
			fmt.Fprintf(conn, "%s\r\n\r\n%s\r\n.\r\n", BuildHeaders(email), strings.Join(split[:lines], "\r\n"))
		case "RETR":
			if len(args) < 1 {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			pop3s.EmailsMutex.Lock()
			if len(pop3s.Emails) <= id-1 {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			email := pop3s.Emails[id-1]
			pop3s.EmailsMutex.Unlock()

			fmt.Fprintf(conn, "+OK\r\n")
			fmt.Fprintf(conn, "%s\r\n\r\n%s\r\n.\r\n", BuildHeaders(email), email.Message)
		case "DELE":
			if len(args) < 1 {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}

			pop3s.EmailsMutex.Lock()
			if len(pop3s.Emails) <= id-1 {
				fmt.Fprintf(conn, "-ERR no such message\r\n")
				continue
			}
			pop3s.EmailsMutex.Unlock()

			shouldDelete = append(shouldDelete, id-1)

			fmt.Fprintf(conn, "+OK\r\n")
		case "QUIT":
			return nil
		default:
			fmt.Fprintf(conn, "-ERR Command not implemented\r\n")
		}
	}

	return nil
}
