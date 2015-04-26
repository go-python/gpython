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
    # ('''(1,)''', "eval"),
    # ('''(1,1)''', "eval"),
    # ('''(1,1,3,1)''', "eval"),
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
    # dict
    ('''{}''', "eval"),
    ('''{1:2,a:b}''', "eval"),
    # set
    # ('''set()''', "eval"),
    ('''{1}''', "eval"),
    ('''{1,2,a,b}''', "eval"),
    # lambda
    ('''lambda: 0''', "eval"),
    #('''lambda x: 2*x''', "eval"),
    # pass statment
    ('''pass''', "exec"),
    # expr statement
    ('''(a+b)''', "exec"),
    #('''(a+\nb+\nc)\n''', "exec"),
    # assert
    ('''assert a, "hello"''', "exec"),
    ('''assert 1, 2''', "exec"),
    ('''assert a''', "exec"),
    ('''assert 1''', "exec"),
    # assign
    ('''a = 1''', "exec"),
    ('''a = b = c = 1''', "exec"),
    # FIXME ('''a[1] = 1''', "exec"),
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
    # FIXME ('''a[1]+=1''', "exec"),
    # delete
    ('''del a''', "exec"),
    ('''del a, b''', "exec"),
    # FIXME ('''del a[1]''', "exec"),
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
    #('''def fn(a,b): a+b+c+d''', "exec"),

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
        return 'py.String("%s")' % x
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
