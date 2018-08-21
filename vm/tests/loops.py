# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="While"
a = 1
while a < 10:
    a += 1
assert a == 10

doc="While else"
a = 1
ok = False
while a < 10:
    a += 1
else:
    ok = True
assert a == 10
assert ok

doc="While break"
a = 1
ok = True
while True:
    if a >= 10:
        break
    a += 1
else:
    ok = False
assert a == 10
assert ok

doc="While continue"
a = 1
while a < 10:
    if a == 5:
        a += 1000
        continue
    a += 1
assert a == 1005

doc="For"
a = 0
for i in (1,2,3,4,5):
    a += i
assert a == 15

doc="For else"
a = 0
ok = False
for i in (1,2,3,4,5):
    a += i
else:
    ok = True
assert a == 15
assert ok

doc="For break"
a = 0
ok = True
for i in (1,2,3,4,5):
    if i >= 3:
        break
    a += i
else:
    ok = False
assert a == 3
assert ok

doc="For continue"
a = 0
for i in (1,2,3,4,5):
    if i == 3:
        continue
    a += i
assert a == 12

doc="For continue in try/finally"
ok = False
a = 0
for i in (1,2,3,4,5):
    if i == 3:
        try:
            continue
        finally:
            ok = True
    a += i
assert a == 12
assert ok

doc="finished"
