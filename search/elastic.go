package search

import (
	"errors"
	"strconv"

	elastigo "github.com/mattbaird/elastigo/lib"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type Elastic struct {
	c *elastigo.Conn
	b *elastigo.BulkIndexer
}

func NewElastic(host string) (*Elastic, error) {
	c := elastigo.NewConn()
	c.Domain = host
	b := c.NewBulkIndexer(10)
	b.Start()
	return &Elastic{c, b}, nil
}

func (e *Elastic) NewTorrent(t *core.Torrent) {
	e.c.Index("thedistributedbay", "torrent", t.Hash, nil, t)
	e.c.Flush()
}

func (e *Elastic) NewBatchedTorrent(t *core.Torrent) {
	e.b.Index("thedistributedbay", "torrent", t.Hash, "", nil, t, false)
}

func (e *Elastic) Flush() {
	e.b.Flush()
}

func (e *Elastic) Exists(h string) error {
	resp, err := e.c.ExistsBool("thedistributedbay", "torrent", h, nil)
	if err != nil {
		return err
	}
	if resp {
		return nil
	}
	return errors.New("no such document")
}

func (e *Elastic) Search(term string, from, size int, categories []uint8, sort string) (elastigo.SearchResult, error) {
	query := ""
	if len(term) > 0 {
		query += "(" + term + ")"
	} else {
		query += "(*)"
	}
	if len(categories) > 0 && categories[0] != 0 {
		query += " AND ("
		for i, category := range categories {
			query += "CategoryID:" + strconv.Itoa((int)(category))
			if i != (len(categories) - 1) {
				query += " OR "
			}
		}
		query += ")"
	}
	params := map[string]interface{}{
		"q":    query,
		"from": from,
		"size": size,
		"sort": sort,
	}
	return e.c.SearchUri("thedistributedbay", "torrent", params)
}
