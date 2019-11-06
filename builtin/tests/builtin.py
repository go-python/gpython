# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaises

doc="abs"
assert abs(0) == 0
assert abs(10) == 10
assert abs(-10) == 10

doc="all"
assert all((0,0,0)) == False
assert all((1,1,0)) == False
assert all(["hello", "world"]) == True
assert all([]) == True

doc="any"
assert any((0,0,0)) == False
assert any((1,1,0)) == True
assert any(["hello", "world"]) == True
assert any([]) == False

doc="ascii"
assert ascii('hello world') == "'hello world'"
assert ascii('안녕 세상') == "'\\uc548\\ub155 \\uc138\\uc0c1'"
assert ascii(chr(0x10001)) == "'\\U00010001'"
assert ascii('안녕 gpython') == "'\\uc548\\ub155 gpython'"

doc="bin"
assert bin(False) == '0b0'
assert bin(True) == '0b1'
assert bin(0) == '0b0'
assert bin(1) == '0b1'
assert bin(-1) == '-0b1'
assert bin(10) == '0b1010'
assert bin(-10) == '-0b1010'
assert bin(2**32) == '0b100000000000000000000000000000000'
assert bin(2**32-1) == '0b11111111111111111111111111111111'
assert bin(-(2**32)) == '-0b100000000000000000000000000000000'
assert bin(-(2**32-1)) == '-0b11111111111111111111111111111111'

doc="chr"
assert chr(65) == "A"
assert chr(163) == "£"
assert chr(0x263A) == "☺"

doc="compile"
code = compile("pass", "<string>", "exec")
assert code is not None
# FIXME

doc="divmod"
assert divmod(34,7) == (4, 6)

doc="enumerate"
a = [3, 4, 5, 6, 7]
for idx, value in enumerate(a):
    assert value == a[idx]

doc="eval"
# smoke test only - see vm/tests/builtin.py for more tests
assert eval("1+2") == 3

doc="exec"
# smoke test only - see vm/tests/builtin.py for more tests
glob = {"a":100}
assert exec("b = a+100", glob) == None
assert glob["b"] == 200

doc="getattr"
class C:
    def __init__(self):
        self.potato = 42
c = C()
assert getattr(c, "potato") == 42
assert getattr(c, "potato", 43) == 42
assert getattr(c, "sausage", 43) == 43

doc="globals"
a = 1
assert globals()["a"] == 1

doc="hasattr"
assert hasattr(c, "potato")
assert not hasattr(c, "sausage")

doc="len"
assert len(()) == 0
assert len((1,2,3)) == 3
assert len("hello") == 5
assert len("£☺") == 2

doc="locals"
def fn(x):
    assert locals()["x"] == 1
fn(1)

def func(p):
   return p[1]

doc="min"
values = (1,2,3)
v = min(values)
assert v == 1
v = min(4,5,6)
assert v == 4
v = min((), default=-1)
assert v == -1
v = min([], default=-1)
assert v == -1
v = min([], key=func, default=(1,3))
assert v == (1,3)
v = min([(1,3), (2,1)], key=func)
assert v == (2,1)
ok = False
try:
    min([(1,3), (2,1)], key=3)
except TypeError:
    ok = True
assert ok, "TypeError not raised"
ok = False
try:
    min([])
except ValueError:
    ok = True
assert ok, "ValueError not raised"

doc="max"
values = (1,2,3)
v = max(values)
assert v == 3
v = max(4,5,6)
assert v == 6
v = max((), default=-1)
assert v == -1
v = max([], default=-1)
assert v == -1
v = max([], key=func, default=(1,3))
assert v == (1,3)
v = max([(1,3), (2,1)], key=func)
assert v == (1,3)
ok = False
try:
    max([(1,3), (2,1)], key=3)
except TypeError:
    ok = True
assert ok, "TypeError not raised"
ok = False
try:
    max([])
except ValueError:
    ok = True
assert ok, "ValueError not raised"

