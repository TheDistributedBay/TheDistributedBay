package core

type Database interface {
	Get(hash string) (*Torrent, error)
	Add(t *Torrent) error
	AddSignature(s *Signature)
	AddTorrentClient(c chan *Torrent)
	GetTorrents(c chan string)
}
