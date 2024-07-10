package main

import (
	"encoding/hex"
	"fmt"
	"os"

	bazbittorrent "github.com/CalvoM/baz_bit-torrent"
)

func main() {
	file, err := os.Open("samples/archlinux20240701.torrent")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bencoding := bazbittorrent.NewDecoder(file)
	v, _ := bencoding.Decode()
	fmt.Println(hex.Dump([]byte(v["info"].(map[string]any)["pieces"].(string))))
}
