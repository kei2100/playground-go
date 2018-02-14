package binary

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

// 現在主流のx86, x86_64はリトルエンディアン
// ARMはどちらでも切り替えられるがAndroid, iOSはリトルエンディアン
//
// しかし、ネットワーク上で転送されるデータの多くは、大きい桁からメモリに格納されるビックエンディアンである
// 転送データのrawバイナリを解析する際はエンディアン変換が必要になる

func TestConvertEndian(t *testing.T) {
	bb := make([]byte, 2)
	binary.BigEndian.PutUint16(bb, 255)
	if g, w := bb, []byte{0, 255}; !reflect.DeepEqual(g, w) {
		t.Errorf(" got %v, want %v", g, w)
	}

	lb := make([]byte, 2)
	binary.LittleEndian.PutUint16(lb, 255)
	if g, w := lb, []byte{255, 0}; !reflect.DeepEqual(g, w) {
		t.Errorf(" got %v, want %v", g, w)
	}

	var bi uint16
	binary.Read(bytes.NewReader(bb), binary.BigEndian, &bi)

	var li uint16
	binary.Read(bytes.NewReader(lb), binary.LittleEndian, &li)

	if bi != 255 || li != 255 {
		t.Errorf("got bi:%v li:%v, want %v", bi, li, 255)
	}
}
