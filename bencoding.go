package bazbittorrent

import (
	"bufio"
	"errors"
	"io"
)

type BencodingDecoder struct {
	r         *bufio.Reader
	bytesRead int
}

func NewDecoder(r io.Reader) *BencodingDecoder {
	return &BencodingDecoder{r: bufio.NewReader(r), bytesRead: 0}
}

func (d *BencodingDecoder) decodeInt() error {
	ch, err := d.readByte()
	if err != nil {
		return err
	}
	if ch != 'i' {
		panic("Expected i for Integer")
	}
	line, err := d.readBytes('e')
	if err != nil {
		return err
	}
	digits := string(line[:len(line)-1])
	return nil
}

func (d *BencodingDecoder) Decode() string {
	next, err := d.peek()
	if err != nil {
		return ""
	}
	switch next {
	case 'i':
		return "integer"
	case 'l':
		return "list"
	case 'd':
		return "dictionary"
	default:
		err = errors.New("invalid Input")
		return string(err.Error())
	}
	return ""
}

func (d *BencodingDecoder) peek() (b byte, err error) {
	ch, err := d.r.Peek(1)
	if err != nil {
		return
	}
	b = ch[0]
	return
}

func (d *BencodingDecoder) readByte() (b byte, err error) {
	b, err = d.r.ReadByte()
	d.bytesRead++
	return
}

func (d *BencodingDecoder) readBytes(delim byte) (line []byte, err error) {
	line, err = d.r.ReadBytes(delim)
	d.bytesRead += len(line)
	return
}

func (d *BencodingDecoder) read(p []byte) (n int, err error) {
	n, err = d.r.Read(p)
	d.bytesRead += n
	return
}
