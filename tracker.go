package bazbittorrent

type Tracker struct{}

func (tracker Tracker) Get(InfoHash [20]byte, PeerID [20]byte) {
}
