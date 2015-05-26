#!/usr/bin/env python3.4
"""
Write symtable_data_test.go
"""

import sys
import ast
import subprocess
import dis
from symtable import symtable
try:
    # run ./build.sh in the readsymtab directory to create this module
    from readsymtab import readsymtab
    use_readsymtab = True
except ImportError:
    import ctypes
    use_readsymtab = False
    print("Using compiler dependent code to read bitfields - compile readsymtab module to be sure!")

inp = [
    ('''1''', "eval"),
    ('''a*b*c''', "eval"),
    # Functions
    ('''def fn(): pass''', "exec"),
    ('''def fn(a,b):\n e=1\n return a*b*c*d*e''', "exec"),
    ('''def fn(*args,a=2,b=3,**kwargs): return (args,a,b,kwargs)''', "exec"),
    ('''def fn(a,b):\n def nested(c,d):\n  return a*b*c*d*e''', "exec"),
    ('''\
def fn(a:A,*arg:ARG,b:B=BB,c:C=CC,**kwargs:KW) -> RET:
    def fn(A,b):
        e=1
        return a*arg*b*c*kwargs*A*e*glob''', "exec"),
    ('''\
def fn(a):
    global b
    b = a''', "exec"),
    ('''\
def fn(a):
    global b
    global b
    return b''', "exec"),
    ('''\
def inner():
  print(x)
  global x
''', "exec"),
    ('''\
def fn(a):
    b = 6
    global b
    b = a''', "exec"),
    ('''\
def fn(a=b,c=1):
    return a+b''', "exec"),
    ('''\
@sausage
@potato(beans)
def outer():
   x = 1
   def inner():
       nonlocal x
       x = 2''', "exec"),
    ('''\
def fn(a):
    nonlocal b
    ''', "exec", SyntaxError),
    ('''\
def outer():
   def inner():
       nonlocal x
       x = 2''', "exec", SyntaxError),
    ('''\
def outer():
   x = 1
   def inner():
       print(x)
       nonlocal x
''', "exec"),
    ('''\
def outer():
   x = 1
   def inner():
       x = 2
       nonlocal x''', "exec"),
    ('''\
def outer():
   x = 1
   def inner(x):
       nonlocal x''', "exec", SyntaxError),
    ('''\
def outer():
   x = 1
   def inner(x):
       global x''', "exec", SyntaxError),
    ('''\
def outer():
   def inner():
       global x
       nonlocal x
       ''', "exec", SyntaxError),
    ('''\
def outer():
   x = 1
   def inner():
       y = 2
       return x + y + z
''', "exec"),
    ('''\
def outer():
   global x
   def inner():
       return x
''', "exec"),
    ('''\
nonlocal x
''', "exec", SyntaxError),
    ('''def fn(a,a): pass''', "exec", SyntaxError),
    # List Comp
    ('''[ x for x in xs ]''', "exec"),
    ('''[ x+y for x in xs for y in ys ]''', "exec"),
    ('''[ x+y+z for x in xs if x if y if z if r for y in ys if x if y if z if p for z in zs if x if y if z if q]''', "exec"),
    ('''[ x+y for x in [ x for x in xs ] ]''', "exec"),
    ('''[ x for x in xs ]\n[ y for y in ys ]''', "exec"),
    # Generator expr
    ('''( x for x in xs )''', "exec"),
    ('''( x+y for x in xs for y in ys )''', "exec"),
    ('''( x+y+z for x in xs if x if y if z if r for y in ys if x if y if z if p for z in zs if x if y if z if q)''', "exec"),
    ('''( x+y for x in ( x for x in xs ) )''', "exec"),
    # Set comp
    ('''{ x for x in xs }''', "exec"),
    ('''{ x+y for x in xs for y in ys }''', "exec"),
    ('''{ x+y+z for x in xs if x if y if z if r for y in ys if x if y if z if p for z in zs if x if y if z if q}''', "exec"),
    ('''{ x+y for x in { x for x in xs } }''', "exec"),
    # Dict comp
    ('''{ x:1 for x in xs }''', "exec"),
    ('''{ x+y:1 for x in xs for y in ys }''', "exec"),
    ('''{ x+y+z:1 for x in xs if x if y if z if r for y in ys if x if y if z if p for z in zs if x if y if z if q}''', "exec"),
    ('''{ x+y:k for k, x in { x:1 for x in xs } }''', "exec"),
    # Class
    ('''\
@potato
@sausage()
class A(a,b,c="1",d="2",*args,**kwargs):
    VAR = x
    def method(self):
        super().method()
        return VAR
''', "exec"),
    # Lambda
    ('''lambda: x''', "exec"),
    ('''lambda y: x+y''', "exec"),
    ('''lambda a,*arg,b=BB,c=CC,**kwargs: POTATO+a+arg+b+c+kwargs''', "exec"),
    # With
    ('''\
with x() as y:
  y.floop()
print(y)
''', "exec"),
    # try, except
    ('''\
try:
  something()
except RandomError as e:
  print(e)
print(e)
''', "exec"),
    # Import
    ('''import potato''', "exec"),
    ('''import potato.sausage''', "exec"),
    ('''from potato import sausage''', "exec"),
    ('''from potato import sausage as salami''', "exec"),
    ('''from potato import *''', "exec"),
    ('''\
def fn():
  from potato import *
''', "exec", SyntaxError),
    # Yield
    ('''\
def f():
    yield
    ''', "exec"),
    ('''\
def f():
    yield from range(10)
    ''', "exec"),
]

