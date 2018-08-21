# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Test functions
doc="fn"
def fn():
    return 1
assert fn() == 1

doc="fn1"
def fn1(x):
    return x+1
assert fn1(1) == 2

doc="fn2"
def fn2(x,y=1):
    return x+y
assert fn2(1) == 2
assert fn2(1,3) == 4
assert fn2(1,y=4) == 5

# Closure

# FIXME something wrong with closures over function arguments...
# doc="counter3"
# def counter3(x):
#     def inc():
#         nonlocal x
#         x += 1
#         return x
#     return inc
# fn3 = counter3(1)
# assert fn3() == 2
# assert fn3() == 3

doc="counter4"
def counter4(initial):
    x = initial
    def inc():
        nonlocal x
        x += 1
        return x
    return inc
fn4 = counter4(1)
assert fn4() == 2
assert fn4() == 3

doc="counter5"
def counter5(initial):
    L = [initial]
    def inc():
        L[0] += 1
        return L[0]
    return inc
fn5 = counter5(1)
assert fn5() == 2
assert fn5() == 3


doc="del_deref6"
def del_deref6(initial):
    x = initial
    def inc():
        nonlocal x
        a = x
        del x
        return a+1
    return inc
fn6 = del_deref6(1)
assert fn6() == 2
try:
    fn6()
except NameError as e:
    pass
else:
    assert False, "NameError not raised"

# check you can't delete it twice!

doc="fn7"
def fn7(b):
 c = 1
 def nested(d):
   nonlocal c
   del c
   del c
 return nested

try:
    fn7(1)(2)
except NameError as e:
    pass
else:
    assert False, "NameError not raised"

# globals

doc="fn8"
a = 1
def fn8():
    global a
    assert a == 1
    a = 2
    assert a == 2
fn8()
assert a == 2

doc="fn9"
def fn9():
    global a
    del a
fn9()
try:
    a
except NameError as e:
    pass
else:
    assert False, "NameError not raised"
try:
    fn9()
except NameError as e:
    pass
else:
    assert False, "NameError not raised"

# delete
doc="fn10"
def fn10():
    a = 1
    assert a == 1
    del a
    try:
        a
    except NameError as e:
        pass
    else:
        assert False, "NameError not raised"
    try:
        del a
    except NameError as e:
        pass
    else:
        assert False, "NameError not raised"
fn10()

# annotations
doc="fn11"
def fn11(a:"A") -> "RET":
    return a+1
assert fn11(1) == 2
# FIXME check annotations are in place

#kwargs
doc="fn12"
def fn12(*args,a=2,b=3,**kwargs):
    return (args,a,b,kwargs)
assert fn12() == ((),2,3,{})
assert fn12(1) == ((1,),2,3,{})
assert fn12(1,2) == ((1,2),2,3,{})
assert fn12(1,2,b=7) == ((1,2),2,7,{})
assert fn12(1,2,a=9,b=7) == ((1,2),9,7,{})
assert fn12(1,2,a=9,b=7,c=10) == ((1,2),9,7,{'c':10})
assert fn12(*(1,2),b=7,**{'a':9,'c':10}) == ((1,2),9,7,{'c':10})
assert fn12(*(1,2),**{'a':9,'b':7,'c':10}) == ((1,2),9,7,{'c':10})

doc="fn13"
def fn13(a,b,*args,c=4,d=5,**kwargs):
    return (a,b,args,c,d,kwargs)
assert fn13(0,1) == (0,1,(),4,5,{})
assert fn13(0,1,2) == (0,1,(2,),4,5,{})
assert fn13(0,1,2,3) == (0,1,(2,3),4,5,{})
assert fn13(0,1,2,3,c=6) == (0,1,(2,3),6,5,{})
assert fn13(0,1,2,3,c=6,d=7) == (0,1,(2,3),6,7,{})
assert fn13(0,*(1,2,3),**{'c':6,'d':7,'e':8}) == (0,1,(2,3),6,7,{'e':8})
assert fn13(*(0,1,2,3),d=7,**{'c':6,'e':8}) == (0,1,(2,3),6,7,{'e':8})
assert fn13(*(0,1,2,3),**{'c':6,'d':7,'e':8}) == (0,1,(2,3),6,7,{'e':8})

doc="Calling errors fn14"
def fn14():
    pass
def ck(fn, text, *args, **kwargs):
    try:
        fn(*args, **kwargs)
    except TypeError as e:
        if e.args[0] != text:
            raise
    else:
        assert False, "TypeError not raised"
            
ck(fn14, "fn14() got an unexpected keyword argument 'a'", a=1)
try:
    fn14(a=1,**{'a':2})
except TypeError as e:
        if e.args[0] != "fn14() got multiple values for keyword argument 'a'":
            raise
else:
    assert False, "Type error not raised"
ck(fn14, "fn14() takes 0 positional arguments but 1 was given", 1)

doc="Calling errors fn15"
def fn15_1(a):
    pass
def fn15_2(a,b):
    pass
def fn15_3(a,b,c):
    pass
def fn15_4(a,b,c,d):
    pass
ck(fn15_1, "fn15_1() missing 1 required positional argument: 'a'")
ck(fn15_2, "fn15_2() missing 2 required positional arguments: 'a' and 'b'")
ck(fn15_3, "fn15_3() missing 3 required positional arguments: 'a', 'b', and 'c'")
ck(fn15_4, "fn15_4() missing 4 required positional arguments: 'a', 'b', 'c', and 'd'")
ck(fn15_4, "fn15_4() missing 3 required positional arguments: 'b', 'c', and 'd'", 1)
ck(fn15_4, "fn15_4() missing 3 required positional arguments: 'a', 'b', and 'd'", c=3)

doc="Calling errors fn16"
def fn16_0(a=1):
    pass
def fn16_1(*,a=1):
    pass
def fn16_2(a,*,b=1):
    pass
def fn16_3(a,b,c,*,d=1,e=2,f=3):
    pass
def fn16_4(*,a):
    pass
def fn16_5(*,a,b):
    pass
def fn16_6(*,a,b,c):
    pass
fn16_0()
ck(fn16_1, "fn16_1() takes 0 positional arguments but 1 was given", 1)
ck(fn16_2, "fn16_2() missing 1 required positional argument: 'a'")
ck(fn16_2, "fn16_2() takes 1 positional argument but 2 were given", 1, 2)
ck(fn16_3, "fn16_3() takes 3 positional arguments but 4 were given", 1, 2, 3, 4)
ck(fn16_4, "fn16_4() missing 1 required keyword-only argument: 'a'")
ck(fn16_5, "fn16_5() missing 2 required keyword-only arguments: 'a' and 'b'")
ck(fn16_6, "fn16_6() missing 3 required keyword-only arguments: 'a', 'b', and 'c'")

#FIXME decorators

doc="finished"
