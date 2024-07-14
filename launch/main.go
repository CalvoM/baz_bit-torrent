package main

import (
	bazbittorrent "github.com/CalvoM/baz_bit-torrent"
)

func main() {
	m := bazbittorrent.MetaInfoFile{}
	m.UnMarshalFile("samples/tears-of-steel.torrent")
	client := bazbittorrent.Client{}
	client.Init(m)
	// var enc bazbittorrent.BencodingEncoder
	// enc.Encode(m)
	// m.UnMarshallToDict()
	// fmt.Println(enc.EncodedData)
}
