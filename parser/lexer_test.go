package parser

import (
	"bytes"
	"log"
	"math"
	"testing"

	"github.com/ncw/gpython/py"
)

func TestLexerIsIdentifier(t *testing.T) {
	for _, test := range []struct {
		in    rune
		start bool
		char  bool
	}{
		{'a', true, true},
		{'r', true, true},
		{'z', true, true},
		{'A', true, true},
		{'R', true, true},
		{'Z', true, true},
		{'0', false, true},
		{'4', false, true},
		{'9', false, true},
		{'_', true, true},
		{'@', false, false},
		{'[', false, false},
		{' ', false, false},
		{'\t', false, false},
		{'µ', true, true},
		{'©', false, false},
		{'—', false, false},
	} {
		got := isIdentifierStart(test.in)
		if got != test.start {
			t.Errorf("isIdentifierStart(%q) got %v expected %v", test.in, got, test.start)
		}
		got = isIdentifierChar(test.in)
		if got != test.char {
			t.Errorf("isIdentifierChar(%q) got %v expected %v", test.in, got, test.char)
		}
	}

}

func TestLexerReadIdentifier(t *testing.T) {
	x := yyLex{}
	for _, test := range []struct {
		in       string
		expected string
		out      string
	}{
		{"", "", ""},
		{"1", "", "1"},
		{"potato", "potato", ""},
		{"break+", "break", "+"},
		{"_aAzZ09ß²¹", "_aAzZ09ß", "²¹"},
		{"123abc", "", "123abc"},
		{" abc", "", " abc"},
		{"+abc", "", "+abc"},
	} {
		x.line = test.in
		got := x.readIdentifier()
		if got != test.expected || x.line != test.out {
			t.Errorf("readIdentifier(%q) got %q remainder %q, expected %q remainder %q", test.in, got, x.line, test.expected, test.out)
		}
	}
}

func TestLexerReadIdentifierOrKeyword(t *testing.T) {
	x := yyLex{}
	for _, test := range []struct {
		in    string
		token int
		value string
		out   string
	}{
		{"", eof, "", ""},
		{"1", eof, "", "1"},
		{"potato", NAME, "potato", ""},
		{"break+", BREAK, "break", "+"},
		{"breaker+", NAME, "breaker", "+"},
		{"_aAzZ09ß²¹", NAME, "_aAzZ09ß", "²¹"},
		{"123abc", eof, "", "123abc"},
		{" abc", eof, "", " abc"},
		{"+abc", eof, "", "+abc"},
	} {
		x.line = test.in
		token, value := x.readIdentifierOrKeyword()
		if token != test.token || value != test.value || x.line != test.out {
			t.Errorf("readIdentifierOrKeyword(%q) got (%q,%q) remainder %q, expected (%q,%q) remainder %q", test.in, tokenToString[token], value, x.line, tokenToString[test.token], test.value, test.out)
		}
	}
}

func TestLexerReadOperator(t *testing.T) {
	x := yyLex{}
	for _, test := range []struct {
		in       string
		expected int
		out      string
	}{
		{"", eof, ""},
		{" <", eof, " <"},
		{"<", '<', ""},
		{"<+", '<', "+"},
		{"<< ", LTLT, " "},
		{"<<=", LTLTEQ, ""},
		{"<<< ", LTLT, "< "},
		{"<==", LTEQ, "="},
		{"/", '/', ""},
		{"//", DIVDIV, ""},
		{"=//", '=', "//"},
		{"//=", DIVDIVEQ, ""},
		{"....", ELIPSIS, "."},
	} {
		x.line = test.in
		got := x.readOperator()
		if got != test.expected || x.line != test.out {
			t.Errorf("readOperator(%q) got %q remainder %q, expected %q remainder %q", test.in, tokenToString[got], x.line, tokenToString[test.expected], test.out)
		}
	}
}

// Whether two floats are more or less the same
func approxEq(a, b float64) bool {
	log.Printf("ApproxEq(a = %#v, b = %#v)", a, b)
	diff := a - b
	log.Printf("ApproxEq(diff = %e)", diff)
	if math.Abs(diff) > 1E-10 {
		log.Printf("ApproxEq(false)")
		return false
	}
	log.Printf("ApproxEq(true)")
	return true
}

