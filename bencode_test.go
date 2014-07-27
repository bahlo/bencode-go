package octorrent

import (
	"bufio"
	"bytes"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func getDecoder(d []byte) decoder {
	buf := bytes.NewBuffer(d)

	return decoder{
		*bufio.NewReader(buf),
	}
}

func TestReadInt(t *testing.T) {
	i, o := "412e51", int64(412)
	dec := getDecoder([]byte(i))

	num, err := dec.readInt()
	if err != nil {
		t.Errorf("readInt returned error: %v", err)
	}

	if num != o {
		t.Errorf("readInt did return %i, but should return %i", num, o)
	}
}

func TestReadIntUntil(t *testing.T) {
	i, o := "142", int64(14)
	dec := getDecoder([]byte(i))

	num, err := dec.readIntUntil('2')
	if err != nil {
		t.Errorf("readIntUntil returned error: %v", err)
	}

	if num != o {
		t.Errorf("readIntUntil did return %i, but should return %i", num, o)
	}

	// Check broken reader
	dec = getDecoder([]byte(""))
	num, err = dec.readIntUntil('t')
	if err == nil {
		t.Errorf("readIntUntil with empty reader returned no error, but should")
	}

	dec = getDecoder([]byte("test"))
	num, err = dec.readIntUntil('e')
	if err == nil {
		t.Errorf("readIntUntil with no digits returned no error, but should")
	}

	// Check unsigned int
	var ui uint64
	ui = math.MaxInt64 + 1
	i = "i" + strconv.FormatUint(ui, 10) + "e"
	dec = getDecoder([]byte(i))
	num, err = dec.readIntUntil('e')
	if err == nil {
		t.Errorf("readIntUntil with unsigned int returned no error, but should")
	}
}

func TestReadString(t *testing.T) {
	const i, o = "4:testfoobar", "test"
	dec := getDecoder([]byte(i))

	str, err := dec.readString()
	if err != nil {
		t.Errorf("readString returned error: %v", err)
	}

	if str != o {
		t.Errorf("readString returned %s, but should return %s", str, o)
	}
}

func TestReadList(t *testing.T) {
	const i = "4:testi4ee"
	o := []interface{}{"test", int64(4)}
	dec := getDecoder([]byte(i))

	list, err := dec.readList()
	if err != nil {
		t.Errorf("readList returned error: %v", err)
	}

	if !reflect.DeepEqual(list, o) {
		t.Errorf("readList returned %v, but should return %v", list, o)
	}
}

func TestIsEnd(t *testing.T) {
	i, o := "4:test", false
	dec := getDecoder([]byte(i))
	end, err := dec.isEnd()
	if err != nil {
		t.Errorf("isEnd returned error: %v", err)
	}
	if end != o {
		t.Errorf("isEnd returned %v, but should return %v", end, o)
	}

	i, o = "e4:test", true
	dec = getDecoder([]byte(i))
	end, err = dec.isEnd()
	if err != nil {
		t.Errorf("isEnd returned error: %v", err)
	}
	if end != o {
		t.Errorf("isEnd returned %v, but should return %v", end, o)
	}
}

func TestReadDictionary(t *testing.T) {
	const i = "4:testi1337e3:foo3:bare"
	o := map[string]interface{}{
		"test": int64(1337),
		"foo":  "bar",
	}
	dec := getDecoder([]byte(i))

	dict, err := dec.readDictionary()
	if err != nil {
		t.Errorf("readDictionary returned error: %v", err)
	}

	if !reflect.DeepEqual(dict, o) {
		t.Errorf("readDictionary returned %v, but should %v", dict, o)
	}
}

func TestDecode(t *testing.T) {
	const i = "d4:testli4ei3ee3:foo3:bare"
	o := map[string]interface{}{
		"test": []interface{}{
			int64(4),
			int64(3),
		},
		"foo": "bar",
	}
	buf := bytes.NewBufferString(i)

	res, err := Decode(buf)
	if err != nil {
		t.Errorf("Decode returned error: %v", err)
	}

	if !reflect.DeepEqual(res, o) {
		t.Errorf("Decode returned %v, but should return %v", res, o)
	}
}
