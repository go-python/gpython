#!/usr/bin/env python3.4

# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""
Write grammar_data_test.go
"""

import sys
import ast
import datetime
import subprocess

inp = [
    # basics
    ("", "exec"),
    ("\n", "exec"),
    ("()", "eval"),
    ("()", "exec"),
    ("[ ]", "exec"),
    ("True\n", "eval"),
    ("False\n", "eval"),
    ("None\n", "eval"),
    ("...", "eval"),
    ("abc123", "eval"),
    ('"abc"', "eval"),
    ('"abc" """123"""', "eval"),
    ("b'abc'", "eval"),
    ("b'abc' b'''123'''", "eval"),
    ("1234", "eval"),
    ("01234", "eval", SyntaxError, "illegal decimal with leading zero"),
    ("1234d", "eval", SyntaxError, "invalid syntax"),
    ("1234d", "exec", SyntaxError),
    ("1234d", "single", SyntaxError),
    ("0x1234", "eval"),
    ("12.34", "eval"),
    ("1,", "eval"),
    ("1,2", "eval"),
    ("1,2,", "eval"),
    ("{ }", "eval"),
    ("{1}", "eval"),
    ("{1,}", "eval"),
    ("{1,2}", "eval"),
    ("{1,2,3,}", "eval"),
    ("{ 'a':1 }", "eval"),
    ("{ 'a':1, 'b':2 }", "eval"),
    ("{ 'a':{'aa':11, 'bb':{'aa':11, 'bb':22}}, 'b':{'aa':11, 'bb':22} }", "eval"),
    ("(1)", "eval"),
    ("(1,)", "eval"),
    ("(1,2)", "eval"),
    ("(1,2,)", "eval"),
    ("{(1,2)}", "eval"),
    ("(((((1,),(2,),),(2,),),((1,),(2,),),),((1,),(2,),))", "eval"),
    ("(((1)))", "eval"),
    ("[1]", "eval"),
    ("[1,]", "eval"),
    ("[1,2]", "eval"),
    ("[1,2,]", "eval"),
    ("[e for e in (1,2,3)]", "eval"),

    # tuple
    ("( a for a in ab )", "eval"),
    ("( a for a, in ab )", "eval"),
    ("( a for a, b in ab )", "eval"),
    ("( a for a in ab if a )", "eval"),
    ("( a for a in ab if a if b if c )", "eval"),
    ("( a for a in ab for A in AB )", "eval"),
    ("( a for a in ab if a if b for A in AB if c )", "eval"),
    ("( a for a in ab if lambda: None )", "eval"),
    ("( a for a in ab if lambda x,y: x+y )", "eval"),

    # list
    ("[ a for a in ab ]", "eval"),
    ("[ a for a, in ab ]", "eval"),
    ("[ a for a, b in ab ]", "eval"),
    ("[ a for a in ab if a ]", "eval"),
    ("[ a for a in ab if a if b if c ]", "eval"),
    ("[ a for a in ab for A in AB ]", "eval"),
    ("[ a for a in ab if a if b for A in AB if c ]", "eval"),

    # set
    ("{ a for a in ab }", "eval"),
    ("{ a for a, in ab }", "eval"),
    ("{ a for a, b in ab }", "eval"),
    ("{ a for a in ab if a }", "eval"),
    ("{ a for a in ab if a if b if c }", "eval"),
    ("{ a for a in ab for A in AB }", "eval"),
    ("{ a for a in ab if a if b for A in AB if c }", "eval"),

    # dict
    ("{ a:b for a in ab }", "eval"),
    ("{ a:b for a, in ab }", "eval"),
    ("{ a:b for a, b in ab }", "eval"),
    ("{ a:b for a in ab if a }", "eval"),
    ("{ a:b for a in ab if a if b if c }", "eval"),
    ("{ a:b for a in ab for A in AB }", "eval"),
    ("{ a:b for a in ab if a if b for A in AB if c }", "eval"),

    # BinOp
    ("a|b", "eval"),
    ("a^b", "eval"),
    ("a&b", "eval"),
    ("a<<b", "eval"),
    ("a>>b", "eval"),
    ("a+b", "eval"),
    ("a-b", "eval"),
    ("a*b", "eval"),
    ("a/b", "eval"),
    ("a//b", "eval"),
    ("a**b", "eval"),

    # UnaryOp
    ("not a", "eval"),
    ("+a", "eval"),
    ("-a", "eval"),
    ("~a", "eval"),

    # BoolOp
    ("a and b", "eval"),
    ("a or b", "eval"),
    ("a or b or c", "eval"),
    ("(a or b) or c", "eval"),
    ("a or (b or c)", "eval"),
    ("a and b and c", "eval"),
    ("(a and b) and c", "eval"),
    ("a and (b and c)", "eval"),

    # Exprs
    ("a+b-c/d", "eval"),
    ("a+b-c/d//e", "eval"),
    ("a+b-c/d//e%f", "eval"),
    ("a+b-c/d//e%f**g", "eval"),
    ("a+b-c/d//e%f**g|h&i^k<<l>>m", "eval"),
    ("a if b else c", "eval"),

    ("a==b", "eval"),
    ("a!=b", "eval"),
    ("a<b", "eval"),
    ("a<=b", "eval"),
    ("a>b", "eval"),
    ("a>=b", "eval"),
    ("a is b", "eval"),
    ("a is not b", "eval"),
    ("a in b", "eval"),
    ("a not in b", "eval"),

    ("a<b<c<d", "eval"),
    ("a==b<c>d", "eval"),
    ("(a==b)<c", "eval"),
    ("a==(b<c)", "eval"),
    ("(a==b)<(c>d)>e", "eval"),

    # trailers
    ("a()", "eval"),
    ("a(b)", "eval"),
    ("a(b,)", "eval"),
    ("a(b,c)", "eval"),
    ("a(b,*c)", "eval"),
    ("a(*b)", "eval"),
    ("a(*b,c)", "eval", SyntaxError),
    ("a(b,*c,**d)", "eval"),
    ("a(b,**c)", "eval"),
    ("a(a=b)", "eval"),
    ("a(a,a=b,*args,**kwargs)", "eval"),
    ("a(a,a=b,*args,e=f,**kwargs)", "eval"),
    ("a(b for c in d)", "eval"),
    ("a.b", "eval"),
    ("a.b.c.d", "eval"),
    ("a.b().c.d()()", "eval"),
    ("x[a]", "eval"),
    ("x[a,]", "eval"),
    ("x[a:b]", "eval"),
    ("x[:b]", "eval"),
    ("x[b:]", "eval"),
    ("x[:]", "eval"),
    ("x[a:b:c]", "eval"),
    ("x[:b:c]", "eval"),
    ("x[a::c]", "eval"),
    ("x[a:b:]", "eval"),
    ("x[::c]", "eval"),
    ("x[:b:]", "eval"),
    ("x[::c]", "eval"),
    ("x[::]", "eval"),
    ("x[a,p]", "eval"),
    ("x[a, b]", "eval"),
    ("x[a, b, c]", "eval"),
    ("x[a, b:c, ::d]", "eval"),
    ("x[a, b:c, ::d]", "eval"),
    ("x[0, 1:2, ::5, ...]", "eval"),

    # yield expressions
    ("(yield a,b)", "eval"),
    ("(yield from a)", "eval"),

    # statements
    ("del a,b", "exec"),
    ("del *a,*b", "exec"),
    ("pass", "exec"),
    ("break", "exec"),
    ("continue", "exec"),
    ("return", "exec"),
    ("return a", "exec"),
    ("return a,", "exec"),
    ("return a,b", "exec"),
    ("raise", "exec"),
    ("raise a", "exec"),
    ("raise a from b", "exec"),
    ("yield", "exec"),
    ("yield a", "exec"),
    ("yield a, b", "exec"),
    ("import a", "exec"),
    ("import a as b, c as d", "exec"),
    ("import a . b,c .d.e", "exec"),
    ("from a import b", "exec"),
    ("from a import b as c, d as e", "exec"),
    ("from a import b, c", "exec"),
    ("from a import (b, c)", "exec"),
    ("from a import *", "exec"),
    ("from . import b", "exec"),
    ("from .. import b", "exec"),
    ("from .a import (b, c,)", "exec"),
    ("from ..a import b", "exec"),
    ("from ...a import b", "exec"),
    ("from ....a import b", "exec"),
    ("from .....a import b", "exec"),
    ("from ......a import b", "exec"),
    ("global a", "exec"),
    ("global a, b", "exec"),
    ("global a, b, c", "exec"),
    ("nonlocal a", "exec"),
    ("nonlocal a, b", "exec"),
    ("nonlocal a, b, c", "exec"),
    ("assert True", "exec"),
    ("assert True, 'Bang'", "exec"),
    ("assert a == b, 'Bang'", "exec"),
    ("pass ; break ; continue", "exec"),

    # Compound statements
    ("while True: pass", "exec"),
    ("while True:\n pass\n", "exec"),
    ("while True:\n pass\nelse:\n return\n", "exec"),
    ("if True: pass", "exec"),
    ("if True:\n pass\n", "exec"),
    ("if True:\n pass\n\n", "exec"),
    ("""\
