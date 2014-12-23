package torrent

import (
  "testing"
)

func BenchmarkScrape(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
    ScrapeTrackers(
      []string{"udp://open.demonii.com:1337/scrape"},
      []string{"A06130D93965BCA27A04CCB9A54CACEB1F5FBCB1"})
	}
}
