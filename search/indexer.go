package search

import (
	"container/heap"
	"log"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type Indexer struct {
	recentIndexes []*PriorityIndex
}

func NewIndexer(db core.Database) *Indexer {
	recentIndexes := make([]*PriorityIndex, 10)
	for i, _ := range recentIndexes {
		index := make(PriorityIndex, 0)
		heap.Init(&index)
		recentIndexes[i] = &index
	}
	indexer := Indexer{
		recentIndexes,
	}
	go indexer.NewTorrent(db)

	return &indexer
}

func (i *Indexer) Index(t *core.Torrent) {
	i.addItemToCategory(t, 0)
	i.addItemToCategory(t, t.CategoryID)
}

func (i *Indexer) addItemToCategory(t *core.Torrent, category uint8) {
	index := i.recentIndexes[category]
	heap.Push(index, &TorrentIndex{
		torrent:  t,
		priority: t.CreatedAt.Unix(),
	})
	if len(*index) > 350 {
		heap.Pop(index)
	}
}

func (i *Indexer) NewTorrent(db core.Database) {
	c := make(chan *core.Torrent, 100)
	db.AddTorrentClient(c)

	for t := range c {
		i.Index(t)
	}

	log.Println("Indexer torrent channel closed. Reopenning.")

	// If the database closes our connection, reopen it.
	go i.NewTorrent(db)
	return
}

// An TorrentIndex is something we manage in a priority queue.
type TorrentIndex struct {
	torrent  *core.Torrent
	priority int64
	index    int
}

// A PriorityIndex implements heap.Interface and holds TorrentIndexs.
type PriorityIndex []*TorrentIndex

func (pq PriorityIndex) Len() int { return len(pq) }

func (pq PriorityIndex) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityIndex) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityIndex) Push(x interface{}) {
	n := len(*pq)
	item := x.(*TorrentIndex)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityIndex) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an TorrentIndex in the queue.
func (pq *PriorityIndex) update(item *TorrentIndex, torrent *core.Torrent, priority int64) {
	item.torrent = torrent
	item.priority = priority
	heap.Fix(pq, item.index)
}
