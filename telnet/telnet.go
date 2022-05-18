package telnet

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

func New(network, address string) (*TelnetServer, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	return &TelnetServer{
		Listener: l,
		MsgCh:    make(chan string),
		Clients:  make([]*net.Conn, 0),
		Mutex:    sync.Mutex{},
	}, nil
}

func (telnets *TelnetServer) Start() error {
	for {
		conn, err := telnets.Listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			if err := telnets.handleConn(conn); err != nil {
				fmt.Println("telnets.handleConn(conn) err:", err)
			}
		}()
	}
}

func (telnets *TelnetServer) Broadcast(line string) {
	telnets.Mutex.Lock()
	for _, client := range telnets.Clients {
		fmt.Fprintf(*client, "%s\r\n", line)
	}
	telnets.Mutex.Unlock()
}

func (telnets *TelnetServer) handleConn(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	connp := &conn

	telnets.Mutex.Lock()
	telnets.Clients = append(telnets.Clients, connp)
	telnets.Mutex.Unlock()

	defer func() {
		telnets.Mutex.Lock()
		for i, client := range telnets.Clients {
			if client == connp {
				telnets.Clients[i] = telnets.Clients[len(telnets.Clients)-1]
				telnets.Clients = telnets.Clients[:len(telnets.Clients)-1]
				break
			}
		}
		telnets.Mutex.Unlock()

		conn.Close()
	}()

	var line string

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}
		// Echo
		fmt.Fprintf(conn, "%s", string(b))

		if b == '\r' {
			fmt.Fprintf(conn, "\n")

			clean_line := strings.Trim(line, "\r \n")

			fmt.Printf("[DEBUG] telnetd: %s\n", clean_line)
			telnets.MsgCh <- clean_line

			line = ""
		} else if b == 8 {
			if len(line) == 0 {
				continue // can't backspace if there's no characters
			}
			line = line[:len(line)-1]
		} else {
			line += string(b)
		}
	}

	return nil
}
