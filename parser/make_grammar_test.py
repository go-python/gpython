#!/usr/bin/env python3
"""
Read in grammar_test.go, and re-write the tests section
"""

import sys
import ast

inp = [
    ("", "exec"),
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

    ("( a for a in ab )", "eval"),
    ("( a for a, in ab )", "eval"),
    ("( a for a, b in ab )", "eval"),
    ("( a for a in ab if a )", "eval"),
    ("( a for a in ab if a if b if c )", "eval"),
    ("( a for a in ab for A in AB )", "eval"),
    ("( a for a in ab if a if b for A in AB if c )", "eval"),

    ("[ a for a in ab ]", "eval"),
    ("[ a for a, in ab ]", "eval"),
    ("[ a for a, b in ab ]", "eval"),
    ("[ a for a in ab if a ]", "eval"),
    ("[ a for a in ab if a if b if c ]", "eval"),
    ("[ a for a in ab for A in AB ]", "eval"),
    ("[ a for a in ab if a if b for A in AB if c ]", "eval"),

    ("{ a for a in ab }", "eval"),
    ("{ a for a, in ab }", "eval"),
    ("{ a for a, b in ab }", "eval"),
    ("{ a for a in ab if a }", "eval"),
    ("{ a for a in ab if a if b if c }", "eval"),
    ("{ a for a in ab for A in AB }", "eval"),
    ("{ a for a in ab if a if b for A in AB if c }", "eval"),

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
    #("a(*b,c)", "eval"), -test error
    ("a(b,*c,**d)", "eval"),
    ("a(b,**c)", "eval"),
    ("a(a=b)", "eval"),
    ("a(a,a=b,*args,**kwargs)", "eval"),
    ("a(a,a=b,*args,e=f,**kwargs)", "eval"),
    ("a.b", "eval"),
    ("a.b.c.d", "eval"),
    ("a.b().c.d()()", "eval"),
    ("x[a]", "eval"),
    ("x[a:b]", "eval"),
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

    ("del a,b", "exec"),
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
    ("import a . b,c .d.e", "exec"),
    ("from a import b", "exec"),
    ("from a import b, c", "exec"),
    ("from a import (b, c)", "exec"),
    ("from a import *", "exec"),
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

    # Compound statements
    ("while True: pass", "exec"),
    ("while True:\n pass\n", "exec"),
    ("while True:\n pass\nelse:\n return\n", "exec"),
    ("if True: pass", "exec"),
    ("if True:\n pass\n", "exec"),
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
]

def dump(source, mode):
    """Dump source after parsing with mode"""
    a = ast.parse(source, mode=mode)
    return ast.dump(a, annotate_fields=True, include_attributes=False)

def escape(x):
    """Encode strings with backslashes for python/go"""
    return x.replace('\\', "\\\\").replace('"', r'\"').replace("\n", r'\n').replace("\t", r'\t')

def main():
    """Read in grammar_test.go, and re-write the tests section"""
    path = "grammar_test.go"
    with open(path) as f:
        grammar_test = f.read()
    lines = grammar_test.split("\n")
    while lines[-1] == "":
        lines = lines[:-1]
    out = []
    in_tests = False
    for line in lines:
        if "START TESTS" in line:
            out.append(line)
            out.append("\t\t// *** Tests auto generated by make_grammar_test.py - do not edit ***")
            for source, mode in inp:
                out.append('\t\t{"%s", "%s", "%s"},' % (escape(source), mode, escape(dump(source, mode))))
            in_tests = True
        elif "END TESTS" in line:
            in_tests = False
        if not in_tests:
            out.append(line)
    print("Rewriting %s" % path)
    with open(path, "w") as f:
        f.write("\n".join(out))
        f.write("\n")

if __name__ == "__main__":
    main()
