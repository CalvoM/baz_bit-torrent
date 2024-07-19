package main

import (
	bazbittorrent "github.com/CalvoM/baz_bit-torrent"
)

func main() {
	m := bazbittorrent.MetaInfoFile{}
	// m.UnMarshalFile("samples/PoI_SSN1.torrent")
	m.UnMarshalFile("samples/sintel.torrent")
	client := bazbittorrent.Client{}
	client.Init(m)
	// var enc bazbittorrent.BencodingEncoder
	// enc.Encode(m)
	// m.UnMarshallToDict()
	// fmt.Println(enc.EncodedData)
}