if True:
    pass
    continue
else:
    break
    pass
""", "exec"),
    ("""\
if a:
    continue
elif b:
    break
elif c:
    pass
elif c:
    continue
    pass
""", "exec"),
    ("""\
if a:
    continue
elif b:
    break
else:
    continue
    pass
""", "exec"),
    ("""\
if a:
    continue
elif b:
    break
elif c:
    pass
else:
    continue
    pass
""", "exec"),
    ("if lambda: None:\n pass\n", "exec"),
    ("for a in b: pass", "exec"),
    ("for a, b in b: pass", "exec"),
    ("for a, b in b:\n pass\nelse: break\n", "exec"),
    ("""\
try:
    pass
except:
    break
""", "exec"),
    ("""\
try:
    pass
except a:
    break
""", "exec"),
    ("""\
try:
    pass
except a as b:
    break
""", "exec"),
    ("""\
try:
    pass
except a:
    break
except:
    continue
except b as c:
    break
else:
    pass
""", "exec"),
    ("""\
try:
    pass
except:
    continue
finally:
    pass
""", "exec"),
    ("""\
try:
    pass
except:
    continue
else:
    break
finally:
    pass
""", "exec"),

    ("""\
with x:
    pass
""", "exec"),
    ("""\
with x as y:
    pass
