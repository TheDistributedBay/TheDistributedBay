package importer

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/crypto"
)

func ProduceTorrents(file string, c chan *core.Torrent, d chan *core.Torrent) {
	log.Println("Reading database dump from:", file)
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer gr.Close()

	cr := csv.NewReader(gr)
	cr.LazyQuotes = true
	cr.Comma = '|'

	for rec, err := cr.Read(); err != io.EOF; rec, err = cr.Read() {
		if err != nil {
			log.Print(err)
			continue
		}

		name := rec[0]
		//created := rec[1]
		magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s", rec[2])

		category := rec[4]
		//seeder := rec[5]
		//leecher := rec[6]

		t := core.CreateTorrent(magnet, name, "from db dump", category, time.Now(), nil)
		if err != nil {
			log.Print(err)
			continue
		}
		c <- t
		d <- t
	}
	close(c)
	close(d)
}

func WriteDbTorrent(db core.Database, c chan *core.Torrent, totalRows int64) {
	start := time.Now()
	count := int64(0)
	for t := range c {
		count++
		db.Add(t)
		if count%1000 == 0 {
			eta := time.Now().Sub(start) / time.Duration(count) * time.Duration(totalRows)
			log.Println("Loaded: ", count, "of", totalRows, "(ETA:", eta.String()+")")
		}
	}
}

func WriteDbSignature(db core.Database, d chan *core.Torrent, totalRows int64) {
	start := time.Now()
	count := int64(0)
	open := true
	b := make([]*core.Torrent, 0)
	k, err := crypto.NewKey()
	if err != nil {
		panic(err)
	}
	for open {
		b = b[:0]
		for i := 0; i < 100 && open; i++ {
			var t *core.Torrent
			t, open = <-d
			b = append(b, t)
		}
		count += int64(len(b))
		s, err := core.SignTorrents(k, b)
		if err != nil {
			panic(err)
		}
		db.AddSignature(s)
		if count%1000 == 0 {
			eta := time.Now().Sub(start) * time.Duration(totalRows) / time.Duration(count)
			log.Println("Signed: ", count, "of", totalRows, "(ETA:", eta.String()+")")
		}
	}
}

func CalculateSize(file string) int64 {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer gr.Close()

	totalRows := int64(0)
	b := make([]byte, 1024)
	var n int
	for n, err = gr.Read(b); err == nil; n, err = gr.Read(b) {
		totalRows += int64(strings.Count(string(b[0:n]), "\n"))
	}
	if err != io.EOF {
		log.Fatal(err)
	}
	return totalRows
}

func Import(file string, db core.Database) {
	c := make(chan *core.Torrent, 2)
	d := make(chan *core.Torrent, 200)
	log.Print("Calculating size")
	totalRows := CalculateSize(file)
	log.Print("Done")
	go ProduceTorrents(file, c, d)
	go WriteDbTorrent(db, c, totalRows)
	go WriteDbSignature(db, d, totalRows)
}
