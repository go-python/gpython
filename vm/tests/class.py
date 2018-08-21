# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

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

# FIXME more corner cases in CLASS_DEREF

doc="CLASS_DEREF"
def classderef(y):
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

# FIXME doesn't work
# doc="CLASS_DEREF2"
# def classderef2(x):
#     class DeRefTest:
#         VAR = x
#         def method1(self, x):
#             "method1"
#             return self.VAR+x
#     return DeRefTest
# x = classderef2(1)
# c = x()
# assert c.method1(1) == 2

doc="finished"