""", "exec"),
    ("""\
with x as y, a as b, c, d as e:
    pass
    continue
""", "exec"),

    # Augmented assign
    ("a += b", "exec"),
    ("a -= b", "exec"),
    ("a *= b", "exec"),
    ("a /= b", "exec"),
    ("a -= b", "exec"),
    ("a %= b", "exec"),
    ("a &= b", "exec"),
    ("a |= b", "exec"),
    ("a ^= b", "exec"),
    ("a <<= b", "exec"),
    ("a >>= b", "exec"),
    ("a **= b", "exec"),
    ("a //= b", "exec"),
    ("a //= yield b", "exec"),
    ("a <> b", "exec", SyntaxError),
    ('''a.b += 1''', "exec"),

    # Assign
    ("a = b", "exec"),
    ("a = 007", "exec", SyntaxError, "illegal decimal with leading zero"),
    ("a = b = c", "exec"),
    ("a, b = 1, 2", "exec"),
    ("a, b = c, d = 1, 2", "exec"),
    ("a, b = *a", "exec"),
    ("a = yield a", "exec"),
    ('''a.b = 1''', "exec"),
    ("[e for e in [1, 2, 3]] = 3", "exec", SyntaxError),
    ("{e for e in [1, 2, 3]} = 3", "exec", SyntaxError),
    ("{e: e**2 for e in [1, 2, 3]} = 3", "exec", SyntaxError),
    ('''f() = 1''', "exec", SyntaxError),
    ('''lambda: x = 1''', "exec", SyntaxError),
    ('''(a + b) = 1''', "exec", SyntaxError),
    ('''(x for x in xs) = 1''', "exec", SyntaxError),
    ('''(yield x) = 1''', "exec", SyntaxError),
    ('''[x for x in xs] = 1''', "exec", SyntaxError),
    ('''{x for x in xs} = 1''', "exec", SyntaxError),
    ('''{x:x for x in xs} = 1''', "exec", SyntaxError),
    ('''{} = 1''', "exec", SyntaxError),
    ('''None = 1''', "exec", SyntaxError),
    ('''... = 1''', "exec", SyntaxError),
    ('''(a < b) = 1''', "exec", SyntaxError),
    ('''(a if b else c) = 1''', "exec", SyntaxError),

    # lambda
    ("lambda: a", "eval"),
    ("lambda: lambda: a", "eval"),
    ("lambda a: a", "eval"),
    ("lambda a, b: a", "eval"),
    ("lambda a, b,: a", "eval"),
    ("lambda a = b: a", "eval"),
    ("lambda a, b=c: a", "eval"),
    ("lambda a, *b: a", "eval"),
    ("lambda a, *b, c=d: a", "eval"),
    ("lambda a, *, c=d: a", "eval"),
    ("lambda a, *b, c=d, **kws: a", "eval"),
    ("lambda a, c=d, **kws: a", "eval"),
    ("lambda *args, c=d: a", "eval"),
    ("lambda *args, c=d, **kws: a", "eval"),
    ("lambda **kws: a", "eval"),

    # function
    ("def fn(): pass", "exec"),
    ("def fn(a): pass", "exec"),
    ("def fn(a, b): pass", "exec"),
    ("def fn(a, b,): pass", "exec"),
    ("def fn(a = b): pass", "exec"),
    ("def fn(a, b=c): pass", "exec"),
    ("def fn(a, *b): pass", "exec"),
    ("def fn(a, *b, c=d): pass", "exec"),
    ("def fn(a, *b, c=d, **kws): pass", "exec"),
    ("def fn(a, c=d, **kws): pass", "exec"),
    ("def fn(*args, c=d): pass", "exec"),
    ("def fn(a, *, c=d): pass", "exec"),
    ("def fn(*args, c=d, **kws): pass", "exec"),
    ("def fn(**kws): pass", "exec"),
    ("def fn() -> None: pass", "exec"),
    ("def fn(a:'potato') -> 'sausage': pass", "exec"),
    ("del f()", "exec", SyntaxError),

    # class
    ("class A: pass", "exec"),
    ("class A(): pass", "exec"),
    ("class A(B): pass", "exec"),
    ("class A(B,C): pass", "exec"),
    ("class A(B,C,D=F): pass", "exec"),
    ("class A(B,C,D=F,*AS,**KWS): pass", "exec"),

    # decorators
    ("""\
