package binary

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestConvertIntToBytes(t *testing.T) {
	b := make([]byte, 2)

	binary.BigEndian.PutUint16(b, 255)
	if g, w := b, []byte{0, 255}; !reflect.DeepEqual(g, w) {
		t.Errorf("b got %v, want %v", g, w)
	}

	binary.BigEndian.PutUint16(b, 256)
	if g, w := b, []byte{1, 0}; !reflect.DeepEqual(g, w) {
		t.Errorf("b got %v, want %v", g, w)
	}

	binary.LittleEndian.PutUint16(b, 255)
	if g, w := b, []byte{255, 0}; !reflect.DeepEqual(g, w) {
		t.Errorf("b got %v, want %v", g, w)
	}

	binary.LittleEndian.PutUint16(b, 256)
	if g, w := b, []byte{0, 1}; !reflect.DeepEqual(g, w) {
		t.Errorf("b got %v, want %v", g, w)
	}

	// byte(=uint8)の最大値は255
	//
	// 現在主流のx86, x86_64はリトルエンディアン
	// ARMはどちらでも切り替えられるがAndroid, iOSはリトルエンディアン
	//
	// しかし、ネットワーク上で転送されるデータの多くは、大きい桁からメモリに格納されるビックエンディアンである
	// 転送データのrawバイナリを解析する際はエンディアン変換が必要になる
}

func TestConvertStringToBytes(t *testing.T) {
	s1 := "テスト"
	b := []byte(*&s1)
	s2 := string(b)

	if g, w := s1, s2; g != w {
		t.Errorf("string got %v, want %v", g, w)
	}
}
