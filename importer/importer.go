package importer

import (
	"github.com/TheDistributedBay/TheDistributedBay/database"
  "log"
  "os"
  "compress/gzip"
  "encoding/csv"
  "io"
)

func Import(file string, db database.Database) {
  log.Println("Reading database dump from:", file)

  f, err := os.Open(file)
  if err != nil {
      log.Fatal(err)
  }
  defer f.Close()
  log.Println("Decompressed")
  gr, err := gzip.NewReader(f)
  if err != nil {
      log.Fatal(err)
  }
  defer gr.Close()

  cr := csv.NewReader(gr)
  cr.LazyQuotes = true
  cr.Comma = '|'

  count := 0


  rec, err := cr.Read()
  for ;err != io.EOF; rec, err = cr.Read() {
    if err != nil {
        log.Print(err)
        continue
    }
    count += 1
    _ = rec
    if count % 10000 == 0 {
      log.Println(count)
    }
  }
  log.Println(err)
  log.Println("TOTAL", count)
}
