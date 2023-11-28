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
        # FIXME(sbinet): implement
        print("  SKIP: NotImplemented")
        continue
    if typ in "bhilqfd":
        arr = array.array(typ, [-1, -2, -3, -4])
    if typ in "BHILQ":
        arr = array.array(typ, [+1, +2, +3, +4])
    print("  array: %s" % (repr(arr),))
    print("  itemsize: %s" % (arr.itemsize,))
    print("  typecode: %s" % (arr.typecode,))
    print("  len:      %s" % (len(arr),))
    print("  arr[0]: %s" % (arr[0],))
    print("  arr[-1]: %s" % (arr[-1],))
    try:
        arr[-10]
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")

    try:
        arr[10]
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")
    arr[-2] = 33
    print("  arr[-2]: %s" % (arr[-2],))

    try:
        arr[-10] = 2
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")

    if typ in "bhilqfd":
        arr.extend([-5,-6])
    if typ in "BHILQ":
        arr.extend([5,6])
    print("  array: %s" % (repr(arr),))
    print("  len:   %s" % (len(arr),))

    if typ in "bhilqfd":
        arr.append(-7)
    if typ in "BHILQ":
        arr.append(7)
    print("  array: %s" % (repr(arr),))
    print("  len:   %s" % (len(arr),))

    try:
        arr.append()
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.append([])
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.append(1, 2)
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.append(None)
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")

    try:
        arr.extend()
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.extend(None)
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")
    try:
        arr.extend([1,None])
        print("  ERROR: expected an exception")
    except:
        print("  caught an exception [ok]")
    pass

print("\n")
print("## testing array.array(...)")
try:
    arr = array.array()
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array(b"d")
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("?")
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("dd")
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("d", initializer=[1,2])
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("d", [1], [])
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("d", 1)
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")

try:
    arr = array.array("d", ["a","b"])
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")

try:
    ## FIXME(sbinet): implement it at some point.
    arr = array.array("u")
    print("ERROR: expected an exception")
except:
    print("caught an exception [ok]")