doc="hex"
assert hex( 0)=="0x0",    "hex(0)"
assert hex( 1)=="0x1",    "hex(1)"
assert hex(42)=="0x2a",   "hex(42)"
assert hex( -0)=="0x0",   "hex(-0)"
assert hex( -1)=="-0x1",  "hex(-1)"
assert hex(-42)=="-0x2a", "hex(-42)"
assert hex( 1<<64) ==  "0x10000000000000000", "hex(1<<64)"
assert hex(-1<<64) == "-0x10000000000000000", "hex(-1<<64)"
assert hex( 1<<128) ==  "0x100000000000000000000000000000000", "hex(1<<128)"
assert hex(-1<<128) == "-0x100000000000000000000000000000000", "hex(-1<<128)"
assertRaises(TypeError, hex, 10.0) ## TypeError: 'float' object cannot be interpreted as an integer
assertRaises(TypeError, hex, float(0)) ## TypeError: 'float' object cannot be interpreted as an integer

doc="isinstance"
class A:
    pass
a = A()
assert isinstance(1, (str, tuple, int))
assert isinstance(a, (str, (tuple, (A, ))))
assertRaises(TypeError, isinstance, 1, (A, ), "foo")
assertRaises(TypeError, isinstance, 1, [A, "foo"])

doc="iter"
cnt = 0
def f():
    global cnt
    cnt += 1
    return cnt

l = list(iter(f,20))
assert len(l) == 19
for idx, v in enumerate(l):
    assert idx + 1 == v

words1 = ['g', 'p', 'y', 't', 'h', 'o', 'n']
words2 = list(iter(words1))
for w1, w2 in zip(words1, words2):
    assert w1 == w2

ok = False
try:
    iter()
except TypeError:
    ok = True
finally:
    assert ok, "TypeError not raised"
    ok = False

try:
    l = [1, 2, 3]
    iter(l, 2)
except TypeError:
    ok = True
finally:
    assert ok, "TypeError not raised"
    ok = False

try:
    iter(f, 2, 3)
except TypeError:
    ok = True
finally:
    assert ok, "TypeError not raised"
    ok = False


doc="next no default"
def gen():
    yield 1
    yield 2
g = gen()
assert next(g) == 1
assert next(g) == 2
ok = False
try:
    next(g)
except StopIteration:
    ok = True
assert ok, "StopIteration not raised"

doc="next with default"
g = gen()
assert next(g, 42) == 1
assert next(g, 42) == 2
assert next(g, 42) == 42
assert next(g, 42) == 42

doc="next no default with exception"
def gen2():
    yield 1
    raise ValueError("potato")
g = gen2()
assert next(g) == 1
ok = False
try:
    next(g)
except ValueError:
    ok = True
assert ok, "ValueError not raised"

doc="next with default and exception"
g = gen2()
assert next(g, 42) == 1
ok = False
try:
    next(g)
except ValueError:
    ok = True
assert ok, "ValueError not raised"

doc="ord"
assert 65 == ord("A")
assert 163 == ord("£")
assert 0x263A == ord("☺")
assert 65 == ord(b"A")
ok = False
try:
    ord("AA")
except TypeError as e:
    if e.args[0] != "ord() expected a character, but string of length 2 found":
        raise
    ok = True
assert ok, "TypeError not raised"
try:
    ord(None)
except TypeError as e:
    if e.args[0] != "ord() expected string of length 1, but NoneType found":
        raise
    ok = True
assert ok, "TypeError not raised"

doc="open"
assert open(__file__) is not None

doc="pow"
assert pow(2, 10) == 1024
assert pow(2, 10, 17) == 4

doc="repr"
assert repr(5) == "5"
assert repr("hello") == "'hello'"

doc="print"
ok = False
try:
    print("hello", sep=1)
except TypeError as e:
    #if e.args[0] != "sep must be None or a string, not int":
    #   raise
    ok = True
assert ok, "TypeError not raised"

try:
    print("hello", sep=" ", end=1)
