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
		t.Errorf("readInt returned error: %s", err)
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
		t.Errorf("readString returned error: %s", err)
	}

	if str != o {
		t.Errorf("readString returned %s, but should return %s", str, o)
	}
}
