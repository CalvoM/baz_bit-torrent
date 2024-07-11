package bazbittorrent

import (
	"encoding/json"
	"os"
)

type MultiFiles struct {
	Length int      `json:"length"`
	Md5sum string   `json:"md5sum"` // optional
	Path   []string `json:"path"`
}
type Info struct {
	PieceLength int          `json:"piece length"` // common
	Pieces      string       `json:"pieces"`       // common
	Private     uint64       `json:"private"`      // common optional
	Name        string       `json:"name"`         // common
	Length      int          `json:"length"`       // single-file
	Md5sum      string       `json:"md5sum"`       // single-file optional
	Files       []MultiFiles `json:"files"`
}

type MetaInfoFile struct {
	Announce     string     `json:"announce"`
	Info         Info       `json:"info"`
	AnnounceList [][]string `json:"announce-list"` // optional
	CreationDate int64      `json:"creation date"` // optional
	Comment      string     `json:"comment"`       // optional
	CreatedBy    string     `json:"created by"`    // optional
	Encoding     string     `json:"encoding"`      // optional
}

func (metainfofile *MetaInfoFile) UnMarshalFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bencoding := NewDecoder(file)
	v, err := bencoding.Decode()
	if err != nil {
		panic(err)
	}
	jstr, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jstr, &metainfofile)
	if err != nil {
		panic(err)
	}
}
