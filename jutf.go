// Package jutf implements the modified UTF-8 scheme as used by Java.
package jutf

import (
	"bytes"
)

//
// https://docs.oracle.com/javase/8/docs/api/java/io/DataInput.html#modified-utf-8
//
// The null byte '\u0000' is encoded in 2-byte format rather than 1-byte, so that the encoded strings never have embedded nulls.
// Only the 1-byte, 2-byte, and 3-byte formats are used.
// Supplementary characters are represented in the form of surrogate pairs.
//
// Modified UTF-8:
//
// Range (hex)          Encoding (binary)
// 0                    11000000 00000000
// $1 .. $7f            0_______
// $80 .. $7ff          110_____ 10______
// $800 .. $ffff        1110____ 10______ 10______
// $10000 .. $10ffff    6 byte (2 surrogates)
//
// Example: U+10437 to surrogate pair:
// 1. $10437 - $10000 = $437
// 2. ($437 shr  10)  + $d800 = $d801
// 3. ($437 and $3ff) + $dc00 = $dc37
//

// Encode returns a string in modified UTF-8 format.
func Encode(s string) []byte {
	// Naive code, optimize later.
	buf := bytes.Buffer{}

	enc3 := func(r rune) {
		buf.WriteByte(byte(0xe0 | ((r >> 12) & 0xf)))
		buf.WriteByte(byte(0x80 | ((r >> 6) & 0x3f)))
		buf.WriteByte(byte(0x80 | (r & 0x3f)))
	}

	for _, r := range s {
		if r == 0 {
			buf.WriteByte(0xc0)
			buf.WriteByte(0x80)
		} else if r >= 1 && r <= 0x7f {
			buf.WriteByte(byte(r))
		} else if r >= 0x80 && r <= 0x7ff {
			buf.WriteByte(byte(0xc0 | (r >> 6)))
			buf.WriteByte(byte(0x80 | (r & 0x3f)))
		} else if r >= 0x800 && r <= 0xffff {
			enc3(r)
		} else if r >= 0x10000 && r <= 0x10ffff {
			enc3(((r - 0x10000) >> 10) + 0xd800)
			enc3(((r - 0x10000) & 0x3ff) + 0xdc00)
		} else {
			// outside utf-8 range, panic()?
		}
	}

	return buf.Bytes()
}

// Decode decodes the modified UTF-8 input in data to a string.
func Decode(data []byte) string {
	// TODO
	for i := 0; i < len(data); i++ {
		// 0xxxxxxx
		// 10xxxxxx
		// 110xxxxx 10xxxxxx
		// 1110xxxx 10xxxxxx 10xxxxxx
	}

	return ""
}
