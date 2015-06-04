package parser

import (
	"bytes"
	"testing"

	"github.com/ncw/gpython/py"
)

func TestDecodeEscape(t *testing.T) {
	for _, test := range []struct {
		in        string
		want      string
		errString string
		byteMode  bool
	}{
		// Stringmode tests
		{``, ``, "", false},
		{`Potato`, `Potato`, "", false},
		{`Potato\`, ``, `Trailing \ in string`, false},
		{`\Potato`, `\Potato`, "", false},
		{`n\\`, `n\`, "", false},
		{`\'x`, `'x`, "", false},
		{`\"`, `"`, "", false},
		{"\\\n", ``, "", false},
		{`\b`, "\010", "", false},
		{`\f`, "\014", "", false},
		{`\t`, "\011", "", false},
		{`\n`, "\012", "", false},
		{`\r`, "\015", "", false},
		{`\v`, "\013", "", false},
		{`\a`, "\007", "", false},
		{`\1`, "\001", "", false},
		{`\12`, "\012", "", false},
		{`\123`, "\123", "", false},
		{`\777`, "\u01ff", "", false},
		{`\1\12\123\1234`, "\001\012\123\123" + "4", "", false},
		{`a\1a\12a\123a`, "a\001a\012a\123a", "", false},
		{`\x`, "", `truncated \x escape at position 0`, false},
		{`\x1`, "", `truncated \x escape at position 0`, false},
		{`\x11`, "\x11", "", false},
		{`\xzz`, "", `invalid \x escape at position 0`, false},
		{`{\x11}`, "{\x11}", "", false},
		{`\x01\x8a\xff`, "\x01\u008a\u00ff", "", false},
		{`\x01\x8A\xFF`, "\x01\u008a\u00ff", "", false},
		{`\u`, "", `truncated \u escape at position 0`, false},
		{`\u1`, "", `truncated \u escape at position 0`, false},
		{`\u12`, "", `truncated \u escape at position 0`, false},
		{`z\u134`, "", `truncated \u escape at position 1`, false},
		{`\u1234`, "\u1234", "", false},
		{`z\uzzzz`, "", `invalid \u escape at position 1`, false},
		{`{\u1234}`, "{\u1234}", "", false},
		{`\U00000001\U0000018a\U000012ff`, "\U00000001\U0000018a\U000012ff", "", false},
		{`\U00000001\U0000018A\U000012FF`, "\U00000001\U0000018a\U000012ff", "", false},
		{`\U0000`, "", `truncated \U escape at position 0`, false},
		{`\U00001`, "", `truncated \U escape at position 0`, false},
		{`\U000012`, "", `truncated \U escape at position 0`, false},
		{`z\U0000134`, "", `truncated \U escape at position 1`, false},
		{`\U00001234`, "\U00001234", "", false},
		{`z\Uzzzzzzzz`, "", `invalid \U escape at position 1`, false},
		{`{\U00001234}`, "{\U00001234}", "", false},
		{`\U00000001\U0000018a\U000012ff`, "\U00000001\U0000018a\U000012ff", "", false},
		{`\U00000001\U0000018A\U000012FF`, "\U00000001\U0000018a\U000012ff", "", false},
		{`\N{potato}`, `\N{potato}`, "", false},

		// Bytemode tests
		{``, ``, "", true},
		{`Potato`, `Potato`, "", true},
		{`Potato\`, ``, `Trailing \ in string`, true},
		{`\Potato`, `\Potato`, "", true},
		{`n\\`, `n\`, "", true},
		{`\'x`, `'x`, "", true},
		{`\"`, `"`, "", true},
		{"\\\n", ``, "", true},
		{`\b`, "\010", "", true},
		{`\f`, "\014", "", true},
		{`\t`, "\011", "", true},
		{`\n`, "\012", "", true},
		{`\r`, "\015", "", true},
		{`\v`, "\013", "", true},
		{`\a`, "\007", "", true},
		{`\1`, "\001", "", true},
		{`\12`, "\012", "", true},
		{`\123`, "\123", "", true},
		{`\777`, "\xff", "", true},
		{`\1\12\123\1234`, "\001\012\123\123" + "4", "", true},
		{`a\1a\12a\123a`, "a\001a\012a\123a", "", true},
		{`\x`, "", `truncated \x escape at position 0`, true},
		{`\x1`, "", `truncated \x escape at position 0`, true},
		{`\x11`, "\x11", "", true},
		{`\xzz`, "", `invalid \x escape at position 0`, true},
		{`{\x11}`, "{\x11}", "", true},
		{`\x01\x8a\xff`, "\x01\x8a\xff", "", true},
		{`\x01\x8A\xFF`, "\x01\x8a\xff", "", true},
		{`\u`, `\u`, "", true},
		{`\u1`, `\u1`, "", true},
		{`\u12`, `\u12`, "", true},
		{`z\u134`, `z\u134`, "", true},
		{`\u1234`, `\u1234`, "", true},
		{`z\uzzzz`, `z\uzzzz`, "", true},
		{`{\u1234}`, `{\u1234}`, "", true},
		{`\U00000001\U0000018a\U000012ff`, `\U00000001\U0000018a\U000012ff`, "", true},
		{`\U00000001\U0000018A\U000012FF`, `\U00000001\U0000018A\U000012FF`, "", true},
		{`\U0000`, `\U0000`, "", true},
		{`\U00001`, `\U00001`, "", true},
		{`\U000012`, `\U000012`, "", true},
		{`z\U0000134`, `z\U0000134`, "", true},
		{`\U00001234`, `\U00001234`, "", true},
		{`z\Uzzzzzzzz`, `z\Uzzzzzzzz`, "", true},
		{`{\U00001234}`, `{\U00001234}`, "", true},
		{`\U00000001\U0000018a\U000012ff`, `\U00000001\U0000018a\U000012ff`, "", true},
		{`\U00000001\U0000018A\U000012FF`, `\U00000001\U0000018A\U000012FF`, "", true},
		{`\N{potato}`, `\N{potato}`, "", true},
	} {
		in := bytes.NewBufferString(test.in)
		out, err := DecodeEscape(in, test.byteMode)
		if err != nil {
			if test.errString == "" {
				t.Errorf("%q: not expecting error but got: %v", test.in, err)
			} else {
				exc := err.(*py.Exception)
				args := exc.Args.(py.Tuple)
				if string(args[0].(py.String)) != test.errString {
					t.Errorf("%q: want error %q but got %q", test.in, test.errString, args[0])
				}
			}
			continue
		}
		if test.errString != "" {
			t.Errorf("%q: expecting error but didn't get one", test.in)
			continue
		}
		got := out.String()
		if test.want != got {
			t.Errorf("%q: want %q but got %q", test.in, test.want, got)
		}
	}
}
