package importer

import (
	"compress/gzip"
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/crypto"
	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func Import(file string, db database.Database) {
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

	count := 0

	// Signing import with initial key
	ecdsa, err := crypto.NewKey()
	if err != nil {
		log.Print(err)
		return
	}

	rec, err := cr.Read()
	for ; err != io.EOF; rec, err = cr.Read() {
		count += 1
		if err != nil {
			log.Print(err)
			log.Print(rec)
			continue
		}

		name := rec[0]
		//created := rec[1]
		magnet := rec[2]

		category := rec[4]
		//seeder := rec[5]
		//leecher := rec[6]

		t, err := crypto.CreateTorrent(ecdsa, magnet, name, "from db dump", category, time.Now(), nil)
		if err != nil {
			log.Print(err)
			continue
		}

		db.Add(t)

		if count%10 == 0 {
			log.Println("Loaded: ", count)
		}
	}
	log.Println("TOTAL", count)
}
