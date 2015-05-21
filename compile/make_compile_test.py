#!/usr/bin/env python3.4
"""
Write compile_data_test.go
"""

import sys
import ast
import subprocess

inp = [
    # Constants
    ('''1''', "eval"),
    ('''"hello"''', "eval"),
    ('''a''', "eval"),
    ('''b"hello"''', "eval"),
    # BinOps - strange operations to defeat constant optimizer!
    ('''"a"+1''', "eval"),
    ('''"a"-1''', "eval"),
    ('''"a"*"b"''', "eval"),
    ('''"a"/1''', "eval"),
    ('''"a"%1''', "eval"),
    ('''"a"**1''', "eval"),
    ('''"a"<<1''', "eval"),
    ('''"a">>1''', "eval"),
    ('''"a"|1''', "eval"),
    ('''"a"^1''', "eval"),
    ('''"a"&1''', "eval"),
    ('''"a"//1''', "eval"),
    ('''a+a''', "eval"),
    ('''"a"*"a"''', "eval"),
    ('''1''', "exec"),
    ('''1\n"hello"''', "exec"),
    ('''a+a''', "exec"),
    # UnaryOps
    ('''~ "a"''', "eval"),
    ('''not "a"''', "eval"),
    ('''+"a"''', "eval"),
    ('''-"a"''', "eval"),
    # Bool Ops
    ('''1 and 2''', "eval"),
    ('''1 and 2 and 3 and 4''', "eval"),
    ('''1 and 2''', "eval"),
    ('''1 or 2''', "eval"),
    ('''1 or 2 or 3 or 4''', "eval"),
    # With brackets
    ('''"1"+"2"*"3"''', "eval"),
    ('''"1"+("2"*"3")''', "eval"),
    ('''(1+"2")*"3"''', "eval"),
    # If expression
    ('''(a if b else c)+0''', "eval"),
    # Compare
    ('''a == b''', "eval"),
    ('''a != b''', "eval"),
    ('''a < b''', "eval"),
    ('''a <= b''', "eval"),
    ('''a > b''', "eval"),
    ('''a >= b''', "eval"),
    ('''a is b''', "eval"),
    ('''a is not b''', "eval"),
    ('''a in b''', "eval"),
    ('''a not in b''', "eval"),
    ('''(a < b < c)+0''', "eval"),
    ('''(a < b < c < d)+0''', "eval"),
    ('''(a < b < c < d < e)+0''', "eval"),
    # tuple
    ('''()''', "eval"),
    #('''(1,)''', "eval"),
    #('''(1,1)''', "eval"),
    #('''(1,1,3,1)''', "eval"),
    ('''(a,)''', "eval"),
    ('''(a,b)''', "eval"),
    ('''(a,b,c,d)''', "eval"),
    # list
    ('''[]''', "eval"),
    ('''[1]''', "eval"),
    ('''[1,1]''', "eval"),
    ('''[1,1,3,1]''', "eval"),
    ('''[a]''', "eval"),
    ('''[a,b]''', "eval"),
    ('''[a,b,c,d]''', "eval"),
    # named constant
    ('''True''', "eval"),
    ('''False''', "eval"),
    ('''None''', "eval"),
    # attribute
    ('''a.b''', "eval"),
    ('''a.b.c''', "eval"),
    ('''a.b.c.d''', "eval"),
    ('''a.b = 1''', "exec"),
    ('''a.b.c.d = 1''', "exec"),
    ('''a.b += 1''', "exec"),
    ('''a.b.c.d += 1''', "exec"),
    ('''del a.b''', "exec"),
    ('''del a.b.c.d''', "exec"),
    # dict
    ('''{}''', "eval"),
    ('''{1:2,a:b}''', "eval"),
    # set
    # ('''set()''', "eval"),
    ('''{1}''', "eval"),
    ('''{1,2,a,b}''', "eval"),
    # lambda
    ('''lambda: 0''', "eval"),
    ('''lambda x: 2*x''', "eval"),
    ('''lambda a,b=42,*args,**kw: a*b*args*kw''', "eval"),
    # pass statment
    ('''pass''', "exec"),
    # expr statement
    ('''(a+b)''', "exec"),
    ('''(a+\nb+\nc)\n''', "exec"),
    # assert
    ('''assert a, "hello"''', "exec"),
    ('''assert 1, 2''', "exec"),
    ('''assert a''', "exec"),
    ('''assert 1''', "exec"),
    # assign
    ('''a = 1''', "exec"),
    ('''a = b = c = 1''', "exec"),
    ('''a[1] = 1''', "exec"),
    # aug assign
    ('''a+=1''', "exec"),
    ('''a-=1''', "exec"),
    ('''a*=b''', "exec"),
    ('''a/=1''', "exec"),
    ('''a%=1''', "exec"),
    ('''a**=1''', "exec"),
    ('''a<<=1''', "exec"),
    ('''a>>=1''', "exec"),
    ('''a|=1''', "exec"),
    ('''a^=1''', "exec"),
    ('''a&=1''', "exec"),
    ('''a//=1''', "exec"),
    ('''a[1]+=1''', "exec"),
    # delete
    ('''del a''', "exec"),
    ('''del a, b''', "exec"),
    ('''del a[1]''', "exec"),
    ('''\
def fn(b):
 global a
 del a
 c = 1
 def nested(d):
   nonlocal b
   e = b+c+d+e
   f(e)
   del b,c,d,e
''', "exec"),
    # raise
    ('''raise''', "exec"),
    ('''raise a''', "exec"),
    ('''raise a from b''', "exec"),
    # if
    ('''if a: b = c''', "exec"),
    ('''if a:\n b = c\nelse:\n c = d\n''', "exec"),
    # while
    ('''while a:\n b = c''', "exec"),
    ('''while a:\n b = c\nelse:\n b = d\n''', "exec"),
    ('''while a:\n if b: break\n b = c\n''', "exec"),
    ('''while a:\n if b: continue\n b = c\n''', "exec"),
    ('''continue''', "exec", SyntaxError),
    ('''break''', "exec", SyntaxError),
    # for
    ('''for a in b: pass''', "exec"),
    ('''for a in b:\n if a:\n  break\n c = e\nelse: c = d\n''', "exec"),
    ('''for a in b:\n if a:\n  continue\n c = e\nelse: c = d\n''', "exec"),
    # call
    ('''f()''', "eval"),
    ('''f(a)''', "eval"),
    ('''f(a,b,c)''', "eval"),
    ('''f(A=a)''', "eval"),
    ('''f(a, b, C=d, D=d)''', "eval"),
    ('''f(*args)''', "eval"),
    ('''f(*args, **kwargs)''', "eval"),
    ('''f(**kwargs)''', "eval"),
    ('''f(a, b, *args)''', "eval"),
    ('''f(a, b, *args, d=e, **kwargs)''', "eval"),
    ('''f(a, d=e, **kwargs)''', "eval"),
    # def
    ('''def fn(): pass''', "exec"),
    ('''def fn(a): pass''', "exec"),
    ('''def fn(a,b,c): pass''', "exec"),
    ('''def fn(a,b=1,c=2): pass''', "exec"),
    ('''def fn(a,*arg,b=1,c=2): pass''', "exec"),
    ('''def fn(a,*arg,b=1,c=2,**kwargs): pass''', "exec"),
    ('''def fn(a:"a",*arg:"arg",b:"b"=1,c:"c"=2,**kwargs:"kw") -> "ret": pass''', "exec"),
    ('''def fn(): a+b''', "exec"),
    ('''def fn(a,b): a+b+c+d''', "exec"),
    ('''\
def fn(a):
    global b
    b = a''', "exec"),
    ('''def fn(): return''', "exec"),
    ('''def fn(): return a''', "exec"),
    ('''def fn():\n "docstring"\n return True''', "exec"),
    ('''\
def outer(o):
    def inner(i):
       x = 2''', "exec"),
    ('''\
def outer(o1,o2):
    x = 1
    def inner(i1,i2):
       nonlocal x
       x = 2
       def inner2(s):
           return 2*s
       f = inner2(x)
       l = o1+o2+i1+i2+f
       return l
    return inner''', "exec"),
    ('''\
def outer(o):
    x = 17
    return lambda a,b=42,*args,**kw: a*b*args*kw*x*o''', "exec"),
    ('''\
@wrap
def fn(o):
    return o''', "exec"),
    ('''\
@wrap1
@wrap2("potato", 2)
@wrap3("sausage")
@wrap4
def fn(o):
    return o''', "exec"),
    ('''\
def outer(o):
    @wrap1
    @wrap2("potato", o)
    def inner(i):
        return o+i''', "exec"),
    # module docstrings
    ('''\
# Module
"""
A module docstring
"""
''', "exec"),
    ('''\
# Empty docstring
""
''', "exec"),
    # class
    ('''\
class Dummy:
    pass
''', "exec"),
    ('''\
@d1
@d2
class Dummy(a,b,c=d):
    "A class"
    pass
''', "exec"),
    ('''\
class Dummy:
    def method(self):
        return self+1
''', "exec"),
    ('''\
@class1
@class2(arg2)
class Dummy:
    "Dummy"
    @fn1
    @fn2(arg2)
    def method(self):
        "method"
        return self+1
    def method2(self, m2):
        "method2"
        return self.method()+m2
''', "exec"),
    ('''\
def closure_class(a):
    b = 42
    class AClass:
        def method(self, c):
            return a+b+c
    return AClass
''', "exec"),
    ('''\
@potato
@sausage()
class A(a,b,c=\"1\",d=\"2\",*args,**kwargs):
    VAR = x
    def method(self):
        super().method()
        return VAR
''', "exec"),
    ('''\
def outer(x):
    class DeRefTest:
        VAR = x
''', "exec"),
    # comprehensions
    ('''[ x for x in xs ]''', "eval"),
    ('''{ x: y for x in xs }''', "eval"),
    ('''{ x for x in xs }''', "eval"),
    ('''( x for x in xs )''', "eval"),
    ('''[ x for x in xs if a ]''', "eval"),
    ('''{ x: y for x in xs if a if b }''', "eval"),
    ('''{ x for x in xs if a}''', "eval"),
    ('''( x for x in xs if a if b if c)''', "eval"),
    ('''{ x for x in [ x for x in xs if c if d ] if a if b}''', "eval"),
    ('''[ (x,y,z) for x in xs for y in ys for z in zs ]''', "eval"),
    ('''{ (x,y,z) for x in xs for y in ys if a if b for z in zs if c if d }''', "eval"),
    ('''{ x:(y,z) for x in xs for y in ys for z in zs }''', "eval"),
    ('''( (x,y,z) for x in xs for y in ys if a if b for z in zs if c if d )''', "eval"),
    # with
    ('''\
with a:
    f()
''', "exec"),
    ('''\
with a() as b:
    f(b)
''', "exec"),
    ('''\
with A() as a, B() as b:
    f(a,b)
''', "exec"),
    ('''\
with A() as a:
    with B() as b:
        f(a,b)
''', "exec"),
    # try/except/finally/else
    ('''\
ok = False
try:
    raise SyntaxError
except SyntaxError:
    ok = True
assert ok
''', "exec"),
    ('''\
ok = False
try:
    raise SyntaxError
except SyntaxError as e:
    ok = True
assert ok
''', "exec"),
    ('''\
try:
    f()
except Exception:
    h()
''', "exec"),
    ('''\
try:
    f()
except Exception as e:
    h(e)
except (Exception1, Exception2) as e:
    i(e)
except:
    j()
else:
    potato()
''', "exec"),
    ('''\
try:
    f()
except:
    j()
except Exception as e:
    h(e)
    ''', "exec", SyntaxError),
    ('''\
try:
    f()
finally:
    j()
    ''', "exec"),
    ('''\
try:
    f()
except Exception as e:
    h(e)
finally:
    j()
    ''', "exec"),
    # import / from import
    ('''import mod''', "exec"),
    ('''import mod1, mod2, mod3''', "exec"),
    ('''import mod as pod, mod2 as pod2''', "exec"),
    ('''import mod1.mod2''', "exec"),
    ('''import mod1.mod2.mod3''', "exec"),
    ('''import mod1.mod2.mod3.mod4''', "exec"),
    ('''import mod1.mod2.mod3.mod4 as potato''', "exec"),
    ('''from mod import a''', "exec"),
    ('''from mod1.mod2.mod3 import *''', "exec"),
    ('''from mod1.mod2.mod3 import a as aa, b as bb, c''', "exec"),
    # yield
    ('''yield''', "exec", SyntaxError),
    ('''yield potato''', "exec", SyntaxError),
    ('''\
def f():
    yield
    ''', "exec"),
    ('''\
def f():
    yield potato
    ''', "exec"),
    # yield from
    ('''yield from range(10)''', "exec", SyntaxError),
    ('''\
def f():
    yield from range(10)
    ''', "exec"),
    # ellipsis
    ('''...''', "exec"),
    # starred...
    ('''*a = t''', "exec", SyntaxError),
    ('''a, *b = t''', "exec"),
    ('''(a, *b) = t''', "exec"),
    ('''[a, *b] = t''', "exec"),
    ('''a, *b, c = t''', "exec"),
    ('''a, *b, *c = t''', "exec", SyntaxError),
    ('''a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,*a = t''', "exec", SyntaxError),
    ('''a, b, *c''', "exec", SyntaxError),
    ('''a, (b, c), d = t''', "exec"),
    # subscript - load
    ("x[a]", "exec"),
    ("x[a:b]", "exec"),
    ("x[:b]", "exec"),
    ("x[b:]", "exec"),
    ("x[:]", "exec"),
    ("x[a:b:c]", "exec"),
    ("x[:b:c]", "exec"),
    ("x[a::c]", "exec"),
    ("x[a:b:]", "exec"),
    ("x[::c]", "exec"),
    ("x[:b:]", "exec"),
    ("x[::c]", "exec"),
    ("x[::]", "exec"),
    ("x[a,p]", "exec"),
    ("x[a, b]", "exec"),
    ("x[a, b, c]", "exec"),
    ("x[a, b:c, ::d]", "exec"),
    ("x[0, 1:2, ::5, ...]", "exec"),
    # subscript - store
    ("x[a] = y", "exec"),
    ("x[a:b] = y", "exec"),
    ("x[:b] = y", "exec"),
    ("x[b:] = y", "exec"),
    ("x[:] = y", "exec"),
    ("x[a:b:c] = y", "exec"),
    ("x[:b:c] = y", "exec"),
    ("x[a::c] = y", "exec"),
    ("x[a:b:] = y", "exec"),
    ("x[::c] = y", "exec"),
    ("x[:b:] = y", "exec"),
    ("x[::c] = y", "exec"),
    ("x[::] = y", "exec"),
    ("x[a,p] = y", "exec"),
    ("x[a, b] = y", "exec"),
    ("x[a, b, c] = y", "exec"),
    ("x[a, b:c, ::d] = y", "exec"),
    ("x[0, 1:2, ::5, ...] = y", "exec"),
    # subscript - aug assign (AugLoad and AugStore)
    ("x[a] += y", "exec"),
    ("x[a:b] += y", "exec"),
    ("x[:b] += y", "exec"),
    ("x[b:] += y", "exec"),
    ("x[:] += y", "exec"),
    ("x[a:b:c] += y", "exec"),
    ("x[:b:c] += y", "exec"),
    ("x[a::c] += y", "exec"),
    ("x[a:b:] += y", "exec"),
    ("x[::c] += y", "exec"),
    ("x[:b:] += y", "exec"),
    ("x[::c] += y", "exec"),
    ("x[::] += y", "exec"),
    ("x[a,p] += y", "exec"),
    ("x[a, b] += y", "exec"),
    ("x[a, b, c] += y", "exec"),
    ("x[a, b:c, ::d] += y", "exec"),
    ("x[0, 1:2, ::5, ...] += y", "exec"),
    # subscript - delete
    ("del x[a]", "exec"),
    ("del x[a:b]", "exec"),
    ("del x[:b]", "exec"),
    ("del x[b:]", "exec"),
    ("del x[:]", "exec"),
    ("del x[a:b:c]", "exec"),
    ("del x[:b:c]", "exec"),
    ("del x[a::c]", "exec"),
    ("del x[a:b:]", "exec"),
    ("del x[::c]", "exec"),
    ("del x[:b:]", "exec"),
    ("del x[::c]", "exec"),
    ("del x[::]", "exec"),
    ("del x[a,p]", "exec"),
    ("del x[a, b]", "exec"),
    ("del x[a, b, c]", "exec"),
    ("del x[a, b:c, ::d]", "exec"),
    ("del x[0, 1:2, ::5, ...]", "exec"),
    # continue
    ('''\
try:
    continue
except:
    pass
    ''', "exec", SyntaxError),
    ('''\
try:
    pass
except:
    continue
    ''', "exec", SyntaxError),
    ('''\
for x in xs:
    try:
        f()
    except:
        continue
    f()
    ''', "exec"),
    ('''\
for x in xs:
    try:
        f()
        continue
    finally:
        f()
    ''', "exec"),
    ('''\
for x in xs:
    try:
        f()
    finally:
        continue
    ''', "exec", SyntaxError),
    ('''\
for x in xs:
    try:
        f()
    finally:
        try:
            continue
        except:
             pass
    ''', "exec", SyntaxError),
    ('''\
try:
    continue
except:
    pass
    ''', "exec", SyntaxError),
    ('''\
try:
    pass
except:
    continue
    ''', "exec", SyntaxError),
    ('''\
while truth():
    try:
        f()
    except:
        continue
    f()
    ''', "exec"),
    ('''\
while truth():
    try:
        f()
        continue
    finally:
        f()
    ''', "exec"),
    ('''\
while truth():
    try:
        f()
    finally:
        continue
    ''', "exec", SyntaxError),
    ('''\
while truth():
    try:
        f()
    finally:
        try:
            continue
        except:
             pass
    ''', "exec", SyntaxError),
    # interactive
    ('''print("hello world!")\n''', "single"),
    # FIXME ('''if True:\n "hello world!"\n''', "single"),
    # FIXME ('''def fn(x):\n "hello world!"\n''', "single"),

 ]

