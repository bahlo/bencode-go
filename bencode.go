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
	res, err := decoder.ReadSlice('e')
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

func (decoder *decoder) readString() (interface{}, error) {
	return "", nil
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
