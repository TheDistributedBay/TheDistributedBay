package database

type Torrentdb struct {
	torrents map[string]*Torrent
}

func (t *Torrentdb) Get(hash string) (*Torrent, error) {
	torrent, ok := t.torrents[id]
	if !ok {
		return nil, errors.New("No such hash stored")
	}
	return torrent, nil
}

func (t *Torrentdb) Add(t Torrent) {
	t.torrents[t.hash] = t
}

func (t *Torrentdb) List() []string {
	t := make([]string, 0, len(t.torrents))
	for _, r := range t.torrents {
		t = append(t, r.Hash)
	}
	return t
}
