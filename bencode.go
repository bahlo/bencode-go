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
	buf := make([]byte, strLen)

	// Get length of string
	l, err := decoder.readIntUntil(':')
	if err != nil {
		return "", err
	}

	// Read into buffer
	strLen := l.(int64)
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
	return []interface{}{1, 2}, nil
}

func (decoder *decoder) readDictionary() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func Decode(reader io.Reader) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
