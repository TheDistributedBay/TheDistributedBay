package torrent

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

// Code to talk to trackers.
// Implements:
//  BEP 12 Multitracker Metadata Extension
//  BEP 15 UDP Tracker Protocol

type InfoHashDetails struct {
	InfoHash                     string
	Seeders, Leechers, Completed uint
}

type InfoHashRange struct {
	InfoHash                     string
	Seeders, Leechers, Completed core.Range
}

type TrackerResponse struct {
	Details []InfoHashDetails
	Error   error
}

func ScrapeTrackers(scrapeList []string, infoHashes []string) ([]*InfoHashRange, error) {
	resp := make(chan TrackerResponse, len(scrapeList))
	for _, tracker := range scrapeList {
		go queryTracker(infoHashes, tracker, resp)
	}
	responseCount := 0
	errorCount := 0
	bestDetails := make([]*InfoHashRange, len(infoHashes))
	for i, hash := range infoHashes {
		bestDetails[i] = &InfoHashRange{InfoHash: hash}
	}
	for trackerResp := range resp {
		if trackerResp.Error != nil {
			errorCount += 1
		} else {
			for i, detail := range trackerResp.Details {
				bestDetail := bestDetails[i]
				bestDetail.Seeders.Update(detail.Seeders)
				bestDetail.Leechers.Update(detail.Leechers)
				bestDetail.Completed.Update(detail.Completed)
			}
		}
		responseCount += 1
		if errorCount >= len(scrapeList) {
			break
		}
		if responseCount >= len(scrapeList) {
			return bestDetails, nil
		}
	}
	return nil, errors.New("Did not successfully contact a tracker.")
}

func queryTracker(infoHashes []string, trackerUrl string, resp chan TrackerResponse) {
	u, err := url.Parse(trackerUrl)
	if err != nil {
		log.Println("Error: Invalid announce URL(", trackerUrl, "):", err)
		resp <- TrackerResponse{nil, err}
		return
	}
	switch u.Scheme {
	case "http":
		fallthrough
	case "https":
		results, err := queryHTTPTracker(infoHashes, u)
		resp <- TrackerResponse{results, err}
	case "udp":
		results, err := queryUDPTracker(infoHashes, u)
		resp <- TrackerResponse{results, err}
	default:
		errorMessage := fmt.Sprintf("Unknown scheme %v in %v", u.Scheme, trackerUrl)
		log.Println(errorMessage)
		resp <- TrackerResponse{nil, errors.New(errorMessage)}
	}
}

func proxyHttpGet(url string) (r *http.Response, e error) {
	return proxyHttpClient().Get(url)
}

func proxyHttpClient() *http.Client {
	tr := &http.Transport{Dial: nil}
	return &http.Client{Transport: tr}
}

func getTrackerInfo(url string) (tr []InfoHashDetails, err error) {
	r, err := proxyHttpGet(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode >= 400 {
		data, _ := ioutil.ReadAll(r.Body)
		reason := "Bad Request " + string(data)
		log.Println(reason)
		err = errors.New(reason)
		return
	}
	//var tr2 InfoHashDetails
	log.Println("HTTP REQ TODO", r.Body)
	/*err = bencode.Unmarshal(r.Body, &tr2)
	r.Body.Close()
	if err != nil {
		return
	}
	tr = &tr2
	*/
	return
}

func queryHTTPTracker(infoHashes []string, u *url.URL) (tr []InfoHashDetails, err error) {
	uq := u.Query()
	for _, infoHash := range infoHashes {
		uq.Add("info_hash", infoHash)
	}

	log.Println("HTTP INFO HASH DEBUG REMOVE THIS", uq)

	// Don't report IPv6 address, the user might prefer to keep
	// that information private when communicating with IPv4 hosts.
	if false {
		ipv6Address, err := findLocalIPV6AddressFor(u.Host)
		if err == nil {
			log.Println("our ipv6", ipv6Address)
			uq.Add("ipv6", ipv6Address)
		}
	}

	// This might reorder the existing query string in the Announce url
	// This might break some broken trackers that don't parse URLs properly.

	u.RawQuery = uq.Encode()

	tr, err = getTrackerInfo(u.String())
	if tr == nil || err != nil {
		log.Println("Error: Could not fetch tracker info:", err)
	}
	return
}

func findLocalIPV6AddressFor(hostAddr string) (local string, err error) {
	// Figure out our IPv6 address to talk to a given host.
	host, hostPort, err := net.SplitHostPort(hostAddr)
	if err != nil {
		host = hostAddr
		hostPort = "1234"
	}
	dummyAddr := net.JoinHostPort(host, hostPort)
	log.Println("Looking for host ", dummyAddr)
	conn, err := net.Dial("udp6", dummyAddr)
	if err != nil {
		log.Println("No IPV6 for host ", host, err)
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr()
	local, _, err = net.SplitHostPort(localAddr.String())
	if err != nil {
		local = localAddr.String()
	}
	return
}

func queryUDPTracker(infoHashes []string, u *url.URL) (tr []InfoHashDetails, err error) {
	serverAddr, err := net.ResolveUDPAddr("udp", u.Host)
	if err != nil {
		return
	}
	con, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return
	}
	defer func() { con.Close() }()

	err = con.SetDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return
	}
	connectionID, err := connectToUDPTracker(con)
	if err != nil {
		return
	}
	return getScrapeFromUDPTracker(con, connectionID, infoHashes)
}

