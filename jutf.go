// Copyright 2019-2020 Anders Bergh <anders1@gmail.com>
// MIT license (see LICENSE).

// Package jutf implements the modified UTF-8 encoding used by Java.
package jutf

import (
	"bytes"
	"errors"
	"unicode/utf8"
)

// not yet exported.
var (
	errInvalidNUL        = errors.New("short NUL codepoint not allowed")
	errTooShort          = errors.New("unexpected end of data")
	errTooShortSurrogate = errors.New("unexpected end of data (missing surrogate)")
	errInvalidEncoding   = errors.New("invalid encoding")
)

//
// https://docs.oracle.com/javase/8/docs/api/java/io/DataInput.html#modified-utf-8
//
// The null byte '\u0000' is encoded in 2-byte format rather than 1-byte, so
// that the encoded strings never have embedded nulls.
// Only the 1-byte, 2-byte, and 3-byte formats are used.
// Supplementary characters are represented in the form of surrogate pairs.
//
// Modified UTF-8:
//
// Range                Encoding
// 0x0                  11000000 00000000
// 0x1 .. 0x7f          0_______
// 0x80 .. 0x7ff        110_____ 10______
// 0x800 .. 0xffff      1110____ 10______ 10______
// 0x10000 .. 0x10ffff  surrogate pair (6 bytes)
//
// Example: U+10437 to surrogate pair:
// 1. 0x10437 - 0x10000 = 0x437
// 2. (0x437 >> 10)   + 0xd800 = 0xd801
// 3. (0x437 & 0x3ff) + 0xdc00 = 0xdc37
//

// Encode returns a string in modified UTF-8 format.
func Encode(s string) []byte {
	buf := bytes.Buffer{}

	// Output will be at least as long as s, potentially longer
	buf.Grow(len(s))

	// up to 6 byte long encoding
	tmp := make([]byte, 6)

	for _, r := range s {
		if r == 0 {
			tmp[0] = 0xc0
			tmp[1] = 0x80

			buf.Write(tmp[0:2])
		} else if r >= 1 && r <= 0x7f {
			buf.WriteByte(byte(r))
		} else if r >= 0x80 && r <= 0x7ff {
			tmp[0] = byte(0xc0 | (r >> 6))
			tmp[1] = byte(0x80 | (r & 0x3f))

			buf.Write(tmp[0:2])
		} else if r >= 0x800 && r <= 0xffff {
			tmp[0] = byte(0xe0 | ((r >> 12) & 0xf))
			tmp[1] = byte(0x80 | ((r >> 6) & 0x3f))
			tmp[2] = byte(0x80 | (r & 0x3f))

			buf.Write(tmp[0:3])
		} else if r >= 0x10000 && r <= 0x10ffff {
			// codepoint 1
			r1 := ((r - 0x10000) >> 10) + 0xd800
			tmp[0] = byte(0xe0 | ((r1 >> 12) & 0xf))
			tmp[1] = byte(0x80 | ((r1 >> 6) & 0x3f))
			tmp[2] = byte(0x80 | (r1 & 0x3f))

			// codepoint 2
			r2 := ((r - 0x10000) & 0x3ff) + 0xdc00
			tmp[3] = byte(0xe0 | ((r2 >> 12) & 0xf))
			tmp[4] = byte(0x80 | ((r2 >> 6) & 0x3f))
			tmp[5] = byte(0x80 | (r2 & 0x3f))

			buf.Write(tmp[0:6])
		} else {
			// panic("out of range rune >0x10ffff")
			buf.Write([]byte("\ufffd")) // replacement character
		}
	}

	return buf.Bytes()
}

// Decode decodes the input array to a UTF-8 string.
func Decode(d []byte) (string, error) {
	// if the input already is a normal UTF-8 string, simply return it
	if utf8.ValidString(string(d)) {
		return string(d), nil
	}

	buf := bytes.Buffer{}
	buf.Grow(len(d)) // the final length of the output should be similar to the input.

	for i := 0; i < len(d); {
		if d[i] == 0 {
			// a short NUL, valid and reasonable except this is Java UTF-8.
			return "", errInvalidNUL
		} else if d[i] < 0x80 {
			// ASCII range, can simply copy it
			buf.WriteByte(d[i])
			i++
		} else if d[i]&0xe0 == 0xc0 {
			// 2 bytes
			if i+1 >= len(d) {
				return "", errTooShort
			}

			if d[i] == 0xc0 && d[i+1] == 0x80 {
				// "overlong" null
				buf.WriteByte(0)
			} else {
				// copy
				buf.WriteByte(d[i])
				buf.WriteByte(d[i+1])
			}

			i += 2
		} else if d[i]&0xf0 == 0xe0 {
			// 3 bytes
			if i+2 >= len(d) {
				return "", errTooShort
			}

			// surrogate pair, first codepoint
			if d[i] == 0xed && d[i+1] >= 0xa0 && d[i+1] <= 0xaf {
				// must be followed by a 3 byte codepoint
				if i+5 >= len(d) {
					return "", errTooShortSurrogate
				}

				// make sure the next codepoint is part of the surrogate pair
				if d[i+3] != 0xed || !(d[i+4] >= 0xb0 && d[i+4] <= 0xbf) {
					return "", errInvalidEncoding
				}

				// decode the whole surrogate pair
				c1 := int32(d[i]&0xf) << 12
				c1 |= int32(d[i+1]&0x3f) << 6
				c1 |= int32(d[i+2] & 0x3f)
				c2 := int32(d[i+3]&0xf) << 12
				c2 |= int32(d[i+4]&0x3f) << 6
				c2 |= int32(d[i+5] & 0x3f)
				cp := 0x10000 + ((c1 - 0xd800) << 10) | (c2 - 0xdc00)

				buf.WriteRune(rune(cp))
				i += 3 // eat the second half
			} else {
				// others can be copied
				buf.Write(d[i : i+3])
			}

			i += 3
		} else {
			// would be >3 bytes (invalid)
			return "", errInvalidEncoding
		}
	}

	return buf.String(), nil
}
