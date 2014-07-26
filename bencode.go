package octorrent

import (
	"bufio"
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

func (decoder *decoder) readList() (interface{}, error) {
	var list []interface{}

	for {
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
		case 'e':
			return list, nil
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

		list = append(list, item)
	}

	return []interface{}{1, 2}, nil
}

func (decoder *decoder) readDictionary() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func Decode(reader io.Reader) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
