package udp

import (
	"encoding/binary"
)

type (
	Peer struct {
		ip   uint32
		port uint16
	}
	AnnounceResponsePayload struct {
		action        Action
		transactionID TransactionID
		interval      uint32
		leechers      uint32
		seeders       uint32
		peers         []Peer
	}
	AnnounceRequestPayload struct {
		connectionID  uint64
		action        Action
		transactionID TransactionID
		infoHash      [20]uint8
		peerID        [20]uint8
		downloaded    uint64
		left          uint64
		uploaded      uint64
		event         Event
		ip            uint32
		key           uint32
		numWant       int32
		port          uint16
	}
)

const (
	announceBasicIPV4RespSize = 20
	peerSize                  = 6
)

func (payload *AnnounceRequestPayload) Build(connectionID uint64, infoHash [20]uint8, peerID [20]byte, left uint64) (buf []byte) {
	payload.connectionID = connectionID
	payload.action = Announce
	transaction := Transaction{}
	payload.transactionID = transaction.New()
	payload.infoHash = infoHash
	payload.peerID = peerID
	payload.downloaded = 0
	payload.left = left
	payload.uploaded = 0
	payload.event = None
	payload.ip = 0
	payload.key = uint32(transaction.New())
	payload.numWant = -1
	payload.port = 80
	buf = make([]byte, 98)
	fourByteBuf := make([]byte, 4)
	eightByteBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(eightByteBuf, payload.connectionID)
	copy(buf[0:], eightByteBuf)
	binary.BigEndian.PutUint32(fourByteBuf, uint32(payload.action))
	copy(buf[8:], fourByteBuf)
	binary.BigEndian.PutUint32(fourByteBuf, uint32(payload.transactionID))
	copy(buf[12:], fourByteBuf)
	copy(buf[16:], infoHash[:])
	copy(buf[36:], peerID[:])
	binary.BigEndian.PutUint64(eightByteBuf, payload.downloaded)
	copy(buf[56:], eightByteBuf)
	binary.BigEndian.PutUint64(eightByteBuf, payload.left)
	copy(buf[64:], eightByteBuf)
	binary.BigEndian.PutUint64(eightByteBuf, payload.uploaded)
	copy(buf[72:], eightByteBuf)
	binary.BigEndian.PutUint32(fourByteBuf, uint32(payload.event))
	copy(buf[80:], fourByteBuf)
	binary.BigEndian.PutUint32(fourByteBuf, payload.ip)
	copy(buf[84:], fourByteBuf)
	binary.BigEndian.PutUint32(fourByteBuf, payload.key)
	copy(buf[88:], fourByteBuf)
	binary.BigEndian.PutUint32(fourByteBuf, uint32(payload.numWant))
	copy(buf[92:], fourByteBuf)
	twoByteBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(twoByteBuf, payload.port)
	copy(buf[96:], fourByteBuf)
	return
}

func (payload *AnnounceResponsePayload) Marshall(buf []byte) {
	payload.action = Action(binary.BigEndian.Uint32(buf[0:]))
	payload.transactionID = TransactionID(binary.BigEndian.Uint32(buf[4:]))
	payload.interval = binary.BigEndian.Uint32(buf[8:])
	payload.leechers = binary.BigEndian.Uint32(buf[12:])
	payload.seeders = binary.BigEndian.Uint32(buf[16:])
	payload.peers = nil
}

func MarshallPeers(buf []byte, peersCount int) (peers []Peer) {
	for i := 0; i < peersCount; i++ {
		var peer Peer
		peer.ip = binary.BigEndian.Uint32(buf[i*6:])
		peer.port = binary.BigEndian.Uint16(buf[(i*6)+4:])
		peers = append(peers, peer)
	}
	return
}
