# test_builtin.py:BuiltinTest.test_map()
from libtest import assertRaises

doc="map"
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

assert list(map(lambda x: x*x, range(1,4))) == [1, 4, 9]
try:
    from math import sqrt
except ImportError:
    def sqrt(x):
        return pow(x, 0.5)
assert list(map(lambda x: list(map(sqrt, x)), [[16, 4], [81, 9]])) == [[4.0, 2.0], [9.0, 3.0]]
assert list(map(lambda x, y: x+y, [1,3,2], [9,1,4])) == [10, 4, 6]

def plus(*v):
    accu = 0
    for i in v: accu = accu + i
    return accu
assert list(map(plus, [1, 3, 7])) == [1, 3, 7]
assert list(map(plus, [1, 3, 7], [4, 9, 2])) == [1+4, 3+9, 7+2]
assert list(map(plus, [1, 3, 7], [4, 9, 2], [1, 1, 0])) == [1+4+1, 3+9+1, 7+2+0]
assert list(map(int, Squares(10))) == [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
def Max(a, b):
    if a is None:
        return b
    if b is None:
        return a
    return max(a, b)
assert list(map(Max, Squares(3), Squares(2))) == [0, 1]
assertRaises(TypeError, map)
assertRaises(TypeError, map, lambda x: x, 42)
class BadSeq:
    def __iter__(self):
        raise ValueError
        yield None
assertRaises(ValueError, list, map(lambda x: x, BadSeq()))
def badfunc(x):
    raise RuntimeError
assertRaises(RuntimeError, list, map(badfunc, range(5)))
doc="finished"
