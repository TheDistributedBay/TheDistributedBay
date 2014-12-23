package torrent

import (
	"testing"
)

func TestTorrentInfoRetrevial(t *testing.T) {
	infoHashes := []string{
		"A06131D93965BCA27A04CCB9A54CACEB1F5FBCB1",
		"18A025ADEE8D14F468A7749AD8FB751CBA1A788C",
		"D18B786E32D43705AD37DF09A10A0F5B9173111E",
		"C4477073A96A9BED8B739E1BAE477731282243C3",
		"2552169E161A1E766AA1DDC5F03B16AC5CB50F68",
		"11A2AC68A11634E980F265CB1433C599D017A759",
		"E47AE3BF2FA5BE8DFD620868BC812BE348DA1C92",
		"205F03DB95617F7EAC3E9ED4415BB89FC6E362A8",
		"D2310F718EB02F98665266786F7D00B42A20F055",
		"7CFA1BE24072795701386ABA248AD5E26C7F18AE",
	}
	details, err := ScrapeTrackers(
		[]string{"udp://open.demonii.com:1337/scrape"},
		infoHashes)
	if err != nil {
		t.Fatal(err)
	}
	if len(details) != len(infoHashes) {
		t.Fatal("Outputted info hash details is not the same length.")
	}
	for i, info := range details {
		if info.InfoHash != infoHashes[i] {
			t.Fatal("Produced info hash does not match inputted one.")
		}
		if info.Seeders == 0 && info.Leechers == 0 && info.Completed == 0 {
			t.Fatal("Returned details has zero values.", info, details)
		}
	}
}
func TestTorrentInfoFallback(t *testing.T) {
	infoHashes := []string{"A06130D93965BCA27A04CCB9A54CACEB1F5FBCB1"}
	details, err := ScrapeTrackers([]string{
		"udp://test.com:80/scrape",
		"udp://open.demonii.com:1337/scrape",
	}, infoHashes)
	if err != nil {
		t.Fatal(err)
	}
	if len(details) != len(infoHashes) {
		t.Fatal("Outputted info hash details is not the same length.")
	}
	if details[0].InfoHash != infoHashes[0] {
		t.Fatal("Produced info hash does not match inputted one.")
	}
	if details[0].Seeders == 0 || details[0].Leechers == 0 || details[0].Completed == 0 {
		t.Fatal("Returned details has a zero value.")
	}
}

func BenchmarkScrape(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ScrapeTrackers(
			[]string{"udp://open.demonii.com:1337/scrape"},
			[]string{"A06130D93965BCA27A04CCB9A54CACEB1F5FBCB1"})
	}
}