@dec
def fn():
    pass
""", "exec"),
    ("""\
@dec()
def fn():
    pass
""", "exec"),
    ("""\
@dec(a,b,c=d,*args,**kwargs)
def fn():
    pass
""", "exec"),
    ("""\
@dec1
@dec2()
@dec3(a)
@dec4(a,b)
def fn():
    pass
""", "exec"),
    ("""\
@dec1
@dec2()
@dec3(a)
@dec4(a,b)
class A(B):
    pass
""", "exec"),

    # single input
    ("", "single", SyntaxError),
    ("\n", "single", SyntaxError),
    ("pass\n", "single"),
    ("if True:\n   pass\n\n", "single"),
    ("while True:\n pass\nelse:\n return\n", "single"),
    # unfinished strings
    ("a='potato", "eval", SyntaxError),
    ("a='potato", "exec", SyntaxError),
    ("a='potato", "single", SyntaxError),
    ("a='''potato", "eval", SyntaxError),
    ("a='''potato", "exec", SyntaxError),
    ("a='''potato", "single", SyntaxError),
]

def dump(source, mode):
    """Dump source after parsing with mode"""
    a = ast.parse(source, mode=mode)
    return ast.dump(a, annotate_fields=True, include_attributes=False)

def escape(x):
    """Encode strings with backslashes for python/go"""
    return x.replace('\\', "\\\\").replace('"', r'\"').replace("\n", r'\n').replace("\t", r'\t')

def main():
    """Write grammar_data_test.go"""
    path = "grammar_data_test.go"
    year = datetime.datetime.now().year
    out = ["""// Copyright {year} The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test data generated by make_grammar_test.py - do not edit

package parser

import (
"github.com/go-python/gpython/py"
)

var grammarTestData = []struct {{
in   string
mode string
out  string
exceptionType *py.Type
errString string
}}{{""".format(year=year)]
    for x in inp:
        source, mode = x[:2]
        if len(x) > 2:
            exc = x[2]
            errString = (x[3] if len(x) > 3 else "")
            try:
                dump(source, mode)
            except exc as e:
                error = e.msg
            else:
                raise ValueError("Expecting exception %s" % exc)
            if errString != "":
                error = errString # override error string
            dmp = ""
            exc_name = "py.%s" % exc.__name__
        else:
            dmp = dump(source, mode)
            exc_name = "nil"
            error = ""
        out.append('{"%s", "%s", "%s", %s, "%s"},' % (escape(source), mode, escape(dmp), exc_name, escape(error)))
    out.append("}")
    print("Writing %s" % path)
    with open(path, "w") as f:
        f.write("\n".join(out))
        f.write("\n")
    subprocess.check_call(["gofmt", "-w", path])

if __name__ == "__main__":
    main()