def string(s):
    if isinstance(s, str):
        return '"%s"' % s
    elif isinstance(s, bytes):
        out = '"'
        for b in s:
            out += "\\x%02x" % b
        out += '"'
        return out
    else:
        raise AssertionError("Unknown string %r" % s)

def strings(ss):
    """Dump a list of py strings into go format"""
    return "[]string{"+",".join(string(s) for s in ss)+"}"

codeObjectType = type(strings.__code__)

def const(x):
    if isinstance(x, str):
        return 'py.String("%s")' % x.encode("unicode-escape").decode("utf-8")
    elif isinstance(x, bool):
        if x:
            return 'py.True'
        return 'py.False'
    elif isinstance(x, int):
        return 'py.Int(%d)' % x
    elif isinstance(x, float):
        return 'py.Float(%g)' % x
    elif isinstance(x, bytes):
        return 'py.Bytes("%s")' % x.decode("latin1")
    elif isinstance(x, tuple):
        return 'py.Tuple{%s}' % ",".join(const(y) for y in x)
    elif isinstance(x, codeObjectType):
        return "\n".join([
            "&py.Code{",
            "Argcount: %s," % x.co_argcount,
            "Kwonlyargcount: %s," % x.co_kwonlyargcount,
            "Nlocals: %s," % x.co_nlocals,
            "Stacksize: %s," % x.co_stacksize,
            "Flags: %s," % x.co_flags,
            "Code: %s," % string(x.co_code),
            "Consts: %s," % consts(x.co_consts),
            "Names: %s," % strings(x.co_names),
            "Varnames: %s," % strings(x.co_varnames),
            "Freevars: %s," % strings(x.co_freevars),
            "Cellvars: %s," % strings(x.co_cellvars),
            # "Cell2arg    []byte // Maps cell vars which are arguments".
            "Filename: %s," % string(x.co_filename),
            "Name: %s," % string(x.co_name),
            "Firstlineno: %d," % x.co_firstlineno,
            "Lnotab: %s," % string(x.co_lnotab),
            "}",
        ])
    elif x is None:
        return 'py.None'
    elif x is ...:
        return 'py.Ellipsis'
    else:
        raise AssertionError("Unknown const %r" % x)

