package network

import (
	"io"
	"log"
	"net"

	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/TheDistributedBay/TheDistributedBay/tls"
)

type ConnectionManager struct {
	db  database.Database
	chs []io.Closer
	l   net.Listener
}

type Connection interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (int, error)
	Close() error
	Protocol() string
}

func NewConnectionManager(db database.Database) *ConnectionManager {
	return &ConnectionManager{db, nil, nil}
}

func (m *ConnectionManager) NumPeers() int {
	return len(m.chs)
}

func (m *ConnectionManager) Listen(l net.Listener) {
	m.l = l
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection on %v : %v", l, err)
			return
		}
		m.Handle(tls.Wrap(c))
	}
}

func (m *ConnectionManager) Handle(c Connection) {
	if c.Protocol() != tls.Proto {
		log.Printf("Unrecognized proto on %v : %v", c, c.Protocol())
	}
	t := NewTranscoder(c)
	ch := NewConnectionHandler(t, m.db)
	m.chs = append(m.chs, ch)
}

func (m *ConnectionManager) Close() error {
	for _, c := range m.chs {
		err := c.Close()
		log.Printf("Closing %v got %v", c, err)
	}
	m.l.Close()
	return nil
}
