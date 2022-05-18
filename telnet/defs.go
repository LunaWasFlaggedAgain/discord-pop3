package telnet

import (
	"net"
	"sync"
)

type TelnetServer struct {
	Listener net.Listener
	MsgCh    chan string
	Clients  []*net.Conn
	Mutex    sync.Mutex
}
