package udp

import (
	"fmt"
	"net"
	"net/url"
)

type Action uint32

const (
	Connect Action = iota
	Announce
	Scrape
	Error
)

type UDPTrackerProtocol struct {
	connDialer         net.Dialer
	trackerURL         *url.URL
	serverConnectionID uint64
	transactionID      uint32
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

func (udp *UDPTrackerProtocol) SendConnectRequest(possibleURLs []*url.URL) error {
	var payload ConnectRequestPayload
	buf := payload.Build()
	udp.connDialer = net.Dialer{}
	conn, err := udp.connDialer.Dial(possibleURLs[3].Scheme, possibleURLs[3].Host)
	if err != nil {
		return err
	}
	d, err := conn.Write(buf)
	checkErr(err)
	if d != len(buf) {
		return fmt.Errorf("not all data sent")
	}
	resp := make([]byte, 16)
	d, err = conn.Read(resp)
	checkErr(err)
	if d != len(resp) {
		return fmt.Errorf("not all data received")
	}
	var responsePayload ConnectResponsePayload
	responsePayload.Marshall(resp)
	udp.serverConnectionID = responsePayload.connectionID
	if payload.transactionID != responsePayload.transactionID {
		return UnEqualTransactionIDError{payload.transactionID, responsePayload.transactionID}
	}
	if payload.action != responsePayload.action {
		return UnEqualActionError{payload.action, responsePayload.action}
	}
	return nil
}
