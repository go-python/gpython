# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaisesText, assertRaises

doc="str()"
assert str() == ""
assert str("hello") == "hello"
# FIXME assert str(b"hello") == "hello"
assert str(5) == "5"
assert str(5.1) == "5.1"
# FIXME assert str((1,2,3)) == "(1, 2, 3)"

class A():
    def __str__(self):
        return "str method"
    def __repr__(self):
        return "repr method"

class B():
    def __repr__(self):
        return "repr method"

class C():
    pass

assert str(A()) == "str method"
assert str(B()) == "repr method"
strC = str(C())
assert " at 0x" in strC
assert "<" in strC
assert ">" in strC

doc="repr()"
assert repr("") == "''"
assert repr("hello") == r"'hello'"
assert repr(r"""hel"lo""") == r"""'hel"lo'"""
assert repr("""he
llo""") == r"""'he\nllo'"""
assert repr(r"""hel'lo""") == r'''"hel'lo"'''
assert repr('\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\x7f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0¡¢£¤¥¦§¨©ª«¬\xad®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ') == r"""'\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\x7f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0¡¢£¤¥¦§¨©ª«¬\xad®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ'"""
assert repr('\u1000\uffff\U00010000\U0010ffff') == r"""'က\uffff𐀀\U0010ffff'"""

doc="comparison"
assert "" < "hello"
assert "HELLO" < "hello"
assert "hello" > "HELLO"
assert "HELLO" != "hello"
assert "hello" == "hello"
assert "HELLO" <= "hello"
assert "hello" <= "hello"
assert "hello" >= "HELLO"
assert "hello" <= "hello"
assertRaises(TypeError, lambda: 1 > "potato")
assertRaises(TypeError, lambda: 1 >= "potato")
assertRaises(TypeError, lambda: 1 < "potato")
assertRaises(TypeError, lambda: 1 <= "potato")
assert not( 1 == "potato")
assert 1 != "potato"

doc="startswith"
assert "HELLO THERE".startswith("HELL")
assert not "HELLO THERE".startswith("THERE")
assert "HELLO".startswith("LLO", 2)
assert "HELLO THERE".startswith(("HERE", "HELL"))

doc="endswith"
assert "HELLO THERE".endswith("HERE")
assert not "HELLO THERE".endswith("HELL")
assert "HELLO THERE".endswith(("HELL", "HERE"))

doc="bool"
assert "true"
assert not ""

doc="add"
a = "potato"
a = a + "sausage"
assert a == "potatosausage"
a = "potato"
a += "sausage"
assert a == "potatosausage"
a = "potato"
a = "sausage" + a
assert a == "sausagepotato"
assertRaises(TypeError, lambda: "sausage"+1)
assertRaises(TypeError, lambda: 1+"sausage")

doc="mul"
assert "a" * 0 == ""
assert "a" * -1 == ""
assert "a" * 5 == "aaaaa"
assertRaises(TypeError, lambda: "a" * 5.0)
assert 3 * "ab" == "ababab"
assertRaises(TypeError, lambda: 3.0 * "ab")
assertRaises(TypeError, lambda: "ab" * None)
a = "100"
a *= 2
assert a == "100100"

doc="in"
assert "x" in "hellox"
assert "x" not in "hello"
assert "el" in "hello"
assert "" in "hello"
assert "hello" not in ""
assertRaisesText(TypeError, "'in <string>' requires string as left operand, not int", lambda: 1 in "hello")

asc="hello"
uni="£100世界𠜎" # 1,2,3,4 byte unicode characters

doc="split"
assert ["0","1","2","4"] == list("0,1,2,4".split(","))
assert ['a', 'b', 'c'] == 'a,b,c'.split(',', -1)
assert ['a,b,c'] == 'a,b,c'.split(',', 0)
assert ['a', 'b', 'c'] == 'a,b,c'.split(',', 2147483647)
assertRaisesText(ValueError, "empty separator", lambda: 'abc'.split(''))
assertRaisesText(ValueError, "empty separator", lambda: ''.split('', 0))
assert [""] == list("".split(","))
assert ['a', 'd,c'] == list("a,d,c".split(",",1))
assert ['a', 'd', 'b'] == list(" a   d   b   ".split())
assert ['a', 'd   b   '] == list(" a   d   b   ".split(None, 1))
assertRaisesText(TypeError, "Can't convert 'int' object to str implicitly", lambda: "0,1,2,4".split(1))

doc="ascii len"
assert len(asc) == 5

doc="unicode len"
assert len(uni) == 7

