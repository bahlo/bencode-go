package octorrent

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func getDecoder(d []byte) decoder {
	buf := bytes.NewBuffer(d)

	return decoder{
		*bufio.NewReader(buf),
	}
}

func TestReadInt(t *testing.T) {
	const i, o = "412e51", int64(412)
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
	const i, o = "142", int64(14)
	dec := getDecoder([]byte(i))

	num, err := dec.readIntUntil('2')
	if err != nil {
		t.Errorf("readIntUntil returned error: %v", err)
	}

	if num != o {
		t.Errorf("readIntUntil did return %i, but should return %i", num, o)
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
