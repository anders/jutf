package jutf

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want []byte
	}{
		// Special Java outputs
		{"NUL", "\x00", []byte{0xc0, 0x80}},
		{"four byte", "\U0001f4a9", []byte{0xed, 0xa0, 0xbd, 0xed, 0xb2, 0xa9}},

		// Regular UTF-8 tests, can simply cast to []byte
		{"one byte", "ASCII", []byte("ASCII")},
		{"two byte", "åäö", []byte("åäö")},
		{"three byte", "日本語", []byte("日本語")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"NULL", []byte{0xc0, 0x80}, "\x00"},
		{"surrogate pair", []byte{0xed, 0xa0, 0xbd, 0xed, 0xb2, 0xa9}, "\U0001f4a9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode(tt.data); got != tt.want {
				t.Errorf("Decode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for n:=0; n < b.N; n++ {
		Encode("Hello\x00Wörld!!! \U0001f4a9")
	}
}
