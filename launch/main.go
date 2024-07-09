package main

import (
	"fmt"
	"os"

	bazbittorrent "github.com/CalvoM/baz_bit-torrent"
)

func main() {
	file, err := os.Open("samples/alice.torrent")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bencoding := bazbittorrent.NewDecoder(file)
	v, _ := bencoding.DecodeToJSON()
	fmt.Println(string(v.(string)))
}
