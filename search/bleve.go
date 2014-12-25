package search

import (
	"errors"
	"os"

	"github.com/blevesearch/bleve"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type Bleve struct {
	i         bleve.Index
	b         *bleve.Batch
	BatchSize int
}

func NewBleve(f string) (*Bleve, error) {
	mapping := bleve.NewIndexMapping()
	doc := bleve.NewDocumentStaticMapping()
	field := bleve.NewTextFieldMapping()
	field.Store = false
	doc.AddFieldMappingsAt("Name", field)
	doc.AddFieldMappingsAt("Description", field)

	mapping.DefaultMapping = doc
	os.RemoveAll(f)
	index, err := bleve.New(f, mapping)
	if err != nil {
		return nil, err
	}
	b := &Bleve{index, bleve.NewBatch(), 100}
	return b, nil
}

func (b *Bleve) NewTorrent(t *core.Torrent) {
	b.i.Index(t.Hash, t)
}

func (b *Bleve) NewBatchedTorrent(t *core.Torrent) {
	b.b.Index(t.Hash, t)
	if b.b.Size() > b.BatchSize {
		b.i.Batch(b.b)
		b.b = bleve.NewBatch()
	}
}

func (b *Bleve) Flush() {
	b.i.Batch(b.b)
	b.b = bleve.NewBatch()
}

func (b *Bleve) Exists(h string) error {
	p, err := b.i.Document(h)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("no such document")
	}
	return nil
}

func (b *Bleve) Search(term string, from, size int) (*bleve.SearchResult, error) {
	q := bleve.NewQueryStringQuery(term)
	r := bleve.NewSearchRequestOptions(q, size, from, false)
	return b.i.Search(r)
}
