package octorrent

import (
	"bytes"
	"fmt"
	"strconv"
)

// DecodeInt decodes a bencoded integer
func DecodeInt(data []byte) (int, error) {
	// Check if bencoded int
	if !bytes.HasPrefix(data, []byte("i")) ||
		!bytes.HasSuffix(data, []byte("e")) {
		return 0, fmt.Errorf("not a bencoded int")
	}

	str := string(data[1 : len(data)-1])
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// DecodeString decodes a bencoded string
func DecodeString(data []byte) (string, error) {
	// Get index of the colon
	i := bytes.IndexByte(data, byte(':'))
	l := len(data)

	// Check if bencoded string
	prefix := []byte(strconv.Itoa(l - 2))
	if !bytes.HasPrefix(data, prefix) || i < 1 {
		return "", fmt.Errorf("not a bencoded string")
	}

	return string(data[i+1 : l]), nil
}

// DecodeList decodes a bencoded list
func DecodeList(data []byte) ([]interface{}, error) {
	// Check if list
	if !bytes.HasPrefix(data, []byte("l")) ||
		!bytes.HasSuffix(data, []byte("e")) {
		return []interface{}{}, fmt.Errorf("not a bencoded list")
	}

	return decode(data[1 : len(data)-1])
}

// DecodeDictionary decodes a bencoded dictionary
func DecodeDictionary(data []byte) (map[string]interface{}, error) {
	// Check if dictionary
	if !bytes.HasPrefix(data, []byte("d")) ||
		!bytes.HasSuffix(data, []byte("e")) {
		return map[string]interface{}{}, fmt.Errorf("not a bencoded dictionary")
	}

	stringSlice, err := decode(data[1 : len(data)-1])
	if err != nil {
		return map[string]interface{}{}, err
	}

	return makeMap(stringSlice)
}

// decode parses joined bencoded values and returns them in a slice
func decode(data []byte) ([]interface{}, error) {
	ret := []interface{}{}

	for len(data) > 0 {
		first := string(data[0])

		switch first {
		default:
			// We most likely have a string
			// Try to find it's end
			i, err := indexStringEnd(data)
			if err != nil {
				return ret, err
			}

			// Increment i to get one char after the last
			i++

			// Try to decode it
			d, err := DecodeString(data[:i])
			if err != nil {
				return ret, err
			}

			ret = append(ret, d)
			data = data[i:]
		case "i":
			// We have an integer!
			// Try to find it's end
			i, err := indexIntEnd(data)
			if err != nil {
				return ret, err
			}

			// Increment i to get one char after the last
			i++

			// Try to decode it
			d, err := DecodeInt(data[:i])
			if err != nil {
				return ret, err
			}

			ret = append(ret, d)
			data = data[i:]
		case "l":
			// We have a list
			// Try to find it's end
			i, err := indexListEnd(data)
			if err != nil {
				return ret, err
			}

			// Increment i to get one char after the last
			i++

			// Try to decode it
			d, err := DecodeList(data[:i])
			if err != nil {
				return ret, nil
			}

			ret = append(ret, d)
			data = data[i:]
		case "d":
			// We have a dictionary
			// Try to find it's end
			i, err := indexDictEnd(data)
			if err != nil {
				return ret, err
			}

			// Increment i to get one char after the last
			i++

			d, err := DecodeDictionary(data[:i])
			if err != nil {
				return ret, err
			}

			ret = append(ret, d)
			data = data[i:]
		}
	}

	return ret, nil
}

// indexStringEnd returns the index of the end of the bencoded string
func indexStringEnd(data []byte) (int, error) {
	c := bytes.IndexByte(data, byte(':'))

	str := string(data[:c])
	e, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}

	return int(e) + c, nil
}

// indexIntEnd returns the index of the end of the bencoded int
func indexIntEnd(data []byte) (int, error) {
	e := bytes.IndexByte(data, byte('e'))

	if e < 0 {
		return -1, fmt.Errorf("could not find end of int")
	}

	return e, nil
}

// indexlistEnd returns the index of the end of the bencoded list
func indexListEnd(data []byte) (int, error) {
	if !bytes.HasPrefix(data, []byte("l")) {
		return -1, fmt.Errorf("not a bencoded list")
	}

	if len(data) < 2 {
		return -1, fmt.Errorf("list to short")
	}

	data = data[1:] // Remove "l"
	e, err := indexEnd(data)

	if err != nil {
		return -1, err
	}

	// + 1 because we strip the leading "l"
	return e + 1, nil
}

// indexDictEnd returns the index of the end of a bencoded dictionary
func indexDictEnd(data []byte) (int, error) {
	if len(data) < 2 {
		return -1, fmt.Errorf("dictionary too short")
	}

	if !bytes.HasPrefix(data, []byte("d")) {
		return -1, fmt.Errorf("not a bencoded dictionary")
	}

	data = data[1:] // Remove "d"
	e, err := indexEnd(data)

	if err != nil {
		return -1, err
	}

	// + 1 because we strip the leading "d"
	return e + 1, nil
}

// indexEnd gets the end of variables chained together
func indexEnd(data []byte) (int, error) {
	end := 0

	for len(data) > 0 {
		// TODO: Find a better way to initialize an error
		e, err := 0, fmt.Errorf("")
		first := string(data[0])

		switch first {
		default:
			// Most likely a string or an error
			e, err = indexStringEnd(data)
		case "i":
			// We have an integer
			e, err = indexIntEnd(data)
		case "l":
			// We have a list
			e, err = indexListEnd(data)
		case "d":
			// We have a dictionary
			e, err = indexDictEnd(data)
		case "e":
			// We found the end!
			return end, nil
		}

		// Check if we got an error
		if err != nil {
			return -1, err
		}

		// Increment e, because we need the next start index
		e++

		// No error returned, set new index
		end += e

		// Set string at new index for next loop
		data = data[e:]
	}

	return -1, fmt.Errorf("could not find end")
}

// makeMap takes a string array and converts them to a map[key]value
func makeMap(stringSlice []interface{}) (map[string]interface{}, error) {
	stringMap := map[string]interface{}{}
	l := len(stringSlice)

	// Check if we even got items
	if l == 0 {
		return stringMap, nil
	}

	// Check if we have an even length
	if l%2 != 0 {
		return stringMap, fmt.Errorf("key %s has no value", stringSlice[l-1])
	}

	for {
		// Check if there are enough items
		if len(stringSlice) < 2 {
			break
		}

		if key, ok := stringSlice[0].(string); ok {
			stringMap[key] = stringSlice[1:2][0]
			stringSlice = stringSlice[2:]
		}
	}

	return stringMap, nil
}