func TestLexerReadNumber(t *testing.T) {
	x := yyLex{}
	for _, test := range []struct {
		in    string
		token int
		value py.Object
		out   string
	}{
		{"", eof, nil, ""},
		{"break", eof, py.Object(nil), "break"},

		{"0o0", NUMBER, py.Int(0), ""},
		{"0O765a", NUMBER, py.Int(0765), "a"},
		{"0o0007779", NUMBER, py.Int(0777), "9"},

		{"0x0", NUMBER, py.Int(0), ""},
		{"0XaBcDeFg", NUMBER, py.Int(0xABCDEF), "g"},
		{"0x000123z", NUMBER, py.Int(0x123), "z"},

		{"0b0", NUMBER, py.Int(0), ""},
		{"0B100", NUMBER, py.Int(4), ""},
		{"0B0001112", NUMBER, py.Int(7), "2"},

		{"1.", NUMBER, py.Float(1), ""},
		{".1", NUMBER, py.Float(.1), ""},
		{"0.1", NUMBER, py.Float(0.1), ""},
		{"00000.10000", NUMBER, py.Float(0.1), ""},
		{"1.E1", NUMBER, py.Float(10), ""},
		{".1e1", NUMBER, py.Float(1), ""},
		{"0.1e-01", NUMBER, py.Float(0.01), ""},
		{"00000.10000E-000001", NUMBER, py.Float(0.01), ""},
		{"1.j", NUMBER, py.Complex(complex(0, 1)), ""},
		{".1j", NUMBER, py.Complex(complex(0, .1)), ""},
		{"0.1j", NUMBER, py.Complex(complex(0, 0.1)), ""},
		{"00000.10000j", NUMBER, py.Complex(complex(0, 0.1)), ""},
		{"1.E1j", NUMBER, py.Complex(complex(0, 10)), ""},
		{".1e1j", NUMBER, py.Complex(complex(0, 1)), ""},
		{"0.1e-01j", NUMBER, py.Complex(complex(0, 0.01)), ""},
		{"00000.10000E-000001j", NUMBER, py.Complex(complex(0, 0.01)), ""},

		{"1", NUMBER, py.Int(1), ""},
		{"1+2", NUMBER, py.Int(1), "+2"},
		{"01", eofError, nil, "01"},
		{"00", NUMBER, py.Int(0), ""},
		{"123", NUMBER, py.Int(123), ""},
		{"0123", eofError, nil, "0123"},
		{"0123j", NUMBER, py.Complex(complex(0, 123)), ""},
		{"00j", NUMBER, py.Complex(complex(0, 0)), ""},
	} {
		x.line = test.in
		token, value := x.readNumber()
		if token != test.token || value != test.value || x.line != test.out {
			t.Errorf("readNumber(%q) got (%q,%T,%#v) remainder %q, expected (%q,%T,%#v) remainder %q", test.in, tokenToString[token], value, value, x.line, tokenToString[test.token], test.value, test.value, test.out)
		}
	}
}

func TestLexerReadString(t *testing.T) {
	for _, test := range []struct {
		in    string
		token int
		value py.Object
		out   string
	}{
		{``, eof, nil, ``},

		{`""a`, STRING, py.String(``), `a`},
		{`u"abc"`, STRING, py.String(`abc`), ``},
		{`"a\nc"`, STRING, py.String(`a\nc`), ``},
		{`r"a\nc"`, STRING, py.String(`a\nc`), ``},
		{`"a\"c"`, STRING, py.String(`a\"c`), ``},
		{`"a\\"+`, STRING, py.String(`a\\`), `+`},
		{`"a`, eofError, nil, `a`},

		{`''a`, STRING, py.String(``), `a`},
		{`U'abc'`, STRING, py.String(`abc`), ``},
		{`'a\nc'`, STRING, py.String(`a\nc`), ``},
		{`R'a\nc'`, STRING, py.String(`a\nc`), ``},
		{`'a\'c'`, STRING, py.String(`a\'c`), ``},
		{`'\n`, eofError, nil, `\n`},
		{`'a`, eofError, nil, `a`},

		{`""""""a`, STRING, py.String(``), `a`},
		{`u"""abc"""`, STRING, py.String(`abc`), ``},
		{`"""a\nc"""`, STRING, py.String(`a\nc`), ``},
		{`r"""a\"""c"""`, STRING, py.String(`a\"""c`), ``},
		{`"""a\"""c"""`, STRING, py.String(`a\"""c`), ``},
		{`"""a`, eofError, nil, `a`},
		{"\"\"\"a\nb\nc\n\"\"\"\n", STRING, py.String("a\nb\nc\n"), "\n"},
		{"\"\"\"a\nb\nc\na", eofError, nil, "a"},

		{`''''''a`, STRING, py.String(``), `a`},
		{`U'''abc'''`, STRING, py.String(`abc`), ``},
		{`'''a\nc'''`, STRING, py.String(`a\nc`), ``},
		{`R'''a\nc'''`, STRING, py.String(`a\nc`), ``},
		{`'''a\'''c'''`, STRING, py.String(`a\'''c`), ``},
		{`'''a`, eofError, nil, `a`},
		{"'''a\nb\nc\n'''\n", STRING, py.String("a\nb\nc\n"), "\n"},
		{"'''a\nb\nc\na", eofError, nil, "a"},

		{`b""a`, STRING, py.Bytes{}, "a"},
		{`b'abc'`, STRING, py.Bytes(string(`abc`)), ``},
		{`B"""a\nc"""`, STRING, py.Bytes(string(`a\nc`)), ``},
		{`B'''a\"c'''`, STRING, py.Bytes(string(`a\"c`)), ``},

		{`rb""a`, STRING, py.Bytes{}, "a"},
		{`bR'abc'`, STRING, py.Bytes(string(`abc`)), ``},
		{`BR"""a\nc"""`, STRING, py.Bytes(string(`a\nc`)), ``},
		{`rB'''a\"c'''`, STRING, py.Bytes(string(`a\"c`)), ``},
	} {
		x := NewLex(bytes.NewBufferString(test.in))
		x.refill()
		token, value := x.readString()
		equal := false
		if valueBytes, ok := value.(py.Bytes); ok {
			if testValueBytes, ok := test.value.(py.Bytes); !ok {
				t.Error("Expecting py.Bytes")
			} else {
				equal = (bytes.Compare(valueBytes, testValueBytes) == 0)
			}
		} else {
			equal = (value == test.value)
		}

		if token != test.token || !equal || x.line != test.out {
			t.Errorf("readString(%q) got (%q,%T,%#v) remainder %q, expected (%q,%T,%#v) remainder %q", test.in, tokenToString[token], value, value, x.line, tokenToString[test.token], test.value, test.value, test.out)
		}
	}
}
