package search

import (
	"os"

	"github.com/blevesearch/bleve"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

type Bleve struct {
	i         bleve.Index
	b         *bleve.Batch
	BatchSize int
}

func NewBleve(f string) (*Bleve, error) {
	mapping := bleve.NewIndexMapping()
	doc := bleve.NewDocumentMapping()
	ignore := bleve.NewTextFieldMapping()
	ignore.Index = false
	ignore.Store = false
	doc.AddFieldMappingsAt("Hash", ignore)
	doc.AddFieldMappingsAt("PublicKey", ignore)
	doc.AddFieldMappingsAt("Signature", ignore)
	mapping.DefaultMapping = doc
	os.RemoveAll("search.bleve")
	index, err := bleve.New("search.bleve", mapping)
	if err != nil {
		return nil, err
	}
	b := &Bleve{index, bleve.NewBatch(), 100}
	return b, nil
}

func (b *Bleve) NewTorrent(t *database.Torrent) {
	b.i.Index(t.Hash, t)
}

func (b *Bleve) NewBatchedTorrent(t *database.Torrent) {
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

func (b *Bleve) Search(term string, from int, size int) (*bleve.SearchResult, error) {
	q := bleve.NewQueryStringQuery(term)
	r := bleve.NewSearchRequestOptions(q, size, from, false)
	return b.i.Search(r)
}