doc="ascii index"
assert asc[0] == "h"
assert asc[1] == "e"
assert asc[2] == "l"
assert asc[3] == "l"
assert asc[4] == "o"
assert asc[-5] == "h"
assert asc[-4] == "e"
assert asc[-3] == "l"
assert asc[-2] == "l"
assert asc[-1] == "o"
def index(s, i):
    return s[i]
indexError = "index out of range"
assertRaisesText(IndexError, indexError, index, asc, -6)
assertRaisesText(IndexError, indexError, index, asc, 5)

doc="unicode index"
assert uni[0] == "£"
assert uni[1] == "1"
assert uni[2] == "0"
assert uni[3] == "0"
assert uni[4] == "世"
assert uni[5] == "界"
assert uni[6] == "𠜎"
assert uni[-7] == "£"
assert uni[-6] == "1"
assert uni[-5] == "0"
assert uni[-4] == "0"
assert uni[-3] == "世"
assert uni[-2] == "界"
assert uni[-1] == "𠜎"
assertRaisesText(IndexError, indexError, index, asc, -8)
assertRaisesText(IndexError, indexError, index, asc, 7)

doc="ascii slice"
assert asc[3:3] == ""
assert asc[:3] == "hel"
assert asc[1:4] == "ell"
assert asc[1:-1] == "ell"
assert asc[2:] == "llo"
assert asc[3:2] == ""
assert asc[-100:100] == "hello"
assert asc[100:200] == ""
assertRaisesText(ValueError, "slice step cannot be zero", lambda: asc[1:2:0])

doc="ascii slice exhaustive"
assert asc[0:0] == ''
assert asc[0:1] == 'h'
assert asc[0:2] == 'he'
assert asc[0:3] == 'hel'
assert asc[0:4] == 'hell'
assert asc[0:5] == 'hello'
assert asc[1:0] == ''
assert asc[1:1] == ''
assert asc[1:2] == 'e'
assert asc[1:3] == 'el'
assert asc[1:4] == 'ell'
assert asc[1:5] == 'ello'
assert asc[2:0] == ''
assert asc[2:1] == ''
assert asc[2:2] == ''
assert asc[2:3] == 'l'
assert asc[2:4] == 'll'
assert asc[2:5] == 'llo'
assert asc[3:0] == ''
assert asc[3:1] == ''
assert asc[3:2] == ''
assert asc[3:3] == ''
assert asc[3:4] == 'l'
assert asc[3:5] == 'lo'
assert asc[4:0] == ''
assert asc[4:1] == ''
assert asc[4:2] == ''
assert asc[4:3] == ''
assert asc[4:4] == ''
assert asc[4:5] == 'o'
assert asc[5:0] == ''
assert asc[5:1] == ''
assert asc[5:2] == ''
assert asc[5:3] == ''
assert asc[5:4] == ''
assert asc[5:5] == ''

doc="unicode slice"
assert uni[3:3] == ""
assert uni[:3] == "£10"
assert uni[1:4] == "100"
assert uni[1:-1] == "100世界"
assert uni[2:] == "00世界𠜎"
assert uni[3:2] == ""
assert uni[-100:100] == "£100世界𠜎"
assert uni[100:200] == ""

doc="unicode slice exhaustive"
assert uni[0:0] == ''
assert uni[0:1] == '£'
assert uni[0:2] == '£1'
assert uni[0:3] == '£10'
assert uni[0:4] == '£100'
assert uni[0:5] == '£100世'
assert uni[0:6] == '£100世界'
assert uni[0:7] == '£100世界𠜎'
assert uni[1:0] == ''
assert uni[1:1] == ''
assert uni[1:2] == '1'
assert uni[1:3] == '10'
assert uni[1:4] == '100'
assert uni[1:5] == '100世'
assert uni[1:6] == '100世界'
assert uni[1:7] == '100世界𠜎'
assert uni[2:0] == ''
assert uni[2:1] == ''
assert uni[2:2] == ''
assert uni[2:3] == '0'
assert uni[2:4] == '00'
assert uni[2:5] == '00世'
assert uni[2:6] == '00世界'
assert uni[2:7] == '00世界𠜎'
assert uni[3:0] == ''
assert uni[3:1] == ''
assert uni[3:2] == ''
assert uni[3:3] == ''
assert uni[3:4] == '0'
assert uni[3:5] == '0世'
assert uni[3:6] == '0世界'
assert uni[3:7] == '0世界𠜎'
assert uni[4:0] == ''
assert uni[4:1] == ''
assert uni[4:2] == ''
assert uni[4:3] == ''
assert uni[4:4] == ''
assert uni[4:5] == '世'
assert uni[4:6] == '世界'
assert uni[4:7] == '世界𠜎'
assert uni[5:0] == ''
assert uni[5:1] == ''
assert uni[5:2] == ''
assert uni[5:3] == ''
assert uni[5:4] == ''
assert uni[5:5] == ''
assert uni[5:6] == '界'
assert uni[5:7] == '界𠜎'
assert uni[6:0] == ''
assert uni[6:1] == ''
assert uni[6:2] == ''
assert uni[6:3] == ''
assert uni[6:4] == ''
assert uni[6:5] == ''
assert uni[6:6] == ''
assert uni[6:7] == '𠜎'

