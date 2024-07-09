package bazbittorrent

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type BencodingDecoder struct {
	r         *bufio.Reader
	bytesRead int
	data      map[string]any
}

func NewDecoder(r io.Reader) *BencodingDecoder {
	return &BencodingDecoder{r: bufio.NewReader(r), bytesRead: 0}
}

func (d *BencodingDecoder) DecodeDict() (value map[any]any, err error) {
	ch, err := d.readByte()
	if err != nil {
		return
	}
	if ch != 'd' {
		panic("expected 'd' for dictionary")
	}
	value = make(map[any]any)
	next, err := d.peek()
	if err != nil {
		return
	}
	var tempHolder []any
	for next != 'e' {
		switch next {
		case 'i':
			i, ierr := d.DecodeInt()
			if ierr != nil {
				return
			}
			tempHolder = append(tempHolder, i)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s, strerr := d.DecodeString()
			if strerr != nil {
				return
			}
			tempHolder = append(tempHolder, s)
		case 'l':
			l, lerr := d.DecodeList()
			if lerr != nil {
				return
			}
			tempHolder = append(tempHolder, l)
		case 'd':
			d, derr := d.DecodeDict()
			if derr != nil {
				return
			}
			tempHolder = append(tempHolder, d)
		}
		if len(tempHolder) == 2 {
			value[tempHolder[0]] = tempHolder[1]
			tempHolder = nil
		}
		next, err = d.peek()
		if err != nil {
			return
		}
	}
	return
}

func (d *BencodingDecoder) DecodeString() (value string, err error) {
	strLen, err := d.readBytes(':')
	if err != nil {
		return "", err
	}
	len, err := strconv.Atoi(string(strLen[:len(strLen)-1]))
	if err != nil {
		return "", err
	}
	stringBuf := make([]byte, len)
	buf, err := d.read(stringBuf)
	if buf != len {
		panic("Buffer reading error")
	}
	if err != nil {
		return "", err
	}
	value = string(stringBuf)
	return
}

func (d *BencodingDecoder) DecodeInt() (value int, err error) {
	ch, err := d.readByte()
	if err != nil {
		return -1, err
	}
	if ch != 'i' {
		panic("Expected i for Integer")
	}
	line, err := d.readBytes('e')
	if err != nil {
		return -1, err
	}
	value, err = strconv.Atoi(string(line[:len(line)-1]))
	return
}

func (d *BencodingDecoder) DecodeList() (value []any, err error) {
	ch, err := d.readByte()
	if ch != 'l' {
		panic("Expected l for List")
	}
	next, err := d.peek()
	if err != nil {
		return
	}
	for next != 'e' {
		switch next {
		case 'i':
			i, ierr := d.DecodeInt()
			if ierr != nil {
				return
			}
			value = append(value, i)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s, strerr := d.DecodeString()
			if strerr != nil {
				return
			}
			value = append(value, s)
		case 'l':
			l, lerr := d.DecodeList()
			if lerr != nil {
				return
			}
			value = append(value, l)

		}
		next, err = d.peek()
		if err != nil {
			return
		}
	}
	return
}

func (d *BencodingDecoder) DecodeToJSON() string {
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
