package bazbittorrent

import (
	"net"
	"net/url"
)

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
}
