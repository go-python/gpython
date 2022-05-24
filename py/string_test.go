// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

import "testing"

func TestStringFind(t *testing.T) {
	for _, tc := range []struct {
		str string
		sub string
		beg int
		end int
		idx int
	}{
		{
			str: "hello world",
			sub: "world",
			idx: 6,
		},
		{
			str: "hello world",
			sub: "o",
			idx: 4,
		},
		{
			str: "hello world",
			sub: "o",
			beg: 5,
			idx: 7,
		},
		{
			str: "hello world",
			sub: "bye",
			idx: -1,
		},
		{
			str: "Hello, 世界",
			sub: "界",
			idx: 8,
		},
		{
			str: "01234 6789",
			sub: " ",
			beg: 6,
			idx: -1,
		},
		{
			str: "0123456789",
			sub: "6",
			beg: 1,
			end: 6,
			idx: -1,
		},
		{
			str: "0123456789",
			sub: "6",
			beg: 1,
			end: 7,
			idx: 6,
		},
		{
			str: "0123456789",
			sub: "6",
			beg: 1,
			end: -1,
			idx: 6,
		},
		{
			str: "0123456789",
			sub: "6",
			beg: 100,
			end: -1,
			idx: -1,
		},
		{
			str: "0123456789",
			sub: "6",
			beg: 2,
			end: 1,
			idx: -1,
		},
	} {
		t.Run(tc.str+":"+tc.sub, func(t *testing.T) {
			beg := tc.beg
			end := tc.end
			if end == 0 {
				end = len(tc.str)
			}
			idx, err := String(tc.str).find(Tuple{String(tc.sub), Int(beg), Int(end)})
			if err != nil {
				t.Fatalf("invalid: %+v", err)
			}
			if got, want := int(idx.(Int)), tc.idx; got != want {
				t.Fatalf("got=%d, want=%d", got, want)
			}
		})
	}
}
