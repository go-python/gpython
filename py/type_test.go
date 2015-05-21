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
