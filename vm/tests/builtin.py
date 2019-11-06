# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="eval"
assert eval("1+2") == 3
glob ={'a':1}
assert eval("a+2", glob) == 3
loc ={'b':2}
assert eval("a+b", glob, loc) == 3
co = compile("a+b+1", "s", "eval")
assert eval(co, glob, loc) == 4
assert eval(b"2+3") == 5

try:
    eval(())
except TypeError as e:
    pass
else:
    assert False, "SyntaxError not raised"
    
try:
    eval("a = 26")
except SyntaxError as e:
    pass
else:
    assert False, "SyntaxError not raised"
    
try:
    eval(1,2,3,4)
except TypeError as e:
    pass
else:
    assert False, "TypeError not raised"
    
try:
    eval("1", object())
except TypeError as e:
    pass
else:
    assert False, "TypeError not raised"

try:
    eval("1", {}, object())
except TypeError as e:
    pass
else:
    assert False, "TypeError not raised"
    
doc="exec"
glob = {"a":100}
assert exec("b = a+100", glob) == None
assert glob["b"] == 200
loc = {"c":23}
assert exec("d = a+b+c", glob, loc) == None
assert loc["d"] == 323
co = compile("d = a+b+c+1", "s", "exec")
assert eval(co, glob, loc) == None
assert loc["d"] == 324

try:
    exec("if")
except SyntaxError as e:
    pass
else:
    assert False, "SyntaxError not raised"

doc="finished"
