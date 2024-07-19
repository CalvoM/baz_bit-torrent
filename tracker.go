package bazbittorrent

import (
	"fmt"
	"net/url"
	"strconv"
)

func (metainfofile *MetaInfoFile) buildTrackerURL(peerID []byte) (string, error) {
	baseURL, err := url.Parse(metainfofile.Announce)
	if err != nil {
		panic(err)
	}
	params := url.Values{
		"info_hash":  []string{fmt.Sprintf("%x", metainfofile.InfoHash())},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{baseURL.Port()},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(metainfofile.Info.Length)},
	}
	baseURL.RawQuery = params.Encode()
	return baseURL.String(), nil
}

func (metainfofile *MetaInfoFile) HostDetails() []*url.URL {
	var validURLs []*url.URL
	baseURL, err := url.Parse(metainfofile.Announce)
	if err != nil {
		panic(err)
	}
	validURLs = append(validURLs, baseURL)
	if len(metainfofile.AnnounceList) > 0 {
		for _, baseURL := range metainfofile.AnnounceList {
			u, err := url.Parse(baseURL[0])
			if err != nil {
				panic(err)
			}
			validURLs = append(validURLs, u)
		}
	}
	return validURLs
}
