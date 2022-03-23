# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import string

print("globals:")
for name in ("whitespace", 
			"ascii_lowercase",
			"ascii_uppercase",
			"ascii_letters",
			"digits",
			"hexdigits",
			"octdigits",
			"punctuation",
            "printable"):
    v = getattr(string, name)
    print("\nstring.%s:\n%s" % (name,repr(v)))

def assertEqual(x, y):
    assert x == y, "got: %s, want: %s" % (repr(x), repr(y))

assertEqual(string.capwords('abc def ghi'), 'Abc Def Ghi')
assertEqual(string.capwords('abc\tdef\nghi'), 'Abc Def Ghi')
assertEqual(string.capwords('abc\t   def  \nghi'), 'Abc Def Ghi')
assertEqual(string.capwords('ABC DEF GHI'), 'Abc Def Ghi')
assertEqual(string.capwords('ABC-DEF-GHI', '-'), 'Abc-Def-Ghi')
assertEqual(string.capwords('ABC-def DEF-ghi GHI'), 'Abc-def Def-ghi Ghi')
assertEqual(string.capwords('   aBc  DeF   '), 'Abc Def')
assertEqual(string.capwords('\taBc\tDeF\t'), 'Abc Def')
assertEqual(string.capwords('\taBc\tDeF\t', '\t'), '\tAbc\tDef\t')

