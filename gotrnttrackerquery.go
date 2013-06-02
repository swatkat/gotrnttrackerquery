package gotrnttrackerquery

import (
	"code.google.com/p/bencode-go"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TrackerResponse struct {
	FailureReason  string "failure reason"
	WarningMessage string "warning message"
	Interval       int64  "interval"
	MinInterval    int64  "min interval"
	TrackerId      string "tracker id"
	Complete       int    "complete"
	Incomplete     int    "incomplete"
	Peers          string "peers"
}

// Build tracker request URL from announce and other params
func (trackerResp *TrackerResponse) GetTrackerInfo(annouceUrl, infoHash, peerId string, port uint64) bool {
	trackerUrl, ok := trackerResp.BuildTrackerRequestUrl(annouceUrl, infoHash, peerId, port)
	if !ok {
		return false
	}

	// Do a HTTP GET on tracker request URL
	httpResp, er := http.Get(trackerUrl)
	if er != nil {
		return false
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return false
	}

	// Unmarshal bencoded HTTP response
	er = bencode.Unmarshal(httpResp.Body, trackerResp)
	if er != nil {
		return false
	}

	return true
}

// Build a list of strings in ip:port format
func (trackerResp *TrackerResponse) GetIpPortListFromPeers() []string {
	var peerIpPortList []string
	peersLen := len(trackerResp.Peers)
	for i, j := 0, 0; i < peersLen; i, j = i+6, j+1 {
		peerIpPort := fmt.Sprintf("%d.%d.%d.%d:%d",
			trackerResp.Peers[i], trackerResp.Peers[i+1],
			trackerResp.Peers[i+2], trackerResp.Peers[i+3],
			uint16(trackerResp.Peers[i+4])<<8|uint16(trackerResp.Peers[i+5]))
		peerIpPortList = append(peerIpPortList, peerIpPort)
	}
	return peerIpPortList
}

func (trackerResp *TrackerResponse) DumpTrackerResponse() {
	fmt.Println("Failure reason:", trackerResp.FailureReason)
	fmt.Println("Warning message:", trackerResp.WarningMessage)
	fmt.Println("Interval:", trackerResp.Interval)
	fmt.Println("Min interval:", trackerResp.MinInterval)
	fmt.Println("Tracker id:", trackerResp.TrackerId)
	fmt.Println("Complete:", trackerResp.Complete)
	fmt.Println("Incomplete:", trackerResp.Incomplete)
	peerIpPortList := trackerResp.GetIpPortListFromPeers()
	fmt.Println("Peers:", peerIpPortList)
}

func (trackerResp *TrackerResponse) BuildTrackerRequestUrl(announceUrl, infoHash, peerId string,
	port uint64) (string, bool) {
	parsedAnnounceUrl, er := url.Parse(announceUrl)
	if er != nil {
		return "", false
	}
	queryUrl := parsedAnnounceUrl.Query()
	queryUrl.Add("info_hash", infoHash)
	queryUrl.Add("peer_id", peerId)
	queryUrl.Add("port", strconv.FormatUint(port, 10))
	queryUrl.Add("uploaded", strconv.FormatInt(0, 10))
	queryUrl.Add("downloaded", strconv.FormatInt(0, 10))
	queryUrl.Add("left", strconv.FormatInt(0, 10))
	queryUrl.Add("compact", "1")
	parsedAnnounceUrl.RawQuery = queryUrl.Encode()
	return parsedAnnounceUrl.String(), true
}
