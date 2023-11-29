# Copyright 2023 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import array

print("globals:")
for name in ("typecodes", "array"):
    v = getattr(array, name)
    print("\narray.%s:\n%s" % (name,repr(v)))
    pass

def assertEqual(x, y):
    assert x == y, "got: %s, want: %s" % (repr(x), repr(y))

assertEqual(array.typecodes, 'bBuhHiIlLqQfd')

for i, typ in enumerate(array.typecodes):
    print("")
    print("typecode '%s'" % (typ,))
    if typ == 'u':
        arr = array.array(typ, "?世界!")
    if typ in "bhilq":
        arr = array.array(typ, [-1, -2, -3, -4])
    if typ in "BHILQ":
        arr = array.array(typ, [+1, +2, +3, +4])
    if typ in "fd":
        arr = array.array(typ, [-1.0, -2.0, -3.0, -4.0])
    print("  array: %s ## repr" % (repr(arr),))
    print("  array: %s ## str" % (str(arr),))
    print("  itemsize: %s" % (arr.itemsize,))
    print("  typecode: %s" % (arr.typecode,))
    print("  len:      %s" % (len(arr),))
    print("  arr[0]: %s" % (arr[0],))
    print("  arr[-1]: %s" % (arr[-1],))
    try:
        arr[-10]
        print("  ERROR1: expected an exception")
    except:
        print("  caught an exception [ok]")

    try:
        arr[10]
        print("  ERROR2: expected an exception")
    except:
        print("  caught an exception [ok]")
    arr[-2] = 33
    if typ in "fd":
        arr[-2] = 0.3
    print("  arr[-2]: %s" % (arr[-2],))

    try:
        arr[-10] = 2
        print("  ERROR3: expected an exception")
    except:
        print("  caught an exception [ok]")

    if typ in "bhilqfd":
        arr.extend([-5,-6])
    if typ in "BHILQ":
        arr.extend([5,6])
    if typ == 'u':
        arr.extend("he")
    print("  array: %s" % (repr(arr),))
    print("  len:   %s" % (len(arr),))

    if typ in "bhilqfd":
        arr.append(-7)
    if typ in "BHILQ":
        arr.append(7)
    if typ == 'u':
        arr.append("l")
    print("  array: %s" % (repr(arr),))
    print("  len:   %s" % (len(arr),))

    try:
        arr.append()
        print("  ERROR4: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.append([])
        print("  ERROR5: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.append(1, 2)
        print("  ERROR6: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.append(None)
        print("  ERROR7: expected an exception")
    except:
        print("  caught an exception [ok]")

    try:
        arr.extend()
        print("  ERROR8: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.extend(None)
        print("  ERROR9: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.extend([1,None])
        print("  ERROR10: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.extend(1,None)
        print("  ERROR11: expected an exception")
    except:
        print("  caught an exception [ok]")

    try:
        arr[0] = object()
        print("  ERROR12: expected an exception")
    except:
        print("  caught an exception [ok]")
    pass

print("\n")
print("## testing array.array(...)")
try:
    arr = array.array()
    print("ERROR1: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array(b"d")
    print("ERROR2: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("?")
    print("ERROR3: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("dd")
    print("ERROR4: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("d", initializer=[1,2])
    print("ERROR5: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("d", [1], [])
    print("ERROR6: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("d", 1)
    print("ERROR7: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("d", ["a","b"])
    print("ERROR8: expected an exception")
except:
    print("caught an exception [ok]")