except TypeError as e:
    #if e.args[0] != "end must be None or a string, not int":
    #   raise
    ok = True
assert ok, "TypeError not raised"

try:
    print("hello", sep=" ", end="\n", file=1)
except AttributeError as e:
    #if e.args[0] != "'int' object has no attribute 'write'":
    #   raise
    ok = True
assert ok, "AttributeError not raised"

with open("testfile", "w") as f:
    print("hello", "world", sep=" ", end="\n", file=f)

with open("testfile", "r") as f:
    assert f.read() == "hello world\n"

with open("testfile", "w") as f:
    print(1,2,3,sep=",",end=",\n", file=f)

with open("testfile", "r") as f:
    assert f.read() == "1,2,3,\n"

doc="round"
assert round(1.1) == 1.0

doc="setattr"
class C: pass
c = C()
assert not hasattr(c, "potato")
setattr(c, "potato", "spud")
assert getattr(c, "potato") == "spud"
assert c.potato == "spud"
delattr(c, "potato")
assert not hasattr(c, "potato")
ok = False
try:
    delattr(c, "potato")
except AttributeError as e:
    ok = True
finally:
    assert ok

doc="sorted"
a = [3, 1.1, 1, 2]
assert sorted(a) == [1, 1.1, 2, 3]
assert sorted(sorted(a)) == [1, 1.1, 2, 3]
assert sorted(a, reverse=True) == [3, 2, 1.1, 1]
assert sorted(a, key=lambda l: l+1) == [1, 1.1, 2, 3]
s = [2.0, 2, 1, 1.0]
assert sorted(s, key=lambda l: 0) == [2.0, 2, 1, 1.0]
assert [type(t) for t in sorted(s, key=lambda l: 0)] == [float, int, int, float]
assert sorted(s) == [1, 1.0, 2.0, 2]
assert [type(t) for t in sorted(s)] == [int, float, float, int]

try:
    sorted([2.0, "abc"])
except TypeError:
    pass
else:
    assert False

assert sorted([]) == []
assert sorted([0]) == [0]
s = [0, 1]
try:
    # Sorting a list of len >= 2 with uncallable key must fail on all Python implementations.
    sorted(s, key=1)
except TypeError:
    pass
else:
    assert False

try:
    sorted(1)
except TypeError:
    pass
else:
    assert False

try:
    sorted()
except TypeError:
    pass
else:
    assert False

doc="sum"
assert sum([1,2,3]) == 6
assert sum([1,2,3], 3) == 9
assert sum((1,2,3)) == 6
assert sum((1,2,3), 3) == 9
assert sum((1, 2.5, 3)) == 6.5
assert sum((1, 2.5, 3), 3) == 9.5

try:
    sum([1,2,3], 'hi')
except TypeError as e:
    if e.args[0] != "sum() can't sum strings [use ''.join(seq) instead]":
        raise
    ok = True
assert ok, "TypeError not raised"

try:
    sum([1,2,3], b'hi')
except TypeError as e:
    if e.args[0] != "sum() can't sum bytes [use b''.join(seq) instead]":
        raise
    ok = True
assert ok, "TypeError not raised"

try:
    sum(['h', 'i'])
except TypeError as e:
    if e.args[0] != "unsupported operand type(s) for +: 'int' and 'str'":
        raise
    ok = True
assert ok, "TypeError not raised"

doc="zip"
ok = False
a = [3, 4, 5, 6, 7]
b = [8, 9, 10, 11, 12]
assert [e for e in zip(a, b)] == [(3,8), (4,9), (5,10), (6,11), (7,12)]
try:
    zip(1,2,3)
except TypeError as e:
    if e.args[0] != "zip argument #1 must support iteration":
        raise
    ok = True
assert ok, "TypeError not raised"

doc="__import__"
lib = __import__("lib")
assert lib.libfn() == 42
assert lib.libvar == 43
assert lib.libclass().method() == 44

doc="finished"
