# test_builtin.py:BuiltinTest.test_filter()
from libtest import assertRaises

doc="filter"
class T0:
    def __bool__(self):
        return True
class T1:
    def __len__(self):
        return 1
class T2:
    def __bool__(self):
        return False
class T3:
    pass
t0, t1, t2, t3 = T0(), T1(), T2(), T3()
assert list(filter(None, [t0, t1, t2, t3])) == [t0, t1, t3]
assert list(filter(None, [1, [], 2, ''])) == [1, 2]

class T3:
    def __len__(self):
        raise ValueError
t3 = T3()
assertRaises(ValueError, list, filter(None, [t3]))

class Squares:
    def __init__(self, max):
        self.max = max
        self.sofar = []

    def __len__(self): return len(self.sofar)

    def __getitem__(self, i):
        if not 0 <= i < self.max: raise IndexError
        n = len(self.sofar)
        while n <= i:
            self.sofar.append(n*n)
            n += 1
        return self.sofar[i]

assert list(filter(lambda c: 'a' <= c <= 'z', 'Hello World')) == list('elloorld')
assert list(filter(None, [1, 'hello', [], [3], '', None, 9, 0])) == [1, 'hello', [3], 9]
assert list(filter(lambda x: x > 0, [1, -3, 9, 0, 2])) == [1, 9, 2]
assert list(filter(None, Squares(10))) == [1, 4, 9, 16, 25, 36, 49, 64, 81]
assert list(filter(lambda x: x%2, Squares(10))) == [1, 9, 25, 49, 81]
def identity(item):
    return 1
filter(identity, Squares(5))
assertRaises(TypeError, filter)
class BadSeq(object):
    def __getitem__(self, index):
        if index<4:
            return 42
        raise ValueError
assertRaises(ValueError, list, filter(lambda x: x, BadSeq()))
def badfunc():
    pass
assertRaises(TypeError, list, filter(badfunc, range(5)))

# test bltinmodule.c::filtertuple()
assert list(filter(None, (1, 2))) == [1, 2]
assert list(filter(lambda x: x>=3, (1, 2, 3, 4))) == [3, 4]
assertRaises(TypeError, list, filter(42, (1, 2)))

doc="finished"
