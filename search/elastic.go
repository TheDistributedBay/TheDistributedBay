package search

import (
	"errors"
  "fmt"
  "encoding/json"

	elastigo "github.com/mattbaird/elastigo/lib"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type Elastic struct {
	c     *elastigo.Conn
	b     *elastigo.BulkIndexer
	index string
}

// Initializes a new elastic search handler.
func NewElastic(host string, index string) (*Elastic, error) {
	c := elastigo.NewConn()
	c.Domain = host
	b := c.NewBulkIndexer(10)
	b.Start()
	return &Elastic{c, b, index}, nil
}

// Indexes a torrent and flushes it.
func (e *Elastic) NewTorrent(t *core.Torrent) {
	e.c.Index(e.index, "torrent", t.Hash, nil, t)
	e.c.Flush()
}

// Indexes a torrent with the bulk indexer
func (e *Elastic) NewBatchedTorrent(t *core.Torrent) {
	e.b.Index(e.index, "torrent", t.Hash, "", nil, t, false)
}

// Flushes the bulk indexer
func (e *Elastic) Flush() {
	e.b.Flush()
}

// Checks if a hash is in the database, returns an error if it's not.
func (e *Elastic) Exists(h string) error {
	resp, err := e.c.ExistsBool(e.index, "torrent", h, nil)
	if err != nil {
		return err
	}
	if resp {
		return nil
	}
	return errors.New("no such document")
}

// Performs a search for torrents with the relevant options.
func (e *Elastic) Search(term string, from, size int, categories []uint8, sort string) (*elastigo.Hits, error) {
	// We can't use the DSL since we're using simple_query_string to avoid errors.
	args := map[string]interface{}{
		"from": from,
		"size": size,
	}
	if len(sort) > 0 {
		args["sort"] = sort
	}
	data := map[string]interface{}{}
	if len(term) > 0 {
		data["query"] = map[string]interface{}{
			"simple_query_string": map[string]interface{}{
				"query": term,
			},
		}
	}
	if len(categories) > 0 && categories[0] != 0 {
		categoryInterface := make([]interface{}, len(categories))
		for i, category := range categories {
			categoryInterface[i] = category
		}
		data["filter"] = map[string]interface{}{
			"terms": map[string]interface{}{
				"CategoryID": categoryInterface,
			},
		}
	}
	results, err := e.c.Search(e.index, "torrent", args, data)
	if err != nil {
		return nil, err
	}
	return &results.Hits, nil
}


func (e *Elastic) MoreLikeThis(hash string) (*elastigo.Hits, error) {
  uriVal := fmt.Sprintf("/%s/torrent/%s/_mlt", e.index, hash)
	var retval elastigo.SearchResult
	body, err := e.c.DoCommand("GET", uriVal, map[string]interface{}{
    "min_term_freq": 1,
    "mlt_fields": "Name,Description",
  }, map[string]interface{}{})
	if err != nil {
		return &retval.Hits, err
	}
	if err == nil {
		// marshall into json
		jsonErr := json.Unmarshal([]byte(body), &retval)
		if jsonErr != nil {
			return &retval.Hits, jsonErr
		}
	}
	retval.RawJSON = body
	return &retval.Hits, err
}
