package main

import (
	"fmt"

	bazbittorrent "github.com/CalvoM/baz_bit-torrent"
)

func main() {
	m := bazbittorrent.MetaInfoFile{}
	m.UnMarshalFile("samples/tears-of-steel.torrent")
	fmt.Println(m)
}
