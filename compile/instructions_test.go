// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compile

import (
	"bytes"
	"testing"
)

func TestLnotab(t *testing.T) {
	for i, test := range []struct {
		instrs Instructions
		want   []byte
	}{
		{
			instrs: Instructions{},
			want:   []byte{},
		},
		{
			instrs: Instructions{
				&Op{pos: pos{n: 1, p: 10, lineno: 1}},
				&Op{pos: pos{n: 0, p: 10, lineno: 0}},
				&Op{pos: pos{n: 1, p: 102, lineno: 1}},
			},
			want: []byte{},
		},
		{
			instrs: Instructions{
				&Op{pos: pos{n: 1, p: 0, lineno: 1}},
				&Op{pos: pos{n: 1, p: 1, lineno: 2}},
				&Op{pos: pos{n: 1, p: 2, lineno: 3}},
			},
			want: []byte{1, 1, 1, 1},
		},
		{
			// Example from lnotab.txt
			instrs: Instructions{
				&Op{pos: pos{n: 1, p: 0, lineno: 1}},
				&Op{pos: pos{n: 1, p: 6, lineno: 2}},
				&Op{pos: pos{n: 1, p: 50, lineno: 7}},
				&Op{pos: pos{n: 1, p: 350, lineno: 307}},
				&Op{pos: pos{n: 1, p: 361, lineno: 308}},
			},
			want: []byte{
				6, 1,
				44, 5,
				255, 0,
				45, 255,
				0, 45,
				11, 1},
		},
	} {
		got := test.instrs.Lnotab()
		if !bytes.Equal(test.want, got) {
			t.Errorf("%d: want %d got %d", i, test.want, got)
		}
	}
}