def consts(xs):
    return "[]py.Object{"+",".join(const(x) for x in xs)+"}"
    
def _compile(source, mode):
    """compile source with mode"""
    a = compile(source=source, filename="<string>", mode=mode, dont_inherit=True, optimize=0)
    return a, const(a)

def escape(x):
    """Encode strings with backslashes for python/go"""
    return x.replace('\\', "\\\\").replace('"', r'\"').replace("\n", r'\n').replace("\t", r'\t')

def main():
    """Write compile_data_test.go"""
    path = "compile_data_test.go"
    out = ["""// Test data generated by make_compile_test.py - do not edit

package compile

import (
"github.com/ncw/gpython/py"
)

var compileTestData = []struct {
in   string
mode string // exec, eval or single
out  *py.Code
exceptionType *py.Type
errString string
}{"""]
    for x in inp:
        source, mode = x[:2]
        if len(x) > 2:
            exc = x[2]
            try:
                _compile(source, mode)
            except exc as e:
                error = e.msg
            else:
                raise ValueError("Expecting exception %s" % exc)
            gostring = "nil"
            exc_name = "py.%s" % exc.__name__
        else:
            code, gostring = _compile(source, mode)
            exc_name = "nil"
            error = ""
        out.append('{"%s", "%s", %s, %s, "%s"},' % (escape(source), mode, gostring, exc_name, escape(error)))
    out.append("}")
    print("Writing %s" % path)
    with open(path, "w") as f:
        f.write("\n".join(out))
        f.write("\n")
    subprocess.check_call(["gofmt", "-w", path])

if __name__ == "__main__":
    main()
