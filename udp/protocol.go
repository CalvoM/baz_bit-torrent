package udp

import (
	"fmt"
	"net"
	"net/url"
)

type (
	Action uint32
	Event  uint32
)

const (
	Connect Action = iota
	Announce
	Scrape
	Error
)

const (
	None Event = iota
	Completed
	Started
	Stopped
)

type UDPTrackerProtocol struct {
	dialer             net.Dialer
	conn               net.Conn
	trackerURL         *url.URL
	serverConnectionID uint64
}
type UnEqualActionError struct {
	Sent     Action
	Received Action
}

func (e UnEqualActionError) Error() string {
	return fmt.Sprintf("Action Sent: %d, and Action received: %d", e.Sent, e.Received)
}

type UnEqualTransactionIDError struct {
	Sent     TransactionID
	Received TransactionID
}

func (e UnEqualTransactionIDError) Error() string {
	return fmt.Sprintf("Action Sent: %d, and Action received: %d", e.Sent, e.Received)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (udp *UDPTrackerProtocol) ConnectToTracker(possibleURLs []*url.URL) error {
	var payload ConnectRequestPayload
	buf := payload.Build()
	udp.dialer = net.Dialer{}
	var err error
	udp.conn, err = udp.dialer.Dial("udp", possibleURLs[3].Host)
	if err != nil {
		return err
	}
	d, err := udp.conn.Write(buf)
	checkErr(err)
	if d != len(buf) {
		return fmt.Errorf("not all data sent")
	}
	resp := make([]byte, 16)
	d, err = udp.conn.Read(resp)
	checkErr(err)
	if d != len(resp) {
		return fmt.Errorf("not all data received")
	}
	var responsePayload ConnectResponsePayload
	responsePayload.Marshall(resp)
	if payload.transactionID != responsePayload.transactionID {
		return UnEqualTransactionIDError{payload.transactionID, responsePayload.transactionID}
	}
	if payload.action != responsePayload.action {
		return UnEqualActionError{payload.action, responsePayload.action}
	}
	udp.serverConnectionID = responsePayload.connectionID
	return nil
}

func (udp *UDPTrackerProtocol) AnnounceToTracker(infoHash [20]byte, peerID []byte, left uint64) ([]Peer, error) {
	var payload AnnounceRequestPayload
	var p [20]byte
	copy(p[:], peerID)
	buf := payload.Build(udp.serverConnectionID, infoHash, p, left)
	d, err := udp.conn.Write(buf)
	checkErr(err)
	if d != len(buf) {
		return nil, fmt.Errorf("not all data sent")
	}
	// We support 200 peers now.
	resp := make([]byte, announceMinIPV4RespSize+(200*peerSize))
	d, err = udp.conn.Read(resp)
	if d <= announceMinIPV4RespSize {
		return nil, fmt.Errorf("did not get minimum announce response (20) bytes, found: %d", d)
	}
	checkErr(err)
	var responsePayload AnnounceResponsePayload
	responsePayload.Marshall(resp)
	if payload.transactionID != responsePayload.transactionID {
		return nil, UnEqualTransactionIDError{payload.transactionID, responsePayload.transactionID}
	}
	if payload.action != responsePayload.action {
		return nil, UnEqualActionError{payload.action, responsePayload.action}
	}
	peersCount := responsePayload.leechers + responsePayload.seeders
	responsePayload.peers = append(responsePayload.peers, MarshallPeers(resp[20:], int(peersCount))...)
	return responsePayload.peers, nil
}

func (udp *UDPTrackerProtocol) ScrapeTracker(infoHash [][20]byte) error {
	var payload ScrapeRequestPayload
	buf := payload.Build(udp.serverConnectionID, infoHash)
	d, err := udp.conn.Write(buf)
	checkErr(err)
	if d != len(buf) {
		return fmt.Errorf("not all data sent")
	}
	return nil
}
