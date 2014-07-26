package octorrent

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type decoder struct {
	bufio.Reader
}

func (decoder *decoder) readInt() (interface{}, error) {
	return decoder.readIntUntil('e')
}

func (decoder *decoder) readIntUntil(b byte) (interface{}, error) {
	res, err := decoder.ReadSlice(b)
	if err != nil {
		return -1, err
	}

	// Read until before the 'e'nd
	str := string(res[:len(res)-1])

	if value, err := strconv.ParseInt(str, 10, 64); err == nil {
		return value, nil
	} else if value, err := strconv.ParseUint(str, 10, 64); err == nil {
		return value, nil
	}

	return -1, err
}

func (decoder *decoder) readString() (string, error) {
	// Get length of string
	l, err := decoder.readIntUntil(':')
	if err != nil {
		return "", err
	}

	// Read into buffer
	strLen := l.(int64)
	buf := make([]byte, strLen)
	pos := int64(0)
	for pos < strLen {
		n, err := decoder.Read(buf[pos:])
		if err != nil {
			return "", err
		}

		pos += int64(n)
	}

	return string(buf), nil
}

func (decoder *decoder) readValue() (interface{}, error) {
	b, err := decoder.ReadByte()
	if err != nil {
		return nil, err
	}

	var item interface{}
	switch b {
	case 'i':
		item, err = decoder.readInt()
	case 'l':
		item, err = decoder.readList()
	case 'd':
		item, err = decoder.readDictionary()
	default:
		// Unread last byte, because it belongs to the lenght
		// of the string
		if err := decoder.UnreadByte(); err != nil {
			return nil, err
		}
		item, err = decoder.readString()
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

// Returns if the next byte is an 'e'
func (decoder *decoder) isEnd() (bool, error) {
	b, err := decoder.ReadByte()
	if err != nil {
		return false, err
	} else if b == 'e' {
		return true, nil
	} else if err := decoder.UnreadByte(); err != nil {
		// Unread last byte since it's not an 'e'
		return false, err
	}

	return false, nil
}

func (decoder *decoder) readList() (interface{}, error) {
	var list []interface{}

	for {
		end, err := decoder.isEnd()
		if err != nil {
			return nil, err
		} else if end {
			return list, nil
		}

		item, err := decoder.readValue()
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	return list, nil
}

func (decoder *decoder) readDictionary() (map[string]interface{}, error) {
	var dict = map[string]interface{}{}

	for {
		// Check if end
		end, err := decoder.isEnd()
		if err != nil {
			return nil, err
		} else if end {
			return dict, nil
		}

		// Get key
		key, err := decoder.readString()
		if err != nil {
			return nil, err
		}

		// Get value
		val, err := decoder.readValue()
		if err != nil {
			return nil, err
		}

		dict[key] = val
	}

	return dict, nil
}

func Decode(reader io.Reader) (map[string]interface{}, error) {
	decoder := decoder{*bufio.NewReader(reader)}

	if b, err := decoder.ReadByte(); err != nil {
		return nil, err
	} else if b != 'd' {
		return nil, errors.New("bencode data must be a dictionary at top level")
	}

	return decoder.readDictionary()
}
