#!/usr/bin/env python3.4

doc="Test class definitions"
class C1:
    "Test 1"
    def method1(self, x):
        "method1"
        return x+1
    def method2(self, m2):
        "method2"
        return self.method1(m2)+m2

c = C1()
assert c.method1(1) == 2
assert c.method2(1) == 3

doc="Test class definitions 2"
class C2:
    "Test 2"
    _VAR = 1
    VAR = _VAR + 1
    def method1(self, x):
        "method1"
        return self.VAR + x
    def method2(self, m2):
        "method2"
        return self.method1(m2)+m2

c = C2()
assert c.method1(1) == 3
assert c.method2(1) == 4

doc="CLASS_DEREF"

# FIXME corner cases in CLASS_DEREF
def classderef(y):
    # FIXME should work on parameter of classderef - y - but doesn't
    x = y
    class DeRefTest:
        VAR = x
        def method1(self, x):
            "method1"
            return self.VAR+x
    return DeRefTest
x = classderef(1)
c = x()
assert c.method1(1) == 2

# End with this
doc="finished"
