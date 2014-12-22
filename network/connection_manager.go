package network

import (
	"io"
	"log"
	"net"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

type ConnectionManager struct {
	db  database.Database
	chs []io.Closer
}

func NewConnectionManager(db database.Database) *ConnectionManager {
	m := &ConnectionManager{db, nil}
	return m
}

func (m *ConnectionManager) Listen(l net.Listener) {
	m.chs = append(m.chs, l)
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection on %v : %v", l, err)
			return
		}
		m.Handle(c)
	}
}

func (m *ConnectionManager) Handle(c net.Conn) {
	t := NewTranscoder(c)
	ch := NewConnectionHandler(t, m.db)
	m.chs = append(m.chs, ch)
}

func (m *ConnectionManager) Close() error {
	for _, c := range m.chs {
		err := c.Close()
		log.Printf("Closing %v got %v", c, err)
	}
	return nil
}
