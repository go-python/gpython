// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

import "testing"

func TestIsSubType(t *testing.T) {
	for _, test := range []struct {
		a    *Type
		b    *Type
		want bool
	}{
		{ValueError, ValueError, true},
		{ValueError, ExceptionType, true},
		{ExceptionType, ValueError, false},
	} {
		got := test.a.IsSubtype(test.b)
		if test.want != got {
			t.Errorf("%v.IsSubtype(%v) want %v got %v", test.a.Name, test.b.Name, test.want, got)
		}
	}
}

func TestMro(t *testing.T) {
	for _, test := range []struct {
		t    *Type
		want []*Type
	}{
		{ObjectType, []*Type{ObjectType}},
		{BaseException, []*Type{BaseException, ObjectType}},
		{ExceptionType, []*Type{ExceptionType, BaseException, ObjectType}},
		{ValueError, []*Type{ValueError, ExceptionType, BaseException, ObjectType}},
	} {
		got := test.t.Mro
		if len(test.want) != len(got) {
			t.Errorf("differing lengths: want %v, got %v", test.want, got)
		} else {
			for i := range got {
				baseGot := got[i].(*Type)
				baseWant := test.want[i]
				if baseGot != baseWant {
					t.Errorf("mro[%d] want %s got %s", i, baseWant.Name, baseGot.Name)
				}

			}
		}
	}
}
