package bazbittorrent

import (
	"net"
	"net/url"
)

const (
	BitTorrentIdentifier = "BitTorrent protocol"
	BazPeerID            = "BazeengaBitTorrent12"
)

type HandShakePayLoad struct {
	Identifier string
	Reserved   []byte
	InfoHash   [20]byte
	PeerID     []byte
}

type Client struct {
	conn         net.Conn
	metaInfoFile MetaInfoFile
}

func (c *Client) Init(metaInfoFile MetaInfoFile) {
	var err error
	c.metaInfoFile = metaInfoFile
	serverURL, _ := url.Parse(c.metaInfoFile.Announce)
	c.conn, err = net.Dial(serverURL.Scheme, serverURL.Host)
	if err != nil {
		panic(err)
	}
	handshake := HandShakePayLoad{Identifier: BitTorrentIdentifier, Reserved: []byte{0, 0, 0, 0, 0, 0, 0, 0}, PeerID: []byte(BazPeerID), InfoHash: c.metaInfoFile.InfoHash()}
	outputBuf := make([]byte, len(handshake.Identifier)+49)
	outputBuf[0] = byte(19)
	curr := 1
	curr += copy(outputBuf[curr:], handshake.Identifier)
	curr += copy(outputBuf[curr:], handshake.Reserved[:])
	curr += copy(outputBuf[curr:], handshake.InfoHash[:])
	curr += copy(outputBuf[curr:], handshake.PeerID[:])
	c.conn.Write(outputBuf)
}
