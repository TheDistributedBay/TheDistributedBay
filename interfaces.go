package TheDistributedBay

import (
	"time"
)

type Torrent struct {
	Id          uint64
	Name        string
	Description string
	CategoryID  uint64
	Size        uint64
	Hash        string
	FileCount   uint64
	CreatedAt   time.Time
	Tags        []string
}

type Database interface {
	Get(id uint64) Torrent
}
