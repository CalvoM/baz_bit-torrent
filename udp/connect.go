package udp

import (
	"encoding/binary"
)

const connectRequestID = 0x41727101980

type ConnectResponsePayload struct {
	action        Action
	transactionID TransactionID
	connectionID  uint64
}

func (payload *ConnectResponsePayload) Marshall(buf []byte) {
	payload.action = Action(binary.BigEndian.Uint32(buf[0:]))
	payload.transactionID = TransactionID(binary.BigEndian.Uint32(buf[4:]))
	payload.connectionID = binary.BigEndian.Uint64(buf[8:])
}

type ConnectRequestPayload struct {
	protocolID    uint64
	action        Action
	transactionID TransactionID
}

func (payload *ConnectRequestPayload) Build() (buf []byte) {
	buf = make([]byte, 16)
	protocolIDBuf := make([]byte, 8)
	transactionBuf := make([]byte, 4)
	actionBuf := make([]byte, 4)
	transaction := Transaction{}
	payload.protocolID = connectRequestID
	payload.transactionID = transaction.New()
	payload.action = Connect
	binary.BigEndian.PutUint64(protocolIDBuf, payload.protocolID)
	binary.BigEndian.PutUint32(actionBuf, uint32(payload.action))
	binary.BigEndian.PutUint32(transactionBuf, uint32(payload.transactionID))
	copy(buf[0:], protocolIDBuf)
	copy(buf[8:], actionBuf)
	copy(buf[12:], transactionBuf)
	return
}