def dump_bool(b):
    return ("true" if b else "false")

def dump_strings(ss):
    return "[]string{"+",".join([ '"%s"' % s for s in ss ])+"}"

# Scope numbers to names (from symtable.h)
SCOPES = {
    1: "ScopeLocal",
    2: "ScopeGlobalExplicit",
    3: "ScopeGlobalImplicit",
    4: "ScopeFree",
    5: "ScopeCell",
}

#def-use flags to names (from symtable.h)
DEF_FLAGS = (
    ("DefGlobal", 1),      # global stmt
    ("DefLocal", 2),       # assignment in code block
    ("DefParam", 2<<1),    # formal parameter
    ("DefNonlocal", 2<<2), # nonlocal stmt
    ("DefUse", 2<<3),      # name is used
    ("DefFree", 2<<4),     # name used but not defined in nested block
    ("DefFreeClass", 2<<5),# free variable from class's method
    ("DefImport", 2<<6),   # assignment occurred via import
)

#opt flags flags to names (from symtable.h)
OPT_FLAGS = (
    ("optImportStar", 1),
    ("optTopLevel", 2),
)

BLOCK_TYPES = {
    "function": "FunctionBlock",
    "class": "ClassBlock",
    "module": "ModuleBlock",
}

def dump_flags(flag_bits, flags_dict):
    """Dump the bits in flag_bits using the flags_dict"""
    flags = []
    for name, mask in flags_dict:
        if (flag_bits & mask) != 0:
            flags.append(name)
    if not flags:
        flags = ["0"]
    return "|".join(flags)

