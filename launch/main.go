package main

import (
	"fmt"

	bazbittorrent "github.com/CalvoM/baz_bit-torrent"
)

func main() {
	m := &bazbittorrent.MetaInfoFile{}
	m.UnMarshalFile("samples/tears-of-steel.torrent")
	// client := bazbittorrent.Client{}
	// client.Init(m)
	var enc bazbittorrent.BencodingEncoder
	dicatable := bazbittorrent.Mapable{m}
	enc.Encode(dicatable)
	m.UnMarshallToDict()
	fmt.Println(enc.EncodedData)
}
