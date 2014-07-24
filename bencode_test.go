package octorrent

import (
	"reflect"
	"testing"
)

func TestDecodeInt(t *testing.T) {
	i := []byte("i20e")
	const o = 20

	e, err := DecodeInt(i)

	if err != nil {
		t.Errorf("DecodeInt(%d) got error: %v", i, err)
	}

	if o != e {
		t.Errorf("DecodeInt(%d) = %d, want %d", i, e, o)
	}
}

func TestDecodeString(t *testing.T) {
	i := []byte("7:testing")
	const o = "testing"

	e, err := DecodeString(i)

	if err != nil {
		t.Errorf("DecodeString(%s) got error: %v", i, err)
	}

	if o != e {
		t.Errorf("DecodeString(%s) = %s, want %s", i, e, o)
	}
}

func TestIndexStringEnd(t *testing.T) {
	i := []byte("7:testingi4e")
	const o = 8

	e, err := indexStringEnd(i)

	if err != nil {
		t.Errorf("indexStringEnd(%s) got error: %v", i, err)
	}

	if o != e {
		t.Errorf("indexStringEnd(%s) = %d, want %d", i, e, o)
	}
}

func TestIndexIntEnd(t *testing.T) {
	i := []byte("i1337e3:foo")
	const o = 5

	e, err := indexIntEnd(i)

	if err != nil {
		t.Errorf("indexIntEnd(%s) got error: %v", i, err)
	}

	if o != e {
		t.Errorf("indexIntEnd(%s) = %d, want %d", i, e, o)
	}
}

func TestIndexListEnd(t *testing.T) {
	i := []byte("l3:foo3:barei3e")
	const o = 11

	e, err := indexListEnd(i)

	if err != nil {
		t.Errorf("indexListEnd(%s) got error: %v", i, err)
	}

	if o != e {
		t.Errorf("indexListEnd(%s) = %d, want %d", i, e, o)
	}
}

func TestIndexDictEnd(t *testing.T) {
	i := []byte("d3:foo3:bar6:numberi3ee4:this3:not")
	const o = 22

	e, err := indexDictEnd(i)

	if err != nil {
		t.Errorf("indexDictEnd(%s) got error: %v", i, err)
	}

	if o != e {
		t.Errorf("indexDictEnd(%s) = %d, want %d", i, e, o)
	}
}

func TestIndexEnd(t *testing.T) {
	i := []byte("i3e3:fooi360e3:bare3:not")
	const o = 18

	e, err := indexEnd(i)

	if err != nil {
		t.Errorf("indexEnd(%s) got error: %v", i, err)
	}

	if o != e {
		t.Errorf("indexEnd(%s) = %d, want %d", i, e, o)
	}
}

func TestDecode(t *testing.T) {
	i := []byte("7:testing7:bencode3:foo3:bar")
	o := []interface{}{"testing", "bencode", "foo", "bar"}

	e, err := decode(i)

	if err != nil {
		t.Errorf("decode(%s) got error: %v", i, err)
	}

	if !reflect.DeepEqual(o, e) {
		t.Errorf("decode(%s) = %v, want %v", i, e, o)
	}
}

func TestMakeMap(t *testing.T) {
	i := []interface{}{"testing", "bencode", "foo", "bar"}
	o := map[string]interface{}{
		"testing": "bencode",
		"foo":     "bar",
	}

	m, err := makeMap(i)

	if err != nil {
		t.Errorf("makeMap(%v) got error: %v", i, err)
	}

	if !reflect.DeepEqual(o, m) {
		t.Errorf("makeMap(%v) = %v, want %v", i, m, o)
	}
}

func TestDecodeList(t *testing.T) {
	i := []byte("l3:foo3:bari4ee")
	o := []interface{}{"foo", "bar", 4}

	e, err := DecodeList(i)

	if err != nil {
		t.Errorf("DecodeList(%s) got error: %v", i, err)
	}

	if !reflect.DeepEqual(o, e) {
		t.Errorf("DecodeList(%s) = %s, want %s", i, e, o)
	}
}

func TestDecodeDictionary(t *testing.T) {
	i := []byte("d7:testing7:bencode3:foo3:bar3:numi16ee")
	o := map[string]interface{}{
		"testing": "bencode",
		"foo":     "bar",
		"num":     16,
	}

	e, err := DecodeDictionary(i)

	if err != nil {
		t.Errorf("EncodeDictionary(%s) got error: %v", i, err)
	}

	if !reflect.DeepEqual(o, e) {
		t.Errorf("EncodeDictionary(%s) = %v, want %v", i, e, o)
	}
}

func TestNestedDictionaries(t *testing.T) {
	i := []byte("d3:foo3:bar4:dictd3:numi5eee")
	o := map[string]interface{}{
		"foo": "bar",
		"dict": map[string]interface{}{
			"num": 5,
		},
	}

	e, err := DecodeDictionary(i)

	if err != nil {
		t.Errorf("DecodeDictionary(%s) got error %v", i, err)
	}

	if !reflect.DeepEqual(o, e) {
		t.Errorf("DecodeDictionary(%s) = %v, want %v", i, e, o)
	}
}

func TestNestedLists(t *testing.T) {
	i := []byte("l3:foo3:barli3ei5eee")
	o := []interface{}{
		"foo",
		"bar",
		[]interface{}{
			3,
			5,
		},
	}

	e, err := DecodeList(i)

	if err != nil {
		t.Errorf("DecodeList(%s) got error: %v", i, err)
	}

	if !reflect.DeepEqual(o, e) {
		t.Errorf("DecodeList(%s) = %v, want %v", i, e, o)
	}
}
