package main

import (
	"fmt"
	"os"

	bazbittorrent "github.com/CalvoM/baz_bit-torrent"
)

func main() {
	file, err := os.Open("samples/numbers.torrent")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bencoding := bazbittorrent.NewDecoder(file)
	fmt.Println(bencoding.Decode())
}
