// Copyright 2026 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

import (
	"fmt"
	"testing"
)

func TestStringSearchBounds(t *testing.T) {
	// Expected results were checked against CPython.
	for _, tc := range []struct {
		s, sub                  string
		start, end, find, count int
	}{
		{"abcabc", "a", -3, 6, 3, 1},
		{"abcabc", "a", 0, -1, 0, 2},
		{"abcabc", "c", 0, -1, 2, 1},
		{"abcabc", "a", -100, 100, 0, 2},
		{"abcabc", "a", 0, -100, -1, 0},
		{"abcabc", "", 7, 100, -1, 0},
		{"abcabc", "", 4, 2, -1, 0},
		{"abcabc", "", 6, 6, 6, 1},
		{"abcabc", "", 0, -1, 0, 6},
		{"é世界é", "é", -1, 4, 3, 1},
		{"é世界é", "界", 0, -1, 2, 1},
		{"é世界é", "", -3, -1, 1, 3},
		{"", "", 0, 0, 0, 1},
		{"", "", 1, 1, -1, 0},
	} {
		for _, method := range []string{"find", "count"} {
			t.Run(fmt.Sprintf("%s/%q/%q/%d:%d", method, tc.s, tc.sub, tc.start, tc.end), func(t *testing.T) {
				defer func() {
					if p := recover(); p != nil {
						t.Errorf("unexpected panic: %v", p)
					}
				}()
				args := Tuple{String(tc.sub), Int(tc.start), Int(tc.end)}
				var got Object
				var err error
				want := tc.find
				if method == "find" {
					got, err = String(tc.s).find(args)
				} else {
					got, err = String(tc.s).Count(args)
					want = tc.count
				}
				if err != nil || got != Int(want) {
					t.Errorf("got %v, %v; want %d", got, err, want)
				}
			})
		}
	}
}
