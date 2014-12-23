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

	"github.com/TheDistributedBay/TheDistributedBay/crypto"
	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func ProduceTorrents(file string, c chan *database.Torrent) {
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

	// Signing import with random initial key
	ecdsa, err := crypto.NewKey()
	if err != nil {
		log.Print(err)
		return
	}

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

		t, err := crypto.CreateTorrent(ecdsa, magnet, name, "from db dump", category, time.Now(), nil)
		if err != nil {
			log.Print(err)
			continue
		}
		c <- t
	}
	close(c)
}

func WriteDb(db database.Database, c chan *database.Torrent, totalRows int64) {
	start := time.Now()
	count := int64(0)
	for t := range c {
		count++
		db.Add(t)
		if count%1000 == 0 {
			eta := time.Now().Sub(start) * time.Duration(totalRows) / time.Duration(count)
			log.Println("Loaded: ", count, "of", totalRows, "(ETA:", eta.String()+")")
		}
	}
	log.Println("TOTAL", count)
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
	for n, err := f.Read(b); err != io.EOF; n, err = f.Read(b) {
		totalRows += int64(strings.Count(string(b[0:n]), "\n"))
	}
	return totalRows
}

func Import(file string, db database.Database) {
	c := make(chan *database.Torrent, 2)
	log.Print("Calculating size")
	totalRows := CalculateSize(file)
	log.Print("Done")
	go ProduceTorrents(file, c)
	go WriteDb(db, c, totalRows)
}
