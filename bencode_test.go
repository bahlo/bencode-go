package octorrent

import (
	"bytes"
	"testing"
)

func getDecoder(d []byte) decoder {
	buf := bytes.NewBuffer(d)

	return decoder{
		*buf,
	}
}

func TestReadInt(t *testing.T) {
	const i, o = "412e51", 412
	dec := getDecoder([]byte(i))

	num, err := dec.readInt()
	if err != nil {
		t.Errorf("readInt returned error: %v", err)
	}

	if num != o {
		t.Errorf("readInt did return %i, but should return %i", num, out)
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
	const i, o = "l4:testi4ee", []interface{}{"test", 4}
	dec := getDecoder([]byte(i))

	list, err := dec.readList()
	if err != nil {
		t.Errorf("readList returned error: %v", err)
	}

	if list != o {
		t.Errorf("readList returned %v, but should return %v", list, o)
	}
}

func TestReadDictionary(t *testing.T) {
	const i, o = "d4:testi1337e3:foo3:bare", map[string]interface{}{
		"test": 1337,
		"foo":  "bar",
	}
	dec := getDecoder([]byte(i))

	dict, err := dec.readDictionary()
	if err != nil {
		t.Errorf("readDictionary returned error: %v", err)
	}

	if dict != o {
		t.Errorf("readDictionary returned %v, but should %v", dict, o)
	}
}

func TestDecode(t *testing.T) {
	const i, o = "d4:testli4ei3ee3:foo3:bar", map[string]interface{}{
		"test": []int{
			4,
			3,
		},
		"foo": "bar",
	}
	buf := bytes.NewBufferString(i)

	res, err := Decode(buf)
	if err != nil {
		t.Errorf("Decode returned error: %v", err)
	}

	if res != o {
		t.Errorf("Decode returned %v, but should return %v", res, o)
	}
}
