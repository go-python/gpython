# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Generate some tests
unops = "- + abs ~ float complex int"
augops = "+ - * / // % & | ^"
binops = augops + " < <= == != > >="
# powmod
# divmod
# lsl
# lsr
# round
# **
# **=

import random

def big():
    return random.randint(-1000000000000000000000000000000, 1000000000000000000000000000000)

def small():
    return random.randint(-1<<63, (1<<63)-1)

def unop(op, a):
    if len(op) == 1:
        expr = "%s%s" % (op, a)
    else:
        expr = "%s(%s)" % (op, a)
    r = eval(expr)
    print("assert (%s) == %s" % (expr, r))

for op in unops.split():
    print("\ndoc='unop %s'" % op)
    for i in range(2):
        a = small()
        unop(op, a)
        a = big()
        unop(op, a)

def binop(op, a, b):
    expr = "%s%s%s" % (a, op, b)
    r = eval(expr)
    print("assert (%s) == %s" % (expr, r))

for op in binops.split():
    print("\ndoc='binop %s'" % op)
    a = small()
    b = small()
    binop(op, a, b)
    a = big()
    b = small()
    binop(op, a, b)
    a = small()
    b = big()
    binop(op, a, b)
    a = big()
    b = big()
    binop(op, a, b)

def augop(op, a, b):
    expr = "%s%s%s" % (a, op, b)
    r = eval(expr)
    print("a = %s\na %s= %s\nassert a == %s" % (a, op, b, r))

for op in augops.split():
    print("\ndoc='augop %s'" % op)
    a = small()
    b = small()
    augop(op, a, b)
    a = big()
    b = small()
    augop(op, a, b)
    a = small()
    b = big()
    augop(op, a, b)
    a = big()
    b = big()
    augop(op, a, b)