doc="ascii slice triple"
assert asc[::-1] == 'olleh'

doc="ascii slice triple exhaustive"
assert asc[0:0:-3] == ''
assert asc[0:0:-2] == ''
assert asc[0:0:-1] == ''
assert asc[0:0:1] == ''
assert asc[0:0:2] == ''
assert asc[0:0:3] == ''
assert asc[0:1:-3] == ''
assert asc[0:1:-2] == ''
assert asc[0:1:-1] == ''
assert asc[0:1:1] == 'h'
assert asc[0:1:2] == 'h'
assert asc[0:1:3] == 'h'
assert asc[0:2:-3] == ''
assert asc[0:2:-2] == ''
assert asc[0:2:-1] == ''
assert asc[0:2:1] == 'he'
assert asc[0:2:2] == 'h'
assert asc[0:2:3] == 'h'
assert asc[0:3:-3] == ''
assert asc[0:3:-2] == ''
assert asc[0:3:-1] == ''
assert asc[0:3:1] == 'hel'
assert asc[0:3:2] == 'hl'
assert asc[0:3:3] == 'h'
assert asc[0:4:-3] == ''
assert asc[0:4:-2] == ''
assert asc[0:4:-1] == ''
assert asc[0:4:1] == 'hell'
assert asc[0:4:2] == 'hl'
assert asc[0:4:3] == 'hl'
assert asc[0:5:-3] == ''
assert asc[0:5:-2] == ''
assert asc[0:5:-1] == ''
assert asc[0:5:1] == 'hello'
assert asc[0:5:2] == 'hlo'
assert asc[0:5:3] == 'hl'
assert asc[1:0:-3] == 'e'
assert asc[1:0:-2] == 'e'
assert asc[1:0:-1] == 'e'
assert asc[1:0:1] == ''
assert asc[1:0:2] == ''
assert asc[1:0:3] == ''
assert asc[1:1:-3] == ''
assert asc[1:1:-2] == ''
assert asc[1:1:-1] == ''
assert asc[1:1:1] == ''
assert asc[1:1:2] == ''
assert asc[1:1:3] == ''
assert asc[1:2:-3] == ''
assert asc[1:2:-2] == ''
assert asc[1:2:-1] == ''
assert asc[1:2:1] == 'e'
assert asc[1:2:2] == 'e'
assert asc[1:2:3] == 'e'
assert asc[1:3:-3] == ''
assert asc[1:3:-2] == ''
assert asc[1:3:-1] == ''
assert asc[1:3:1] == 'el'
assert asc[1:3:2] == 'e'
assert asc[1:3:3] == 'e'
assert asc[1:4:-3] == ''
assert asc[1:4:-2] == ''
assert asc[1:4:-1] == ''
assert asc[1:4:1] == 'ell'
assert asc[1:4:2] == 'el'
assert asc[1:4:3] == 'e'
assert asc[1:5:-3] == ''
assert asc[1:5:-2] == ''
assert asc[1:5:-1] == ''
assert asc[1:5:1] == 'ello'
assert asc[1:5:2] == 'el'
assert asc[1:5:3] == 'eo'
assert asc[2:0:-3] == 'l'
assert asc[2:0:-2] == 'l'
assert asc[2:0:-1] == 'le'
assert asc[2:0:1] == ''
assert asc[2:0:2] == ''
assert asc[2:0:3] == ''
assert asc[2:1:-3] == 'l'
assert asc[2:1:-2] == 'l'
assert asc[2:1:-1] == 'l'
assert asc[2:1:1] == ''
assert asc[2:1:2] == ''
assert asc[2:1:3] == ''
assert asc[2:2:-3] == ''
assert asc[2:2:-2] == ''
assert asc[2:2:-1] == ''
assert asc[2:2:1] == ''
assert asc[2:2:2] == ''
assert asc[2:2:3] == ''
assert asc[2:3:-3] == ''
assert asc[2:3:-2] == ''
assert asc[2:3:-1] == ''
assert asc[2:3:1] == 'l'
assert asc[2:3:2] == 'l'
assert asc[2:3:3] == 'l'
assert asc[2:4:-3] == ''
assert asc[2:4:-2] == ''
assert asc[2:4:-1] == ''
assert asc[2:4:1] == 'll'
assert asc[2:4:2] == 'l'
assert asc[2:4:3] == 'l'
assert asc[2:5:-3] == ''
assert asc[2:5:-2] == ''
assert asc[2:5:-1] == ''
assert asc[2:5:1] == 'llo'
assert asc[2:5:2] == 'lo'
assert asc[2:5:3] == 'l'
assert asc[3:0:-3] == 'l'
assert asc[3:0:-2] == 'le'
assert asc[3:0:-1] == 'lle'
assert asc[3:0:1] == ''
assert asc[3:0:2] == ''
assert asc[3:0:3] == ''
assert asc[3:1:-3] == 'l'
assert asc[3:1:-2] == 'l'
assert asc[3:1:-1] == 'll'
assert asc[3:1:1] == ''
assert asc[3:1:2] == ''
assert asc[3:1:3] == ''
assert asc[3:2:-3] == 'l'
assert asc[3:2:-2] == 'l'
assert asc[3:2:-1] == 'l'
assert asc[3:2:1] == ''
assert asc[3:2:2] == ''
assert asc[3:2:3] == ''
assert asc[3:3:-3] == ''
assert asc[3:3:-2] == ''
assert asc[3:3:-1] == ''
assert asc[3:3:1] == ''
assert asc[3:3:2] == ''
assert asc[3:3:3] == ''
assert asc[3:4:-3] == ''
assert asc[3:4:-2] == ''
assert asc[3:4:-1] == ''
assert asc[3:4:1] == 'l'
assert asc[3:4:2] == 'l'
assert asc[3:4:3] == 'l'
assert asc[3:5:-3] == ''
assert asc[3:5:-2] == ''
assert asc[3:5:-1] == ''
assert asc[3:5:1] == 'lo'
assert asc[3:5:2] == 'l'
assert asc[3:5:3] == 'l'
assert asc[4:0:-3] == 'oe'
assert asc[4:0:-2] == 'ol'
assert asc[4:0:-1] == 'olle'
assert asc[4:0:1] == ''
assert asc[4:0:2] == ''
assert asc[4:0:3] == ''
assert asc[4:1:-3] == 'o'
assert asc[4:1:-2] == 'ol'
assert asc[4:1:-1] == 'oll'
assert asc[4:1:1] == ''
assert asc[4:1:2] == ''
assert asc[4:1:3] == ''
assert asc[4:2:-3] == 'o'
assert asc[4:2:-2] == 'o'
assert asc[4:2:-1] == 'ol'
assert asc[4:2:1] == ''
assert asc[4:2:2] == ''
assert asc[4:2:3] == ''
assert asc[4:3:-3] == 'o'
assert asc[4:3:-2] == 'o'
assert asc[4:3:-1] == 'o'
assert asc[4:3:1] == ''
assert asc[4:3:2] == ''
assert asc[4:3:3] == ''
assert asc[4:4:-3] == ''
assert asc[4:4:-2] == ''
assert asc[4:4:-1] == ''
assert asc[4:4:1] == ''
assert asc[4:4:2] == ''
assert asc[4:4:3] == ''
assert asc[4:5:-3] == ''
assert asc[4:5:-2] == ''
assert asc[4:5:-1] == ''
assert asc[4:5:1] == 'o'
assert asc[4:5:2] == 'o'
assert asc[4:5:3] == 'o'
assert asc[5:0:-3] == 'oe'
assert asc[5:0:-2] == 'ol'
assert asc[5:0:-1] == 'olle'
assert asc[5:0:1] == ''
assert asc[5:0:2] == ''
assert asc[5:0:3] == ''
assert asc[5:1:-3] == 'o'
assert asc[5:1:-2] == 'ol'
assert asc[5:1:-1] == 'oll'
assert asc[5:1:1] == ''
assert asc[5:1:2] == ''
assert asc[5:1:3] == ''
assert asc[5:2:-3] == 'o'
assert asc[5:2:-2] == 'o'
assert asc[5:2:-1] == 'ol'
assert asc[5:2:1] == ''
assert asc[5:2:2] == ''
assert asc[5:2:3] == ''
assert asc[5:3:-3] == 'o'
assert asc[5:3:-2] == 'o'
assert asc[5:3:-1] == 'o'
assert asc[5:3:1] == ''
assert asc[5:3:2] == ''
assert asc[5:3:3] == ''
assert asc[5:4:-3] == ''
assert asc[5:4:-2] == ''
assert asc[5:4:-1] == ''
assert asc[5:4:1] == ''
assert asc[5:4:2] == ''
assert asc[5:4:3] == ''
assert asc[5:5:-3] == ''
assert asc[5:5:-2] == ''
assert asc[5:5:-1] == ''
assert asc[5:5:1] == ''
assert asc[5:5:2] == ''
assert asc[5:5:3] == ''

