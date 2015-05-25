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
# FIXME not implemented assert fn2(1,y=4) == 5

# FIXME check *arg and **kwarg


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
def fn12(*args,a=2,b=3,**kwargs) -> "RET":
    return args
# FIXME this blows up: assert fn12() == ()
# FIXME check kwargs passing


#FIXME decorators

doc="finished"
