package bazbittorrent

import (
	"fmt"
	"net"

	"github.com/CalvoM/baz_bit-torrent/udp"
)

const (
	BitTorrentIdentifier = "BitTorrent protocol"
	BazPeerID            = "BazeengaBitTorrent24"
)

type HandShakePayLoad struct {
	identifier string
	reserved   []byte
	infoHash   [20]byte
	peerID     [20]byte
}

type Client struct {
	conn         net.Conn
	metaInfoFile MetaInfoFile
}

func (handShake *HandShakePayLoad) Build(PeerID []byte, InfoHash [20]byte) []byte {
	handShake.identifier = BitTorrentIdentifier
	handShake.reserved = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	handShake.infoHash = [20]byte(InfoHash)
	handShake.peerID = [20]byte(PeerID)
	outputBuf := make([]byte, len(handShake.identifier)+49)
	outputBuf[0] = byte(19)
	curr := 1
	curr += copy(outputBuf[curr:], handShake.identifier)
	curr += copy(outputBuf[curr:], handShake.reserved[:])
	curr += copy(outputBuf[curr:], handShake.infoHash[:])
	curr += copy(outputBuf[curr:], handShake.peerID[:])
	return outputBuf
}

func (c *Client) sendHandShakeToClient(peer udp.Peer) {
	conn, err := net.Dial("tcp", peer.URL())
	if err != nil {
		fmt.Println(err)
	}
	handshake := HandShakePayLoad{}
	peerID := "BazeengaBitTorrent24"
	buf := handshake.Build([]byte(peerID), c.metaInfoFile.InfoHash())
	conn.Write(buf)
}

func (c *Client) Init(metaInfoFile MetaInfoFile) {
	c.metaInfoFile = metaInfoFile
	c.SendConnectRequest()
}

func (c *Client) SendConnectRequest() {
	udpTracker := udp.UDPTrackerProtocol{}
	err := udpTracker.ConnectToTracker(c.metaInfoFile.HostDetails())
	if err != nil {
		panic(err)
	}
	peerID := "BazeengaBitTorrent24"
	peers, err := udpTracker.AnnounceToTracker(c.metaInfoFile.InfoHash(), []byte(peerID), uint64(c.metaInfoFile.Info.Length))
	fmt.Println(err)
	for _, peer := range peers {
		c.sendHandShakeToClient(peer)
		fmt.Println(peer.IP(), peer.Port())
	}
}
