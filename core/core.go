package core

type TorrentWriter interface {
	NewTorrent(t *Torrent)
	NewSignature(s *Signature)
}

type Database interface {
	Get(hash string) (*Torrent, error)
	Add(t *Torrent)
	AddSignature(s *Signature)
	List() []string
	AddClient(w TorrentWriter)
}