def dump_symtable(st):
    """Dump the symtable"""
    out = "&SymTable{\n"
    out += 'Type:%s,\n' % BLOCK_TYPES[st.get_type()] # Return the type of the symbol table. Possible values are 'class', 'module', and 'function'.
    out += 'Name:"%s",\n' % st.get_name() # Return the table’s name. This is the name of the class if the table is for a class, the name of the function if the table is for a function, or 'top' if the table is global (get_type() returns 'module').

    out += 'Lineno:%s,\n' % st.get_lineno() # Return the number of the first line in the block this table represents.
    out += 'Unoptimized:%s,\n' % dump_flags(st._table.optimized, OPT_FLAGS) # Return False if the locals in this table can be optimized.
    out += 'Nested:%s,\n' % dump_bool(st.is_nested()) # Return True if the block is a nested class or function.

    if use_readsymtab:
    # Use readsymtab modules to read the bitfields which aren't normally exported
        free, child_free, generator, varargs, varkeywords, returns_value, needs_class_closure =  readsymtab(st._table)
        out += 'Free:%s,\n' % dump_bool(free)
        out += 'ChildFree:%s,\n' % dump_bool(child_free)
        out += 'Generator:%s,\n' % dump_bool(generator)
        out += 'Varargs:%s,\n' % dump_bool(varargs)
        out += 'Varkeywords:%s,\n' % dump_bool(varkeywords)
        out += 'ReturnsValue:%s,\n' % dump_bool(returns_value)
        out += 'NeedsClassClosure:%s,\n' % dump_bool(needs_class_closure)
    else:
        # Use ctypes to read the bitfields which aren't normally exported
        # FIXME compiler dependent!
        base_addr = id(st._table) + ctypes.sizeof(ctypes.c_long)*8+ ctypes.sizeof(ctypes.c_int)*3
        flags = int.from_bytes(ctypes.c_int.from_address(base_addr), sys.byteorder)
        out += 'Free:%s,\n' % dump_bool(flags & (1 << 0))
        out += 'ChildFree:%s,\n' % dump_bool(flags & (1 << 1))
        out += 'Generator:%s,\n' % dump_bool(flags & (1 << 2))
        out += 'Varargs:%s,\n' % dump_bool(flags & (1 << 3))
        out += 'Varkeywords:%s,\n' % dump_bool(flags & (1 << 4))
        out += 'ReturnsValue:%s,\n' % dump_bool(flags & (1 << 5))
        out += 'NeedsClassClosure:%s,\n' % dump_bool(flags & (1 << 6))

    #out += 'Exec:%s,\n' % dump_bool(st.has_exec()) # Return True if the block uses exec.
    #out += 'ImportStar:%s,\n' % dump_bool(st.has_import_star()) # Return True if the block uses a starred from-import.
    out += 'Varnames:%s,\n' % dump_strings(st._table.varnames)
    out += 'Symbols: Symbols{\n'
    for name in sorted(st.get_identifiers()):
        s = st.lookup(name)
        out += '"%s":%s,\n' % (name, dump_symbol(s))
    out += '},\n'
    out += 'Children:Children{\n'
    for symtable in st.get_children():
        out += '%s,\n' % dump_symtable(symtable)
    out += '},\n'
    out += "}"
    return out

def dump_symbol(s):
    """Dump a symbol"""
    #class symtable.Symbol
    # An entry in a SymbolTable corresponding to an identifier in the source. The constructor is not public.
    out = "Symbol{\n"
    out += 'Flags:%s,\n' % dump_flags(s._Symbol__flags, DEF_FLAGS)
    scope = SCOPES.get(s._Symbol__scope, "scopeUnknown")
    out += 'Scope:%s,\n' % scope
    out += "}"
    return out

def escape(x):
    """Encode strings with backslashes for python/go"""
    return x.replace('\\', "\\\\").replace('"', r'\"').replace("\n", r'\n').replace("\t", r'\t')

def main():
    """Write symtable_data_test.go"""
    path = "symtable_data_test.go"
    out = ["""// Test data generated by make_symtable_test.py - do not edit

package symtable

import (
"github.com/ncw/gpython/py"
)

var symtableTestData = []struct {
in   string
mode string // exec, eval or single
out  *SymTable
exceptionType *py.Type
errString string
}{"""]
    for x in inp:
        source, mode = x[:2]
        if len(x) > 2:
            exc = x[2]
            try:
                table = symtable(source, "<string>", mode)
            except exc as e:
                error = e.msg
            else:
                raise ValueError("Expecting exception %s" % exc)
            dumped_symtable = "nil"
            gostring = "nil"
            exc_name = "py.%s" % exc.__name__
        else:
            table = symtable(source, "<string>", mode)
            exc_name = "nil"
            error = ""
            dumped_symtable = dump_symtable(table)
        out.append('{"%s", "%s", %s, %s, "%s"},' % (escape(source), mode, dumped_symtable, exc_name, escape(error)))
    out.append("}\n")
    print("Writing %s" % path)
    with open(path, "w") as f:
        f.write("\n".join(out))
        f.write("\n")
    subprocess.check_call(["gofmt", "-w", path])

if __name__ == "__main__":
    main()
