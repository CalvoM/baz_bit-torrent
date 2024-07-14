package bazbittorrent

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
)

type BencodingDecoder struct {
	r         *bufio.Reader
	bytesRead int
	data      map[string]any
}

type BencodingEncoder struct {
	EncodedData string
}

func NewDecoder(r io.Reader) *BencodingDecoder {
	return &BencodingDecoder{r: bufio.NewReader(r), bytesRead: 0}
}

func (d *BencodingDecoder) DecodeDict() (value map[string]any, err error) {
	ch, err := d.readByte()
	if err != nil {
		return
	}
	if ch != 'd' {
		panic("expected 'd' for dictionary")
	}
	value = make(map[string]any)
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
			dict, derr := d.DecodeDict()
			if derr != nil {
				return
			}
			tempHolder = append(tempHolder, dict)
		}
		if len(tempHolder) == 2 {
			value[tempHolder[0].(string)] = tempHolder[1]
			tempHolder = nil
		}
		next, err = d.peek()
		if err != nil {
			return
		}
	}
	d.readBytes('e')
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

		case 'd':
			dict, derr := d.DecodeDict()
			if derr != nil {
				return
			}
			value = append(value, dict)
		}
		next, err = d.peek()
		if err != nil {
			return
		}
	}
	d.readBytes('e')
	return
}

func (d *BencodingDecoder) Decode() (map[string]any, error) {
	next, err := d.peek()
	if err != nil {
		return nil, err
	}
	var finalVal map[string]any
	switch next {
	case 'd':
		finalVal, _ = d.DecodeDict()
	default:
		err = errors.New("invalid Start of bencoding")
		return nil, err
	}
	return finalVal, err
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
	n, err = io.ReadFull(d.r, p)
	d.bytesRead += n
	return
}

func (e *BencodingEncoder) Encode(decoded Mapable) {
	// TODO: Expand to other data types, now we handle structs.
	supportedTypes := []reflect.Kind{reflect.Struct}
	k := reflect.TypeOf(decoded).Kind()
	if !slices.Contains(supportedTypes, k) {
		errMsg := fmt.Sprintf("Data type not supported: %v", k)
		panic(errMsg)
	}
	decodedMap := decoded.UnMarshallToDict()
	e.EncodedData += e.bencodeMap(decodedMap)
}

func (e BencodingEncoder) bencodeMap(decodeMap map[string]any) string {
	bencodeMap := "d"
	for key, value := range decodeMap {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			bencodeMap += e.bencodeStr(key) + e.bencodeStr(value.(string))
		case reflect.Float64:
			bencodeMap += e.bencodeStr(key) + e.bencodeInt(value.(float64))
		case reflect.Slice:
			bencodeMap += e.bencodeStr(key) + e.bencodeSlice(value.([]any))
		case reflect.Map:
			bencodeMap += e.bencodeStr(key) + e.bencodeMap(value.(map[string]any))
		default:
			fmt.Printf("Found this new %v: %v of type: %v\r\n", key, value, reflect.TypeOf(value).Kind())
		}
	}
	bencodeMap += "e"
	return bencodeMap
}

func (e BencodingEncoder) bencodeStr(decodeStr string) string {
	return fmt.Sprintf("%d:%s", len(decodeStr), decodeStr)
}

func (e BencodingEncoder) bencodeInt(decodeInt float64) string {
	return fmt.Sprintf("i%fe", decodeInt)
}

func (e BencodingEncoder) bencodeSlice(decodeSlice []any) string {
	bencodedList := "l"
	for _, item := range decodeSlice {
		switch reflect.TypeOf(item).Kind() {
		case reflect.String:
			bencodedList += e.bencodeStr(item.(string))
		case reflect.Slice:
			bencodedList += e.bencodeSlice(item.([]any))
		case reflect.Int:
			bencodedList += e.bencodeInt(item.(float64))
		case reflect.Map:
			bencodedList += e.bencodeMap(item.(map[string]any))
		default:
			fmt.Printf("Found this new list item type: %v\r\n", reflect.TypeOf(item).Kind())
		}
	}
	bencodedList += "e"
	return bencodedList
}
