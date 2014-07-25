package octorrent

import (
	"bytes"
	"testing"
)

func TestReadInt(t testing.T) {
	const in, out = "412e51", 412
	buf := bytes.NewBufferString(in)
	dec := decoder{
		*buf,
	}

	out, err := dec.readInt()
	if err != nil {
		t.Errorf("readInt returned error: %s", err)
	}

	if out != in {
		t.Errorf("readInt should return %i, but did return %i", in, out)
	}
}
