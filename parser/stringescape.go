// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"strconv"

	"github.com/go-python/gpython/py"
)

// DecodeEscape unescapes a backslash-escaped buffer
//
// byteMode indicates whether we are creating a unicode string or a bytes output
func DecodeEscape(in *bytes.Buffer, byteMode bool) (out *bytes.Buffer, err error) {
	// Early exit if no escape sequences
	// NB in.Bytes() is cheap
	inBytes := in.Bytes()
	if bytes.IndexRune(inBytes, '\\') < 0 {
		return in, nil
	}
	out = new(bytes.Buffer)
	runes := bytes.Runes(inBytes)
	decodeHex := func(what byte, i, size int) error {
		i++
		if i+size <= len(runes) {
			cout, err := strconv.ParseInt(string(runes[i:i+size]), 16, 32)
			if err != nil {
				return py.ExceptionNewf(py.ValueError, "invalid \\%c escape at position %d", what, i-2)
			}
			if byteMode {
				out.WriteByte(byte(cout))
			} else {
				out.WriteRune(rune(cout))
			}
		} else {
			return py.ExceptionNewf(py.ValueError, "truncated \\%c escape at position %d", what, i-2)
		}
		return nil
	}
	ignoreEscape := false
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if c != '\\' {
			out.WriteRune(c)
			continue
		}
		i++
		if i >= len(runes) {
			return nil, py.ExceptionNewf(py.ValueError, "Trailing \\ in string")
		}
		c = runes[i]
		switch c {
		case '\n':
		case '\\':
			out.WriteRune('\\')
		case '\'':
			out.WriteRune('\'')
		case '"':
			out.WriteRune('"')
		case 'b':
			out.WriteRune('\b')
		case 'f':
			out.WriteRune('\014') // FF
		case 't':
			out.WriteRune('\t')
		case 'n':
			out.WriteRune('\n')
		case 'r':
			out.WriteRune('\r')
		case 'v':
			out.WriteRune('\013') // VT
		case 'a':
			out.WriteRune('\007') // BEL, not classic C
		case '0', '1', '2', '3', '4', '5', '6', '7':
			// 1 to 3 characters of octal escape
			cout := c - '0'
			if i+1 < len(runes) && '0' <= runes[i+1] && runes[i+1] <= '7' {
				i++
				cout = (cout << 3) + runes[i] - '0'
				if i+1 < len(runes) && '0' <= runes[i+1] && runes[i+1] <= '7' {
					i++
					cout = (cout << 3) + runes[i] - '0'
				}
			}
			if byteMode {
				out.WriteByte(byte(cout))
			} else {
				out.WriteRune(cout)
			}
		case 'x':
			// \xhh exactly 2 characters of hex
			err = decodeHex('x', i, 2)
			if err != nil {
				return nil, err
			}
			i += 2
			// FIXME In a bytes literal, hexadecimal and
			// octal escapes denote the byte with the
			// given value. In a string literal, these
			// escapes denote a Unicode character with the
			// given value.
		case 'u':
			// \uxxxx	Character with 16-bit hex value xxxx - 4 characters required
			if byteMode {
				ignoreEscape = true
				break
			}
			err = decodeHex('u', i, 4)
			if err != nil {
				return nil, err
			}
			i += 4
		case 'U':
			// \Uxxxxxxxx	Character with 32-bit hex value xxxxxxxx - 8 characters required
			if byteMode {
				ignoreEscape = true
				break
			}

			err = decodeHex('U', i, 8)
			if err != nil {
				return nil, err
			}
			i += 8
		case 'N':
			// \N{name}	Character named name in the Unicode database
			if byteMode {
				ignoreEscape = true
				break
			}
			// FIXME go can't do this as builtin so ignore for the moment
			ignoreEscape = true
		default:
			ignoreEscape = true
			break
		}
		// ignore unrecognised escape
		if ignoreEscape {
			i--
			out.WriteRune('\\')
			ignoreEscape = false
		}
	}
	return out, nil
}