doc="unicode triple"
assert uni[::-1] == "𠜎界世001£"

doc="unicode triple exhaustive"
assert uni[0:0:-3] == ''
assert uni[0:0:-2] == ''
assert uni[0:0:-1] == ''
assert uni[0:0:1] == ''
assert uni[0:0:2] == ''
assert uni[0:0:3] == ''
assert uni[0:1:-3] == ''
assert uni[0:1:-2] == ''
assert uni[0:1:-1] == ''
assert uni[0:1:1] == '£'
assert uni[0:1:2] == '£'
assert uni[0:1:3] == '£'
assert uni[0:2:-3] == ''
assert uni[0:2:-2] == ''
assert uni[0:2:-1] == ''
assert uni[0:2:1] == '£1'
assert uni[0:2:2] == '£'
assert uni[0:2:3] == '£'
assert uni[0:3:-3] == ''
assert uni[0:3:-2] == ''
assert uni[0:3:-1] == ''
assert uni[0:3:1] == '£10'
assert uni[0:3:2] == '£0'
assert uni[0:3:3] == '£'
assert uni[0:4:-3] == ''
assert uni[0:4:-2] == ''
assert uni[0:4:-1] == ''
assert uni[0:4:1] == '£100'
assert uni[0:4:2] == '£0'
assert uni[0:4:3] == '£0'
assert uni[0:5:-3] == ''
assert uni[0:5:-2] == ''
assert uni[0:5:-1] == ''
assert uni[0:5:1] == '£100世'
assert uni[0:5:2] == '£0世'
assert uni[0:5:3] == '£0'
assert uni[0:6:-3] == ''
assert uni[0:6:-2] == ''
assert uni[0:6:-1] == ''
assert uni[0:6:1] == '£100世界'
assert uni[0:6:2] == '£0世'
assert uni[0:6:3] == '£0'
assert uni[0:7:-3] == ''
assert uni[0:7:-2] == ''
assert uni[0:7:-1] == ''
assert uni[0:7:1] == '£100世界𠜎'
assert uni[0:7:2] == '£0世𠜎'
assert uni[0:7:3] == '£0𠜎'
assert uni[1:0:-3] == '1'
assert uni[1:0:-2] == '1'
assert uni[1:0:-1] == '1'
assert uni[1:0:1] == ''
assert uni[1:0:2] == ''
assert uni[1:0:3] == ''
assert uni[1:1:-3] == ''
assert uni[1:1:-2] == ''
assert uni[1:1:-1] == ''
assert uni[1:1:1] == ''
assert uni[1:1:2] == ''
assert uni[1:1:3] == ''
assert uni[1:2:-3] == ''
assert uni[1:2:-2] == ''
assert uni[1:2:-1] == ''
assert uni[1:2:1] == '1'
assert uni[1:2:2] == '1'
assert uni[1:2:3] == '1'
assert uni[1:3:-3] == ''
assert uni[1:3:-2] == ''
assert uni[1:3:-1] == ''
assert uni[1:3:1] == '10'
assert uni[1:3:2] == '1'
assert uni[1:3:3] == '1'
assert uni[1:4:-3] == ''
assert uni[1:4:-2] == ''
assert uni[1:4:-1] == ''
assert uni[1:4:1] == '100'
assert uni[1:4:2] == '10'
assert uni[1:4:3] == '1'
assert uni[1:5:-3] == ''
assert uni[1:5:-2] == ''
assert uni[1:5:-1] == ''
assert uni[1:5:1] == '100世'
assert uni[1:5:2] == '10'
assert uni[1:5:3] == '1世'
assert uni[1:6:-3] == ''
assert uni[1:6:-2] == ''
assert uni[1:6:-1] == ''
assert uni[1:6:1] == '100世界'
assert uni[1:6:2] == '10界'
assert uni[1:6:3] == '1世'
assert uni[1:7:-3] == ''
assert uni[1:7:-2] == ''
assert uni[1:7:-1] == ''
assert uni[1:7:1] == '100世界𠜎'
assert uni[1:7:2] == '10界'
assert uni[1:7:3] == '1世'
assert uni[2:0:-3] == '0'
assert uni[2:0:-2] == '0'
assert uni[2:0:-1] == '01'
assert uni[2:0:1] == ''
assert uni[2:0:2] == ''
assert uni[2:0:3] == ''
assert uni[2:1:-3] == '0'
assert uni[2:1:-2] == '0'
assert uni[2:1:-1] == '0'
assert uni[2:1:1] == ''
assert uni[2:1:2] == ''
assert uni[2:1:3] == ''
assert uni[2:2:-3] == ''
assert uni[2:2:-2] == ''
assert uni[2:2:-1] == ''
assert uni[2:2:1] == ''
assert uni[2:2:2] == ''
assert uni[2:2:3] == ''
assert uni[2:3:-3] == ''
assert uni[2:3:-2] == ''
assert uni[2:3:-1] == ''
assert uni[2:3:1] == '0'
assert uni[2:3:2] == '0'
assert uni[2:3:3] == '0'
assert uni[2:4:-3] == ''
assert uni[2:4:-2] == ''
assert uni[2:4:-1] == ''
assert uni[2:4:1] == '00'
assert uni[2:4:2] == '0'
assert uni[2:4:3] == '0'
assert uni[2:5:-3] == ''
assert uni[2:5:-2] == ''
assert uni[2:5:-1] == ''
assert uni[2:5:1] == '00世'
assert uni[2:5:2] == '0世'
assert uni[2:5:3] == '0'
assert uni[2:6:-3] == ''
assert uni[2:6:-2] == ''
assert uni[2:6:-1] == ''
assert uni[2:6:1] == '00世界'
assert uni[2:6:2] == '0世'
assert uni[2:6:3] == '0界'
assert uni[2:7:-3] == ''
assert uni[2:7:-2] == ''
assert uni[2:7:-1] == ''
assert uni[2:7:1] == '00世界𠜎'
assert uni[2:7:2] == '0世𠜎'
assert uni[2:7:3] == '0界'
assert uni[3:0:-3] == '0'
assert uni[3:0:-2] == '01'
assert uni[3:0:-1] == '001'
assert uni[3:0:1] == ''
assert uni[3:0:2] == ''
assert uni[3:0:3] == ''
assert uni[3:1:-3] == '0'
assert uni[3:1:-2] == '0'
assert uni[3:1:-1] == '00'
assert uni[3:1:1] == ''
assert uni[3:1:2] == ''
assert uni[3:1:3] == ''
assert uni[3:2:-3] == '0'
assert uni[3:2:-2] == '0'
assert uni[3:2:-1] == '0'
assert uni[3:2:1] == ''
assert uni[3:2:2] == ''
assert uni[3:2:3] == ''
assert uni[3:3:-3] == ''
assert uni[3:3:-2] == ''
assert uni[3:3:-1] == ''
assert uni[3:3:1] == ''
assert uni[3:3:2] == ''
assert uni[3:3:3] == ''
assert uni[3:4:-3] == ''
assert uni[3:4:-2] == ''
assert uni[3:4:-1] == ''
assert uni[3:4:1] == '0'
assert uni[3:4:2] == '0'
assert uni[3:4:3] == '0'
assert uni[3:5:-3] == ''
assert uni[3:5:-2] == ''
assert uni[3:5:-1] == ''
assert uni[3:5:1] == '0世'
assert uni[3:5:2] == '0'
assert uni[3:5:3] == '0'
assert uni[3:6:-3] == ''
assert uni[3:6:-2] == ''
assert uni[3:6:-1] == ''
assert uni[3:6:1] == '0世界'
assert uni[3:6:2] == '0界'
assert uni[3:6:3] == '0'
assert uni[3:7:-3] == ''
assert uni[3:7:-2] == ''
assert uni[3:7:-1] == ''
assert uni[3:7:1] == '0世界𠜎'
assert uni[3:7:2] == '0界'
assert uni[3:7:3] == '0𠜎'
assert uni[4:0:-3] == '世1'
assert uni[4:0:-2] == '世0'
assert uni[4:0:-1] == '世001'
assert uni[4:0:1] == ''
assert uni[4:0:2] == ''
assert uni[4:0:3] == ''
assert uni[4:1:-3] == '世'
assert uni[4:1:-2] == '世0'
assert uni[4:1:-1] == '世00'
assert uni[4:1:1] == ''
assert uni[4:1:2] == ''
assert uni[4:1:3] == ''
assert uni[4:2:-3] == '世'
assert uni[4:2:-2] == '世'
assert uni[4:2:-1] == '世0'
assert uni[4:2:1] == ''
assert uni[4:2:2] == ''
assert uni[4:2:3] == ''
assert uni[4:3:-3] == '世'
assert uni[4:3:-2] == '世'
assert uni[4:3:-1] == '世'
assert uni[4:3:1] == ''
assert uni[4:3:2] == ''
assert uni[4:3:3] == ''
assert uni[4:4:-3] == ''
assert uni[4:4:-2] == ''
assert uni[4:4:-1] == ''
assert uni[4:4:1] == ''
assert uni[4:4:2] == ''
assert uni[4:4:3] == ''
assert uni[4:5:-3] == ''
assert uni[4:5:-2] == ''
assert uni[4:5:-1] == ''
assert uni[4:5:1] == '世'
assert uni[4:5:2] == '世'
assert uni[4:5:3] == '世'
assert uni[4:6:-3] == ''
assert uni[4:6:-2] == ''
assert uni[4:6:-1] == ''
assert uni[4:6:1] == '世界'
assert uni[4:6:2] == '世'
assert uni[4:6:3] == '世'
assert uni[4:7:-3] == ''
assert uni[4:7:-2] == ''
assert uni[4:7:-1] == ''
assert uni[4:7:1] == '世界𠜎'
assert uni[4:7:2] == '世𠜎'
assert uni[4:7:3] == '世'
assert uni[5:0:-3] == '界0'
assert uni[5:0:-2] == '界01'
assert uni[5:0:-1] == '界世001'
assert uni[5:0:1] == ''
assert uni[5:0:2] == ''
assert uni[5:0:3] == ''
assert uni[5:1:-3] == '界0'
assert uni[5:1:-2] == '界0'
assert uni[5:1:-1] == '界世00'
assert uni[5:1:1] == ''
assert uni[5:1:2] == ''
assert uni[5:1:3] == ''
assert uni[5:2:-3] == '界'
assert uni[5:2:-2] == '界0'
assert uni[5:2:-1] == '界世0'
assert uni[5:2:1] == ''
assert uni[5:2:2] == ''
assert uni[5:2:3] == ''
assert uni[5:3:-3] == '界'
assert uni[5:3:-2] == '界'
assert uni[5:3:-1] == '界世'
assert uni[5:3:1] == ''
assert uni[5:3:2] == ''
assert uni[5:3:3] == ''
assert uni[5:4:-3] == '界'
assert uni[5:4:-2] == '界'
assert uni[5:4:-1] == '界'
assert uni[5:4:1] == ''
assert uni[5:4:2] == ''
assert uni[5:4:3] == ''
assert uni[5:5:-3] == ''
assert uni[5:5:-2] == ''
assert uni[5:5:-1] == ''
assert uni[5:5:1] == ''
assert uni[5:5:2] == ''
assert uni[5:5:3] == ''
assert uni[5:6:-3] == ''
assert uni[5:6:-2] == ''
assert uni[5:6:-1] == ''
assert uni[5:6:1] == '界'
assert uni[5:6:2] == '界'
assert uni[5:6:3] == '界'
assert uni[5:7:-3] == ''
assert uni[5:7:-2] == ''
assert uni[5:7:-1] == ''
assert uni[5:7:1] == '界𠜎'
assert uni[5:7:2] == '界'
assert uni[5:7:3] == '界'
assert uni[6:0:-3] == '𠜎0'
assert uni[6:0:-2] == '𠜎世0'
assert uni[6:0:-1] == '𠜎界世001'
assert uni[6:0:1] == ''
assert uni[6:0:2] == ''
assert uni[6:0:3] == ''
assert uni[6:1:-3] == '𠜎0'
assert uni[6:1:-2] == '𠜎世0'
assert uni[6:1:-1] == '𠜎界世00'
assert uni[6:1:1] == ''
assert uni[6:1:2] == ''
assert uni[6:1:3] == ''
assert uni[6:2:-3] == '𠜎0'
assert uni[6:2:-2] == '𠜎世'
assert uni[6:2:-1] == '𠜎界世0'
assert uni[6:2:1] == ''
assert uni[6:2:2] == ''
assert uni[6:2:3] == ''
assert uni[6:3:-3] == '𠜎'
assert uni[6:3:-2] == '𠜎世'
assert uni[6:3:-1] == '𠜎界世'
assert uni[6:3:1] == ''
assert uni[6:3:2] == ''
assert uni[6:3:3] == ''
assert uni[6:4:-3] == '𠜎'
assert uni[6:4:-2] == '𠜎'
assert uni[6:4:-1] == '𠜎界'
assert uni[6:4:1] == ''
assert uni[6:4:2] == ''
assert uni[6:4:3] == ''
assert uni[6:5:-3] == '𠜎'
assert uni[6:5:-2] == '𠜎'
assert uni[6:5:-1] == '𠜎'
assert uni[6:5:1] == ''
assert uni[6:5:2] == ''
assert uni[6:5:3] == ''
assert uni[6:6:-3] == ''
assert uni[6:6:-2] == ''
assert uni[6:6:-1] == ''
assert uni[6:6:1] == ''
assert uni[6:6:2] == ''
assert uni[6:6:3] == ''
assert uni[6:7:-3] == ''
assert uni[6:7:-2] == ''
assert uni[6:7:-1] == ''
assert uni[6:7:1] == '𠜎'
assert uni[6:7:2] == '𠜎'
assert uni[6:7:3] == '𠜎'
assert uni[7:0:-3] == '𠜎0'
assert uni[7:0:-2] == '𠜎世0'
assert uni[7:0:-1] == '𠜎界世001'
assert uni[7:0:1] == ''
assert uni[7:0:2] == ''
assert uni[7:0:3] == ''
assert uni[7:1:-3] == '𠜎0'
assert uni[7:1:-2] == '𠜎世0'
assert uni[7:1:-1] == '𠜎界世00'
assert uni[7:1:1] == ''
assert uni[7:1:2] == ''
assert uni[7:1:3] == ''
assert uni[7:2:-3] == '𠜎0'
assert uni[7:2:-2] == '𠜎世'
assert uni[7:2:-1] == '𠜎界世0'
assert uni[7:2:1] == ''
assert uni[7:2:2] == ''
assert uni[7:2:3] == ''
assert uni[7:3:-3] == '𠜎'
assert uni[7:3:-2] == '𠜎世'
assert uni[7:3:-1] == '𠜎界世'
assert uni[7:3:1] == ''
assert uni[7:3:2] == ''
assert uni[7:3:3] == ''
assert uni[7:4:-3] == '𠜎'
assert uni[7:4:-2] == '𠜎'
assert uni[7:4:-1] == '𠜎界'
assert uni[7:4:1] == ''
assert uni[7:4:2] == ''
assert uni[7:4:3] == ''
assert uni[7:5:-3] == '𠜎'
assert uni[7:5:-2] == '𠜎'
assert uni[7:5:-1] == '𠜎'
assert uni[7:5:1] == ''
assert uni[7:5:2] == ''
assert uni[7:5:3] == ''
assert uni[7:6:-3] == ''
assert uni[7:6:-2] == ''
assert uni[7:6:-1] == ''
assert uni[7:6:1] == ''
assert uni[7:6:2] == ''
assert uni[7:6:3] == ''
assert uni[7:7:-3] == ''
assert uni[7:7:-2] == ''
assert uni[7:7:-1] == ''
assert uni[7:7:1] == ''
assert uni[7:7:2] == ''
assert uni[7:7:3] == ''

