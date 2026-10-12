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

func TestStringUpper(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want Object
	}{{
		name: "abc",
		s:    String("abc"),
		want: String("ABC")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Upper()
			if err != nil {
				t.Fatalf("Upper() error = %v", err)
			}
			if got.(String) != tt.want.(String) {
				t.Fatalf("Upper() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringLower(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want Object
	}{{
		name: "ABC",
		s:    String("ABC"),
		want: String("abc")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Lower()
			if err != nil {
				t.Fatalf("Lower() error = %v", err)
			}
			if got.(String) != tt.want.(String) {
				t.Fatalf("Lower() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSplitLimits(t *testing.T) {
	for _, limit := range []Int{-100, -2, -1, 0, 1, 2, Int(^uint(0) >> 1)} {
		got, err := String("a,b,c").Split(Tuple{String(","), limit}, nil)
		if err != nil {
			t.Fatal(err)
		}
		want := Tuple{String("a"), String("b"), String("c")}
		if limit == 0 {
			want = Tuple{String("a,b,c")}
		} else if limit == 1 {
			want = Tuple{String("a"), String("b,c")}
		}
		items := got.(*List).Items
		if len(items) != len(want) {
			t.Errorf("limit %d: got %v, want %v", limit, items, want)
			continue
		}
		for i := range want {
			if items[i] != want[i] {
				t.Errorf("limit %d: got %v, want %v", limit, items, want)
			}
		}
	}
	for _, s := range []String{"", "abc"} {
		for _, limit := range []Int{-1, 0, 1} {
			_, err := s.Split(Tuple{String(""), limit}, nil)
			if !IsException(ValueError, err) {
				t.Errorf("%q split empty with limit %d: %v, want ValueError", s, limit, err)
			}
		}
	}
}
