package udp

import "encoding/binary"

type ScrapeRequestPayload struct {
	connectionID  uint64
	action        Action
	transactionID TransactionID
	infoHashes    [][20]uint8
}

func (payload *ScrapeRequestPayload) Build(connectionID uint64, infoHash [][20]uint8) (buf []byte) {
	payload.connectionID = connectionID
	payload.action = Scrape
	transaction := Transaction{}
	payload.transactionID = transaction.New()
	payload.infoHashes = append(payload.infoHashes, infoHash...)
	buf = make([]byte, 16+20*len(payload.infoHashes))
	fourByteBuf := make([]byte, 4)
	eightByteBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(eightByteBuf, payload.connectionID)
	copy(buf[0:], eightByteBuf)
	binary.BigEndian.PutUint32(fourByteBuf, uint32(payload.action))
	copy(buf[8:], fourByteBuf)
	binary.BigEndian.PutUint32(fourByteBuf, uint32(payload.transactionID))
	copy(buf[12:], fourByteBuf)
	for i, hash := range payload.infoHashes {
		copy(buf[16+(i*20):], hash[:])
	}
	return
}
