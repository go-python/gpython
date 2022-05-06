// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

import (
	"fmt"
	"testing"
)

func TestParseTupleAndKeywords(t *testing.T) {
	for _, tc := range []struct {
		args    Tuple
		kwargs  StringDict
		format  string
		kwlist  []string
		results []Object
		err     error
	}{
		{
			args:    Tuple{String("a")},
			format:  "O:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{None},
			format:  "Z:func",
			results: []Object{None},
		},
		{
			args:    Tuple{String("a")},
			format:  "Z:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{Int(42)},
			format:  "Z:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str or None, not int'"),
		},
		{
			args:    Tuple{None},
			format:  "Z*:func", // FIXME(sbinet): invalid format.
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str or None, not NoneType'"),
		},
		{
			args:    Tuple{None},
			format:  "Z#:func",
			results: []Object{None},
		},
		{
			args:    Tuple{String("a")},
			format:  "Z#:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{Int(42)},
			format:  "Z#:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str or None, not int'"),
		},
		{
			args:    Tuple{None},
			format:  "z:func",
			results: []Object{None},
		},
		{
			args:    Tuple{String("a")},
			format:  "z:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{Int(42)},
			format:  "z:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str or None, not int'"),
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "z:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str or None, not bytes'"),
		},
		{
			args:    Tuple{None},
			format:  "z*:func",
			results: []Object{None},
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "z*:func",
			results: []Object{Bytes("a")},
		},
		{
			args:    Tuple{String("a")},
			format:  "z*:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{Int(42)},
			format:  "z*:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str, bytes-like or None, not int'"),
		},
		{
			args:    Tuple{None},
			format:  "z#:func",
			results: []Object{None},
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "z#:func",
			results: []Object{Bytes("a")},
		},
		{
			args:    Tuple{String("a")},
			format:  "z#:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{Int(42)},
			format:  "z#:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str, bytes-like or None, not int'"),
		},
		{
			args:    Tuple{String("a")},
			format:  "s:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{None},
			format:  "s:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str, not NoneType'"),
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "s:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str, not bytes'"),
		},
		{
			args:    Tuple{String("a")},
			format:  "s#:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "s#:func",
			results: []Object{Bytes("a")},
		},
		{
			args:    Tuple{None},
			format:  "s#:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str or bytes-like, not NoneType'"),
		},
		{
			args:    Tuple{String("a")},
			format:  "s*:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "s*:func",
			results: []Object{Bytes("a")},
		},
		{
			args:    Tuple{None},
			format:  "s*:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str or bytes-like, not NoneType'"),
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "y:func",
			results: []Object{Bytes("a")},
		},
		{
			args:    Tuple{None},
			format:  "y:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be bytes-like, not NoneType'"),
		},
		{
			args:    Tuple{String("a")},
			format:  "y:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be bytes-like, not str'"),
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "y#:func",
			results: []Object{Bytes("a")},
		},
		{
			args:    Tuple{String("a")},
			format:  "y#:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be bytes-like, not str'"),
		},
		{
			args:    Tuple{None},
			format:  "y#:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be bytes-like, not NoneType'"),
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "y*:func",
			results: []Object{Bytes("a")},
		},
		{
			args:    Tuple{String("a")},
			format:  "y*:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be bytes-like, not str'"),
		},
		{
			args:    Tuple{None},
			format:  "y*:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be bytes-like, not NoneType'"),
		},
		{
			args:    Tuple{String("a")},
			format:  "U:func",
			results: []Object{String("a")},
		},
		{
			args:    Tuple{None},
			format:  "U:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str, not NoneType'"),
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "U:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str, not bytes'"),
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "U*:func", // FIXME(sbinet): invalid format
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str, not bytes'"),
		},
		{
			args:    Tuple{Bytes("a")},
			format:  "U#:func", // FIXME(sbinet): invalid format
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be str, not bytes'"),
		},
		{
			args:    Tuple{Int(42)},
			format:  "i:func",
			results: []Object{Int(42)},
		},
		{
			args:    Tuple{None},
			format:  "i:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be int, not NoneType'"),
		},
		{
			args:    Tuple{Int(42)},
			format:  "n:func",
			results: []Object{Int(42)},
		},
		{
			args:    Tuple{None},
			format:  "n:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be int, not NoneType'"),
		},
		{
			args:    Tuple{Bool(true)},
			format:  "p:func",
			results: []Object{Bool(true)},
		},
		{
			args:    Tuple{None},
			format:  "p:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be bool, not NoneType'"),
		},
		{
			args:    Tuple{Float(42)},
			format:  "d:func",
			results: []Object{Float(42)},
		},
		{
			args:    Tuple{Int(42)},
			format:  "d:func",
			results: []Object{Float(42)},
		},
		{
			args:    Tuple{None},
			format:  "p:func",
			results: []Object{nil},
			err:     fmt.Errorf("TypeError: 'func() argument 1 must be bool, not NoneType'"),
		},
	} {
		t.Run(tc.format, func(t *testing.T) {
			results := make([]*Object, len(tc.results))
			for i := range tc.results {
				results[i] = &tc.results[i]
			}
			err := ParseTupleAndKeywords(tc.args, tc.kwargs, tc.format, tc.kwlist, results...)
			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %s\nwant=%s", got, want)
				}
			case err != nil && tc.err == nil:
				t.Fatalf("could not parse tuple+kwargs: %+v", err)
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error (got=nil): %+v", tc.err)
			case err == nil && tc.err == nil:
				// ok.
			}
			// FIXME(sbinet): check results
		})
	}
}