doc="string strip methods"
a = "  adfasd   "
assert a.rstrip() == "  adfasd"
assert a.lstrip() == "adfasd   "
assert a.strip() == "adfasd"

a = " a bada a"
assert a.rstrip("a ") == " a bad"
assert a.lstrip("a ") == "bada a"
assert a.strip("a ") == "bad"

doc="upper"
a = "abc"
assert a.upper() == "ABC"

doc="lower"
a = "ABC"
assert a.lower() == "abc"

doc="join"
assert ",".join(['a', 'b', 'c']) == "a,b,c"
assert " ".join(('a', 'b', 'c')) == "a b c"
assert " ".join("abc") == "a b c"
assert "".join(['a', 'b', 'c']) == "abc"
assert ",".join([]) == ""
assert ",".join(()) == ""
assertRaises(TypeError, lambda: ",".join([1, 2, 3]))

class Index:
    def __index__(self):
        return 1

a = '012345678910'
b = Index()
assert a[b] == '1'
assert a[b:10] == a[1:10]
assert a[10:b:-1] == a[10:1:-1]

class NonIntegerIndex:
    def __index__(self):
        return 1.1

a = '012345678910'
b = NonIntegerIndex()
try:
    a[b]
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

try:
    a[b:10]
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

doc="count"
assert 'hello world'.count('l') == 3
assert 'hello world'.count('l', 3) == 2
assert 'hello world'.count('l', 3, 10) == 2
assert 'hello world'.count('l', 3, 100) == 2
assert 'hello world'.count('l', 3, 5) == 1
assert 'hello world'.count('l', 3, 1) == 0
assert 'hello world'.count('z') == 0


doc="finished"