func connectToUDPTracker(con *net.UDPConn) (connectionID uint64, err error) {
	var connectionRequest_connectionID uint64 = 0x41727101980
	var action uint32 = 0
	transactionID := rand.Uint32()

	connectionRequest := new(bytes.Buffer)
	err = binary.Write(connectionRequest, binary.BigEndian, connectionRequest_connectionID)
	if err != nil {
		return
	}
	err = binary.Write(connectionRequest, binary.BigEndian, action)
	if err != nil {
		return
	}
	err = binary.Write(connectionRequest, binary.BigEndian, transactionID)
	if err != nil {
		return
	}

	_, err = con.Write(connectionRequest.Bytes())
	if err != nil {
		return
	}

	connectionResponseBytes := make([]byte, 16)

	var connectionResponseLen int
	connectionResponseLen, err = con.Read(connectionResponseBytes)
	if err != nil {
		return
	}
	if connectionResponseLen != 16 {
		err = fmt.Errorf("Unexpected response size %d", connectionResponseLen)
		return
	}
	connectionResponse := bytes.NewBuffer(connectionResponseBytes)
	var connectionResponseAction uint32
	err = binary.Read(connectionResponse, binary.BigEndian, &connectionResponseAction)
	if err != nil {
		return
	}
	if connectionResponseAction != 0 {
		err = fmt.Errorf("Unexpected response action %d", connectionResponseAction)
		return
	}
	var connectionResponseTransactionID uint32
	err = binary.Read(connectionResponse, binary.BigEndian, &connectionResponseTransactionID)
	if err != nil {
		return
	}
	if connectionResponseTransactionID != transactionID {
		err = fmt.Errorf("Unexpected response transactionID %x != %x",
			connectionResponseTransactionID, transactionID)
		return
	}

	err = binary.Read(connectionResponse, binary.BigEndian, &connectionID)
	if err != nil {
		return
	}
	return
}

func getScrapeFromUDPTracker(con *net.UDPConn, connectionID uint64, infoHashes []string) (tr []InfoHashDetails, err error) {
	transactionID := rand.Uint32()

	announcementRequest := new(bytes.Buffer)
	err = binary.Write(announcementRequest, binary.BigEndian, connectionID)
	if err != nil {
		return
	}
	var action uint32 = 2
	err = binary.Write(announcementRequest, binary.BigEndian, action)
	if err != nil {
		return
	}
	err = binary.Write(announcementRequest, binary.BigEndian, transactionID)
	if err != nil {
		return
	}

	for _, infoHash := range infoHashes {
		var binaryInfoHash []byte
		binaryInfoHash, err = hex.DecodeString(infoHash)
		err = binary.Write(announcementRequest, binary.BigEndian, binaryInfoHash)
		if err != nil {
			return
		}
	}

	_, err = con.Write(announcementRequest.Bytes())
	if err != nil {
		return
	}

	torrentRequestCount := len(infoHashes)

	const minimumResponseLen = 8
	const torrentsDataSize = 12
	expectedResponseLen := minimumResponseLen + torrentsDataSize*torrentRequestCount
	responseBytes := make([]byte, expectedResponseLen)

	var responseLen int
	responseLen, err = con.Read(responseBytes)
	if err != nil {
		return
	}
	if responseLen < minimumResponseLen {
		err = fmt.Errorf("Unexpected response size %d", responseLen)
		return
	}
	response := bytes.NewBuffer(responseBytes)
	var responseAction uint32
	err = binary.Read(response, binary.BigEndian, &responseAction)
	if err != nil {
		return
	}
	if responseAction != action {
		err = fmt.Errorf("Unexpected response action %d", action)
		return
	}
	var responseTransactionID uint32
	err = binary.Read(response, binary.BigEndian, &responseTransactionID)
	if err != nil {
		return
	}
	if transactionID != responseTransactionID {
		err = fmt.Errorf("Unexpected response transactionID %x", responseTransactionID)
		return
	}
	tr = make([]InfoHashDetails, torrentRequestCount)
	for i, infoHash := range infoHashes {
		var seeders uint32
		err = binary.Read(response, binary.BigEndian, &seeders)
		if err != nil {
			return
		}
		var completed uint32
		err = binary.Read(response, binary.BigEndian, &completed)
		if err != nil {
			return
		}
		var leechers uint32
		err = binary.Read(response, binary.BigEndian, &leechers)
		if err != nil {
			return
		}
		tr[i] = InfoHashDetails{
			InfoHash:  infoHash,
			Seeders:   uint(seeders),
			Completed: uint(completed),
			Leechers:  uint(leechers),
		}
	}
	return
}
