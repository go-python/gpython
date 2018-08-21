# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Converted from cpython/Lib/test/test_math.py
# FIXME ignoring ./Lib/test/math_testcases.txt, /Lib/test/cmath_testcases.txt for the moment
#
# Python test set -- math module
# XXXX Should not do tests around zero only

import math
from libtest import *
from libulp import *

eps = 1E-15
NAN = float('nan')
INF = float('inf')
NINF = float('-inf')

# detect evidence of double-rounding: fsum is not always correctly
# rounded on machines that suffer from double rounding.
x, y = 1e16, 2.9999 # use temporary values to defeat peephole optimizer
HAVE_DOUBLE_ROUNDING = (x + y == 1e16 + 4)

def ftest(name, got, want, eps=None):
    # Use %r instead of %f so the error message
    # displays full precision. Otherwise discrepancies
    # in the last few bits will lead to very confusing
    # error messages
    what = '%s got %r, want %r' % (name, got, want)
    if eps is None:
        ulps_check(what, want, got)
    else:
        assert abs(got-want) <= eps, what

doc="to_ulps"
assert to_ulps(                       0) ==  0x0000000000000000
assert to_ulps(                       0) ==  0x0000000000000000
assert to_ulps(                     100) ==  0x4059000000000000
assert to_ulps(             100.0000001) ==  0x40590000006b5fca
assert to_ulps(                    -100) == -0x4059000000000000
assert to_ulps(            -100.0000001) == -0x40590000006b5fca
assert to_ulps(                  1e+307) ==  0x7fac7b1f3cac7433
assert to_ulps(                 -1e+307) == -0x7fac7b1f3cac7433
assert to_ulps(                1.5e+308) ==  0x7feab36d48e1acf0
assert to_ulps(               -1.5e+308) == -0x7feab36d48e1acf0
# FIXME(different in py3) assert to_ulps(                     NAN) ==  0x7ff8000000000001
assert to_ulps(                     INF) ==  0x7ff0000000000000
assert to_ulps(                    NINF) == -0x7ff0000000000000
assert to_ulps(                  5e-324) ==  0x0000000000000001
assert to_ulps(                 -5e-324) == -0x0000000000000001
assert to_ulps( 1.7976931348623157e+308) ==  0x7fefffffffffffff
assert to_ulps(-1.7976931348623157e+308) == -0x7fefffffffffffff

doc="Constants"
ftest('pi', math.pi, 3.141592653589793)
ftest('e', math.e,   2.718281828459045)

doc="Acos"
assertRaises(TypeError, math.acos)
ftest('acos(-1)', math.acos(-1), math.pi)
ftest('acos(0)', math.acos(0), math.pi/2)
ftest('acos(1)', math.acos(1), 0)
assertRaises(ValueError, math.acos, INF)
assertRaises(ValueError, math.acos, NINF)
assertTrue(math.isnan(math.acos(NAN)))

doc="Acosh"
assertRaises(TypeError, math.acosh)
ftest('acosh(1)', math.acosh(1), 0)
ftest('acosh(2)', math.acosh(2), 1.3169578969248168)
assertRaises(ValueError, math.acosh, 0)
assertRaises(ValueError, math.acosh, -1)
assertEqual(math.acosh(INF), INF)
assertRaises(ValueError, math.acosh, NINF)
assertTrue(math.isnan(math.acosh(NAN)))

doc="Asin"
assertRaises(TypeError, math.asin)
ftest('asin(-1)', math.asin(-1), -math.pi/2)
ftest('asin(0)', math.asin(0), 0)
ftest('asin(1)', math.asin(1), math.pi/2)
assertRaises(ValueError, math.asin, INF)
assertRaises(ValueError, math.asin, NINF)
assertTrue(math.isnan(math.asin(NAN)))

doc="Asinh"
assertRaises(TypeError, math.asinh)
ftest('asinh(0)', math.asinh(0), 0)
ftest('asinh(1)', math.asinh(1), 0.88137358701954305)
ftest('asinh(-1)', math.asinh(-1), -0.88137358701954305)
assertEqual(math.asinh(INF), INF)
assertEqual(math.asinh(NINF), NINF)
assertTrue(math.isnan(math.asinh(NAN)))

doc="Atan"
assertRaises(TypeError, math.atan)
ftest('atan(-1)', math.atan(-1), -math.pi/4)
ftest('atan(0)', math.atan(0), 0)
ftest('atan(1)', math.atan(1), math.pi/4)
ftest('atan(inf)', math.atan(INF), math.pi/2)
ftest('atan(-inf)', math.atan(NINF), -math.pi/2)
assertTrue(math.isnan(math.atan(NAN)))

doc="Atanh"
assertRaises(TypeError, math.atan)
ftest('atanh(0)', math.atanh(0), 0)
ftest('atanh(0.5)', math.atanh(0.5), 0.54930614433405489)
ftest('atanh(-0.5)', math.atanh(-0.5), -0.54930614433405489)
assertRaises(ValueError, math.atanh, 1)
assertRaises(ValueError, math.atanh, -1)
assertRaises(ValueError, math.atanh, INF)
assertRaises(ValueError, math.atanh, NINF)
assertTrue(math.isnan(math.atanh(NAN)))

doc="Atan2"
assertRaises(TypeError, math.atan2)
ftest('atan2(-1, 0)', math.atan2(-1, 0), -math.pi/2)
ftest('atan2(-1, 1)', math.atan2(-1, 1), -math.pi/4)
ftest('atan2(0, 1)', math.atan2(0, 1), 0)
ftest('atan2(1, 1)', math.atan2(1, 1), math.pi/4)
ftest('atan2(1, 0)', math.atan2(1, 0), math.pi/2)

# math.atan2(0, x)
ftest('atan2(0., -inf)', math.atan2(0., NINF), math.pi)
ftest('atan2(0., -2.3)', math.atan2(0., -2.3), math.pi)
ftest('atan2(0., -0.)', math.atan2(0., -0.), math.pi)
assertEqual(math.atan2(0., 0.), 0.)
assertEqual(math.atan2(0., 2.3), 0.)
assertEqual(math.atan2(0., INF), 0.)
assertTrue(math.isnan(math.atan2(0., NAN)))
# math.atan2(-0, x)
ftest('atan2(-0., -inf)', math.atan2(-0., NINF), -math.pi)
ftest('atan2(-0., -2.3)', math.atan2(-0., -2.3), -math.pi)
ftest('atan2(-0., -0.)', math.atan2(-0., -0.), -math.pi)
assertEqual(math.atan2(-0., 0.), -0.)
assertEqual(math.atan2(-0., 2.3), -0.)
assertEqual(math.atan2(-0., INF), -0.)
assertTrue(math.isnan(math.atan2(-0., NAN)))
# math.atan2(INF, x)
ftest('atan2(inf, -inf)', math.atan2(INF, NINF), math.pi*3/4)
ftest('atan2(inf, -2.3)', math.atan2(INF, -2.3), math.pi/2)
ftest('atan2(inf, -0.)', math.atan2(INF, -0.0), math.pi/2)
ftest('atan2(inf, 0.)', math.atan2(INF, 0.0), math.pi/2)
ftest('atan2(inf, 2.3)', math.atan2(INF, 2.3), math.pi/2)
ftest('atan2(inf, inf)', math.atan2(INF, INF), math.pi/4)
assertTrue(math.isnan(math.atan2(INF, NAN)))
# math.atan2(NINF, x)
ftest('atan2(-inf, -inf)', math.atan2(NINF, NINF), -math.pi*3/4)
ftest('atan2(-inf, -2.3)', math.atan2(NINF, -2.3), -math.pi/2)
ftest('atan2(-inf, -0.)', math.atan2(NINF, -0.0), -math.pi/2)
ftest('atan2(-inf, 0.)', math.atan2(NINF, 0.0), -math.pi/2)
ftest('atan2(-inf, 2.3)', math.atan2(NINF, 2.3), -math.pi/2)
ftest('atan2(-inf, inf)', math.atan2(NINF, INF), -math.pi/4)
assertTrue(math.isnan(math.atan2(NINF, NAN)))
# math.atan2(+finite, x)
ftest('atan2(2.3, -inf)', math.atan2(2.3, NINF), math.pi)
ftest('atan2(2.3, -0.)', math.atan2(2.3, -0.), math.pi/2)
ftest('atan2(2.3, 0.)', math.atan2(2.3, 0.), math.pi/2)
assertEqual(math.atan2(2.3, INF), 0.)
assertTrue(math.isnan(math.atan2(2.3, NAN)))
# math.atan2(-finite, x)
ftest('atan2(-2.3, -inf)', math.atan2(-2.3, NINF), -math.pi)
ftest('atan2(-2.3, -0.)', math.atan2(-2.3, -0.), -math.pi/2)
ftest('atan2(-2.3, 0.)', math.atan2(-2.3, 0.), -math.pi/2)
assertEqual(math.atan2(-2.3, INF), -0.)
assertTrue(math.isnan(math.atan2(-2.3, NAN)))
# math.atan2(NAN, x)
assertTrue(math.isnan(math.atan2(NAN, NINF)))
assertTrue(math.isnan(math.atan2(NAN, -2.3)))
assertTrue(math.isnan(math.atan2(NAN, -0.)))
assertTrue(math.isnan(math.atan2(NAN, 0.)))
assertTrue(math.isnan(math.atan2(NAN, 2.3)))
assertTrue(math.isnan(math.atan2(NAN, INF)))
assertTrue(math.isnan(math.atan2(NAN, NAN)))

doc="Ceil"
assertRaises(TypeError, math.ceil)
assertEqual(int, type(math.ceil(0.5)))
ftest('ceil(0.5)', math.ceil(0.5), 1)
ftest('ceil(1.0)', math.ceil(1.0), 1)
ftest('ceil(1.5)', math.ceil(1.5), 2)
ftest('ceil(-0.5)', math.ceil(-0.5), 0)
ftest('ceil(-1.0)', math.ceil(-1.0), -1)
ftest('ceil(-1.5)', math.ceil(-1.5), -1)
#assertEqual(math.ceil(INF), INF)
#assertEqual(math.ceil(NINF), NINF)
#assertTrue(math.isnan(math.ceil(NAN)))

class TestCeil:
    def __ceil__(self):
        return 42
class TestNoCeil:
    pass
ftest('ceil(TestCeil())', math.ceil(TestCeil()), 42)
assertRaises(TypeError, math.ceil, TestNoCeil())

t = TestNoCeil()
t.__ceil__ = lambda *args: args
# FIXME assertRaises(TypeError, math.ceil, t)
assertRaises(TypeError, math.ceil, t, 0)

doc="Copysign"
assertEqual(math.copysign(1, 42), 1.0)
assertEqual(math.copysign(0., 42), 0.0)
assertEqual(math.copysign(1., -42), -1.0)
assertEqual(math.copysign(3, 0.), 3.0)
assertEqual(math.copysign(4., -0.), -4.0)

assertRaises(TypeError, math.copysign)
# copysign should let us distinguish signs of zeros
assertEqual(math.copysign(1., 0.), 1.)
assertEqual(math.copysign(1., -0.), -1.)
assertEqual(math.copysign(INF, 0.), INF)
assertEqual(math.copysign(INF, -0.), NINF)
assertEqual(math.copysign(NINF, 0.), INF)
assertEqual(math.copysign(NINF, -0.), NINF)
# and of infinities
assertEqual(math.copysign(1., INF), 1.)
assertEqual(math.copysign(1., NINF), -1.)
assertEqual(math.copysign(INF, INF), INF)
assertEqual(math.copysign(INF, NINF), NINF)
assertEqual(math.copysign(NINF, INF), INF)
assertEqual(math.copysign(NINF, NINF), NINF)
assertTrue(math.isnan(math.copysign(NAN, 1.)))
assertTrue(math.isnan(math.copysign(NAN, INF)))
assertTrue(math.isnan(math.copysign(NAN, NINF)))
assertTrue(math.isnan(math.copysign(NAN, NAN)))
# copysign(INF, NAN) may be INF or it may be NINF, since
# we don't know whether the sign bit of NAN is set on any
# given platform.
assertTrue(math.isinf(math.copysign(INF, NAN)))
# similarly, copysign(2., NAN) could be 2. or -2.
assertEqual(abs(math.copysign(2., NAN)), 2.)

doc="Cos"
assertRaises(TypeError, math.cos)
ftest('cos(-pi/2)', math.cos(-math.pi/2), 0, eps=eps)
ftest('cos(0)', math.cos(0), 1)
ftest('cos(pi/2)', math.cos(math.pi/2), 0, eps=eps)
ftest('cos(pi)', math.cos(math.pi), -1)
try:
    assertTrue(math.isnan(math.cos(INF)))
    assertTrue(math.isnan(math.cos(NINF)))
except ValueError:
    assertRaises(ValueError, math.cos, INF)
    assertRaises(ValueError, math.cos, NINF)
assertTrue(math.isnan(math.cos(NAN)))

doc="Cosh"
assertRaises(TypeError, math.cosh)
ftest('cosh(0)', math.cosh(0), 1)
ftest('cosh(2)-2*cosh(1)**2', math.cosh(2)-2*math.cosh(1)**2, -1) # Thanks to Lambert
assertEqual(math.cosh(INF), INF)
assertEqual(math.cosh(NINF), INF)
assertTrue(math.isnan(math.cosh(NAN)))

doc="Degrees"
assertRaises(TypeError, math.degrees)
ftest('degrees(pi)', math.degrees(math.pi), 180.0)
ftest('degrees(pi/2)', math.degrees(math.pi/2), 90.0)
ftest('degrees(-pi/4)', math.degrees(-math.pi/4), -45.0)

doc="Exp"
assertRaises(TypeError, math.exp)
ftest('exp(-1)', math.exp(-1), 1/math.e)
ftest('exp(0)', math.exp(0), 1)
ftest('exp(1)', math.exp(1), math.e)
assertEqual(math.exp(INF), INF)
assertEqual(math.exp(NINF), 0.)
assertTrue(math.isnan(math.exp(NAN)))

doc="Fabs"
assertRaises(TypeError, math.fabs)
ftest('fabs(-1)', math.fabs(-1), 1)
ftest('fabs(0)', math.fabs(0), 0)
ftest('fabs(1)', math.fabs(1), 1)

doc="Factorial"
assertEqual(math.factorial(0), 1)
assertEqual(math.factorial(0.0), 1)
total = 1
for i in range(1, 1000):
    total *= i
    # print("total", str(total))
    # print("fact ", str(math.factorial(i)))
    assertEqual(math.factorial(i), total)
    assertEqual(math.factorial(float(i)), total)
# assertRaises(ValueError, math.factorial, -1)
# assertRaises(ValueError, math.factorial, -1.0)
# assertRaises(ValueError, math.factorial, math.pi)
# assertRaises(OverflowError, math.factorial, 1<<63)
# assertRaises(OverflowError, math.factorial, 10e100)

doc="Floor"
assertRaises(TypeError, math.floor)
assertEqual(int, type(math.floor(0.5)))
ftest('floor(0.5)', math.floor(0.5), 0)
ftest('floor(1.0)', math.floor(1.0), 1)
ftest('floor(1.5)', math.floor(1.5), 1)
ftest('floor(-0.5)', math.floor(-0.5), -1)
ftest('floor(-1.0)', math.floor(-1.0), -1)
ftest('floor(-1.5)', math.floor(-1.5), -2)
# pow() relies on floor() to check for integers
# This fails on some platforms - so check it here
ftest('floor(1.23e167)', math.floor(1.23e167), 1.23e167)
ftest('floor(-1.23e167)', math.floor(-1.23e167), -1.23e167)
#assertEqual(math.ceil(INF), INF)
#assertEqual(math.ceil(NINF), NINF)
#assertTrue(math.isnan(math.floor(NAN)))

class TestFloor:
    def __floor__(self):
        return 42
class TestNoFloor:
    pass
ftest('floor(TestFloor())', math.floor(TestFloor()), 42)
assertRaises(TypeError, math.floor, TestNoFloor())

t = TestNoFloor()
t.__floor__ = lambda *args: args
# FIXME assertRaises(TypeError, math.floor, t)
assertRaises(TypeError, math.floor, t, 0)

doc="Fmod"
assertRaises(TypeError, math.fmod)
ftest('fmod(10, 1)', math.fmod(10, 1), 0.0)
ftest('fmod(10, 0.5)', math.fmod(10, 0.5), 0.0)
ftest('fmod(10, 1.5)', math.fmod(10, 1.5), 1.0)
ftest('fmod(-10, 1)', math.fmod(-10, 1), -0.0)
ftest('fmod(-10, 0.5)', math.fmod(-10, 0.5), -0.0)
ftest('fmod(-10, 1.5)', math.fmod(-10, 1.5), -1.0)
assertTrue(math.isnan(math.fmod(NAN, 1.)))
assertTrue(math.isnan(math.fmod(1., NAN)))
assertTrue(math.isnan(math.fmod(NAN, NAN)))
assertRaises(ValueError, math.fmod, 1., 0.)
assertRaises(ValueError, math.fmod, INF, 1.)
assertRaises(ValueError, math.fmod, NINF, 1.)
assertRaises(ValueError, math.fmod, INF, 0.)
assertEqual(math.fmod(3.0, INF), 3.0)
assertEqual(math.fmod(-3.0, INF), -3.0)
assertEqual(math.fmod(3.0, NINF), 3.0)
assertEqual(math.fmod(-3.0, NINF), -3.0)
assertEqual(math.fmod(0.0, 3.0), 0.0)
assertEqual(math.fmod(0.0, NINF), 0.0)

doc="Frexp"
assertRaises(TypeError, math.frexp)

def testfrexp(name, result, expected):
    (mant, exp), (emant, eexp) = result, expected
    if abs(mant-emant) > eps or exp != eexp:
        fail('%s returned %r, expected %r'%\
                  (name, result, expected))

testfrexp('frexp(-1)', math.frexp(-1), (-0.5, 1))
testfrexp('frexp(0)', math.frexp(0), (0, 0))
testfrexp('frexp(1)', math.frexp(1), (0.5, 1))
testfrexp('frexp(2)', math.frexp(2), (0.5, 2))

assertEqual(math.frexp(INF)[0], INF)
assertEqual(math.frexp(NINF)[0], NINF)
assertTrue(math.isnan(math.frexp(NAN)[0]))

doc="Fsum"
# math.fsum relies on exact rounding for correct operation.
# There's a known problem with IA32 floating-point that causes
# inexact rounding in some situations, and will cause the
# math.fsum tests below to fail; see issue #2937.  On non IEEE
# 754 platforms, and on IEEE 754 platforms that exhibit the
# problem described in issue #2937, we simply skip the whole
# test.

# Python version of math.fsum, for comparison.  Uses a
# different algorithm based on frexp, ldexp and integer
# arithmetic.

#from sys import float_info
mant_dig = 53 # FIXME float_info.mant_dig
etiny = -1074 # FIXME float_info.min_exp - mant_dig

def msum(iterable):
    """Full precision summation.  Compute sum(iterable) without any
    intermediate accumulation of error.  Based on the 'lsum' function
    at http://code.activestate.com/recipes/393090/

    """
    tmant, texp = 0, 0
    for x in iterable:
        mant, exp = math.frexp(x)
        mant, exp = int(math.ldexp(mant, mant_dig)), exp - mant_dig
        if texp > exp:
            tmant <<= texp-exp
            texp = exp
        else:
            mant <<= exp-texp
        tmant += mant
    # Round tmant * 2**texp to a float.  The original recipe
    # used float(str(tmant)) * 2.0**texp for this, but that's
    # a little unsafe because str -> float conversion can't be
    # relied upon to do correct rounding on all platforms.
    tail = max(len(bin(abs(tmant)))-2 - mant_dig, etiny - texp)
    if tail > 0:
        h = 1 << (tail-1)
        tmant = tmant // (2*h) + bool(tmant & h and tmant & 3*h-1)
        texp += tail
    return math.ldexp(tmant, texp)

test_values = [
    ([], 0.0),
    ([0.0], 0.0),
    ([1e100, 1.0, -1e100, 1e-100, 1e50, -1.0, -1e50], 1e-100),
    ([2.0**53, -0.5, -2.0**-54], 2.0**53-1.0),
    ([2.0**53, 1.0, 2.0**-100], 2.0**53+2.0),
    ([2.0**53+10.0, 1.0, 2.0**-100], 2.0**53+12.0),
    ([2.0**53-4.0, 0.5, 2.0**-54], 2.0**53-3.0),
    ([1./n for n in range(1, 1001)],
     # float.fromhex('0x1.df11f45f4e61ap+2')),
     7.4854708605503450513651842),
    ([(-1.)**n/n for n in range(1, 1001)],
     # float.fromhex('-0x1.62a2af1bd3624p-1')),
     -0.6926474305598202541034425),
    ([1.7**(i+1)-1.7**i for i in range(1000)] + [-1.7**1000], -1.0),
    ([1e16, 1., 1e-16], 10000000000000002.0),
    ([1e16-2., 1.-2.**-53, -(1e16-2.), -(1.-2.**-53)], 0.0),
    # exercise code for resizing partials array
    ([2.**n - 2.**(n+50) + 2.**(n+52) for n in range(-1074, 972, 2)] +
     [-2.**1022],
     # float.fromhex('0x1.5555555555555p+970')),
    1.330560206356479800576653e+292),
    ]

i = 0
for vals, expected in test_values:
    try:
        actual = math.fsum(vals)
    except OverflowError:
        fail("test %d failed: got OverflowError, expected %r "
                  "for math.fsum(%.100r)" % (i, expected, vals))
    except ValueError:
        fail("test %d failed: got ValueError, expected %r "
                  "for math.fsum(%.100r)" % (i, expected, vals))
    #print("want %g got %g" % (actual, expected))
    assertEqual(actual, expected)
    i+=1

# FIXME
# from random import random, gauss, shuffle
# for j in range(1000):
#     vals = [7, 1e100, -7, -1e100, -9e-20, 8e-20] * 10
#     s = 0
#     for i in range(200):
#         v = gauss(0, random()) ** 7 - s
#         s += v
#         vals.append(v)
#     shuffle(vals)
#
#     s = msum(vals)
#     assertEqual(msum(vals), math.fsum(vals))

doc="Hypot"
assertRaises(TypeError, math.hypot)
ftest('hypot(0,0)', math.hypot(0,0), 0)
ftest('hypot(3,4)', math.hypot(3,4), 5)
assertEqual(math.hypot(NAN, INF), INF)
assertEqual(math.hypot(INF, NAN), INF)
assertEqual(math.hypot(NAN, NINF), INF)
assertEqual(math.hypot(NINF, NAN), INF)
assertTrue(math.isnan(math.hypot(1.0, NAN)))
assertTrue(math.isnan(math.hypot(NAN, -2.0)))

doc="Ldexp"
assertRaises(TypeError, math.ldexp)
ftest('ldexp(0,1)', math.ldexp(0,1), 0)
ftest('ldexp(1,1)', math.ldexp(1,1), 2)
ftest('ldexp(1,-1)', math.ldexp(1,-1), 0.5)
ftest('ldexp(-1,1)', math.ldexp(-1,1), -2)
assertRaises(OverflowError, math.ldexp, 1., 1000000)
assertRaises(OverflowError, math.ldexp, -1., 1000000)
assertEqual(math.ldexp(1., -1000000), 0.)
assertEqual(math.ldexp(-1., -1000000), -0.)
assertEqual(math.ldexp(INF, 30), INF)
assertEqual(math.ldexp(NINF, -213), NINF)
assertTrue(math.isnan(math.ldexp(NAN, 0)))

# large second argument
for n in [10**5, 10**10, 10**20, 10**40]:
    assertEqual(math.ldexp(INF, -n), INF)
    assertEqual(math.ldexp(NINF, -n), NINF)
    assertEqual(math.ldexp(1., -n), 0.)
    assertEqual(math.ldexp(-1., -n), -0.)
    assertEqual(math.ldexp(0., -n), 0.)
    assertEqual(math.ldexp(-0., -n), -0.)
    assertTrue(math.isnan(math.ldexp(NAN, -n)))

    assertRaises(OverflowError, math.ldexp, 1., n)
    assertRaises(OverflowError, math.ldexp, -1., n)
    assertEqual(math.ldexp(0., n), 0.)
    assertEqual(math.ldexp(-0., n), -0.)
    assertEqual(math.ldexp(INF, n), INF)
    assertEqual(math.ldexp(NINF, n), NINF)
    assertTrue(math.isnan(math.ldexp(NAN, n)))

doc="Log"
assertRaises(TypeError, math.log)
ftest('log(1/e)', math.log(1/math.e), -1)
ftest('log(1)', math.log(1), 0)
ftest('log(e)', math.log(math.e), 1)
ftest('log(32,2)', math.log(32,2), 5)
ftest('log(10**40, 10)', math.log(10**40, 10), 40)
ftest('log(10**40, 10**20)', math.log(10**40, 10**20), 2)
ftest('log(10**1000)', math.log(10**1000),
           2302.5850929940457)
assertRaises(ValueError, math.log, -1.5)
assertRaises(ValueError, math.log, -10**1000)
assertRaises(ValueError, math.log, NINF)
assertEqual(math.log(INF), INF)
assertTrue(math.isnan(math.log(NAN)))

doc="Log1p"
assertRaises(TypeError, math.log1p)
n= 2**90
assertAlmostEqual(math.log1p(n), math.log1p(float(n)))

doc="Log2"
assertRaises(TypeError, math.log2)

# Check some integer values
assertEqual(math.log2(1), 0.0)
assertEqual(math.log2(2), 1.0)
assertEqual(math.log2(4), 2.0)

# Large integer values
assertEqual(math.log2(2**1023), 1023.0)
assertEqual(math.log2(2**1024), 1024.0)
assertEqual(math.log2(2**2000), 2000.0)

assertRaises(ValueError, math.log2, -1.5)
assertRaises(ValueError, math.log2, NINF)
assertTrue(math.isnan(math.log2(NAN)))

doc="Log2Exact"
# Check that we get exact equality for log2 of powers of 2.
actual = [math.log2(math.ldexp(1.0, n)) for n in range(-1074, 1024)]
expected = [float(n) for n in range(-1074, 1024)]
assertEqual(actual, expected)

doc="Log10"
assertRaises(TypeError, math.log10)
ftest('log10(0.1)', math.log10(0.1), -1)
ftest('log10(1)', math.log10(1), 0)
ftest('log10(10)', math.log10(10), 1)
ftest('log10(10**1000)', math.log10(10**1000), 1000.0)
assertRaises(ValueError, math.log10, -1.5)
assertRaises(ValueError, math.log10, -10**1000)
assertRaises(ValueError, math.log10, NINF)
assertEqual(math.log(INF), INF)
assertTrue(math.isnan(math.log10(NAN)))

doc="Modf"
assertRaises(TypeError, math.modf)

def testmodf(name, result, expected):
    (v1, v2), (e1, e2) = result, expected
    if abs(v1-e1) > eps or abs(v2-e2):
        fail('%s returned %r, expected %r'%\
                  (name, result, expected))

testmodf('modf(1.5)', math.modf(1.5), (0.5, 1.0))
testmodf('modf(-1.5)', math.modf(-1.5), (-0.5, -1.0))

assertEqual(math.modf(INF), (0.0, INF))
assertEqual(math.modf(NINF), (-0.0, NINF))

modf_nan = math.modf(NAN)
assertTrue(math.isnan(modf_nan[0]))
assertTrue(math.isnan(modf_nan[1]))

doc="Pow"
assertRaises(TypeError, math.pow)
ftest('pow(0,1)', math.pow(0,1), 0)
ftest('pow(1,0)', math.pow(1,0), 1)
ftest('pow(2,1)', math.pow(2,1), 2)
ftest('pow(2,-1)', math.pow(2,-1), 0.5)
assertEqual(math.pow(INF, 1), INF)
assertEqual(math.pow(NINF, 1), NINF)
assertEqual((math.pow(1, INF)), 1.)
assertEqual((math.pow(1, NINF)), 1.)
assertTrue(math.isnan(math.pow(NAN, 1)))
assertTrue(math.isnan(math.pow(2, NAN)))
assertTrue(math.isnan(math.pow(0, NAN)))
assertEqual(math.pow(1, NAN), 1)

# pow(0., x)
assertEqual(math.pow(0., INF), 0.)
assertEqual(math.pow(0., 3.), 0.)
assertEqual(math.pow(0., 2.3), 0.)
assertEqual(math.pow(0., 2.), 0.)
assertEqual(math.pow(0., 0.), 1.)
assertEqual(math.pow(0., -0.), 1.)
assertRaises(ValueError, math.pow, 0., -2.)
assertRaises(ValueError, math.pow, 0., -2.3)
assertRaises(ValueError, math.pow, 0., -3.)
assertRaises(ValueError, math.pow, 0., NINF)
assertTrue(math.isnan(math.pow(0., NAN)))

# pow(INF, x)
assertEqual(math.pow(INF, INF), INF)
assertEqual(math.pow(INF, 3.), INF)
assertEqual(math.pow(INF, 2.3), INF)
assertEqual(math.pow(INF, 2.), INF)
assertEqual(math.pow(INF, 0.), 1.)
assertEqual(math.pow(INF, -0.), 1.)
assertEqual(math.pow(INF, -2.), 0.)
assertEqual(math.pow(INF, -2.3), 0.)
assertEqual(math.pow(INF, -3.), 0.)
assertEqual(math.pow(INF, NINF), 0.)
assertTrue(math.isnan(math.pow(INF, NAN)))

# pow(-0., x)
assertEqual(math.pow(-0., INF), 0.)
assertEqual(math.pow(-0., 3.), -0.)
assertEqual(math.pow(-0., 2.3), 0.)
assertEqual(math.pow(-0., 2.), 0.)
assertEqual(math.pow(-0., 0.), 1.)
assertEqual(math.pow(-0., -0.), 1.)
assertRaises(ValueError, math.pow, -0., -2.)
assertRaises(ValueError, math.pow, -0., -2.3)
assertRaises(ValueError, math.pow, -0., -3.)
assertRaises(ValueError, math.pow, -0., NINF)
assertTrue(math.isnan(math.pow(-0., NAN)))

# pow(NINF, x)
assertEqual(math.pow(NINF, INF), INF)
assertEqual(math.pow(NINF, 3.), NINF)
assertEqual(math.pow(NINF, 2.3), INF)
assertEqual(math.pow(NINF, 2.), INF)
assertEqual(math.pow(NINF, 0.), 1.)
assertEqual(math.pow(NINF, -0.), 1.)
assertEqual(math.pow(NINF, -2.), 0.)
assertEqual(math.pow(NINF, -2.3), 0.)
assertEqual(math.pow(NINF, -3.), -0.)
assertEqual(math.pow(NINF, NINF), 0.)
assertTrue(math.isnan(math.pow(NINF, NAN)))

# pow(-1, x)
assertEqual(math.pow(-1., INF), 1.)
assertEqual(math.pow(-1., 3.), -1.)
assertRaises(ValueError, math.pow, -1., 2.3)
assertEqual(math.pow(-1., 2.), 1.)
assertEqual(math.pow(-1., 0.), 1.)
assertEqual(math.pow(-1., -0.), 1.)
assertEqual(math.pow(-1., -2.), 1.)
assertRaises(ValueError, math.pow, -1., -2.3)
assertEqual(math.pow(-1., -3.), -1.)
assertEqual(math.pow(-1., NINF), 1.)
assertTrue(math.isnan(math.pow(-1., NAN)))

# pow(1, x)
assertEqual(math.pow(1., INF), 1.)
assertEqual(math.pow(1., 3.), 1.)
assertEqual(math.pow(1., 2.3), 1.)
assertEqual(math.pow(1., 2.), 1.)
assertEqual(math.pow(1., 0.), 1.)
assertEqual(math.pow(1., -0.), 1.)
assertEqual(math.pow(1., -2.), 1.)
assertEqual(math.pow(1., -2.3), 1.)
assertEqual(math.pow(1., -3.), 1.)
assertEqual(math.pow(1., NINF), 1.)
assertEqual(math.pow(1., NAN), 1.)

# pow(x, 0) should be 1 for any x
assertEqual(math.pow(2.3, 0.), 1.)
assertEqual(math.pow(-2.3, 0.), 1.)
assertEqual(math.pow(NAN, 0.), 1.)
assertEqual(math.pow(2.3, -0.), 1.)
assertEqual(math.pow(-2.3, -0.), 1.)
assertEqual(math.pow(NAN, -0.), 1.)

# pow(x, y) is invalid if x is negative and y is not integral
assertRaises(ValueError, math.pow, -1., 2.3)
assertRaises(ValueError, math.pow, -15., -3.1)

# pow(x, NINF)
assertEqual(math.pow(1.9, NINF), 0.)
assertEqual(math.pow(1.1, NINF), 0.)
assertEqual(math.pow(0.9, NINF), INF)
assertEqual(math.pow(0.1, NINF), INF)
assertEqual(math.pow(-0.1, NINF), INF)
assertEqual(math.pow(-0.9, NINF), INF)
assertEqual(math.pow(-1.1, NINF), 0.)
assertEqual(math.pow(-1.9, NINF), 0.)

# pow(x, INF)
assertEqual(math.pow(1.9, INF), INF)
assertEqual(math.pow(1.1, INF), INF)
assertEqual(math.pow(0.9, INF), 0.)
assertEqual(math.pow(0.1, INF), 0.)
assertEqual(math.pow(-0.1, INF), 0.)
assertEqual(math.pow(-0.9, INF), 0.)
assertEqual(math.pow(-1.1, INF), INF)
assertEqual(math.pow(-1.9, INF), INF)

# pow(x, y) should work for x negative, y an integer
ftest('(-2.)**3.', math.pow(-2.0, 3.0), -8.0)
ftest('(-2.)**2.', math.pow(-2.0, 2.0), 4.0)
ftest('(-2.)**1.', math.pow(-2.0, 1.0), -2.0)
ftest('(-2.)**0.', math.pow(-2.0, 0.0), 1.0)
ftest('(-2.)**-0.', math.pow(-2.0, -0.0), 1.0)
ftest('(-2.)**-1.', math.pow(-2.0, -1.0), -0.5)
ftest('(-2.)**-2.', math.pow(-2.0, -2.0), 0.25)
ftest('(-2.)**-3.', math.pow(-2.0, -3.0), -0.125)
assertRaises(ValueError, math.pow, -2.0, -0.5)
assertRaises(ValueError, math.pow, -2.0, 0.5)

# the following tests have been commented out since they don't
# really belong here:  the implementation of ** for floats is
# independent of the implementation of math.pow
assertEqual(1**NAN, 1)
assertEqual(1**INF, 1)
assertEqual(1**NINF, 1)
assertEqual(1**0, 1)
assertEqual(1.**NAN, 1)
assertEqual(1.**INF, 1)
assertEqual(1.**NINF, 1)
assertEqual(1.**0, 1)

doc="Radians"
assertRaises(TypeError, math.radians)
ftest('radians(180)', math.radians(180), math.pi)
ftest('radians(90)', math.radians(90), math.pi/2)
ftest('radians(-45)', math.radians(-45), -math.pi/4)

doc="Sin"
assertRaises(TypeError, math.sin)
ftest('sin(0)', math.sin(0), 0)
ftest('sin(pi/2)', math.sin(math.pi/2), 1)
ftest('sin(-pi/2)', math.sin(-math.pi/2), -1)
try:
    assertTrue(math.isnan(math.sin(INF)))
    assertTrue(math.isnan(math.sin(NINF)))
except ValueError:
    assertRaises(ValueError, math.sin, INF)
    assertRaises(ValueError, math.sin, NINF)
assertTrue(math.isnan(math.sin(NAN)))

doc="Sinh"
assertRaises(TypeError, math.sinh)
ftest('sinh(0)', math.sinh(0), 0)
ftest('sinh(1)**2-cosh(1)**2', math.sinh(1)**2-math.cosh(1)**2, -1)
ftest('sinh(1)+sinh(-1)', math.sinh(1)+math.sinh(-1), 0)
assertEqual(math.sinh(INF), INF)
assertEqual(math.sinh(NINF), NINF)
assertTrue(math.isnan(math.sinh(NAN)))

doc="Sqrt"
assertRaises(TypeError, math.sqrt)
ftest('sqrt(0)', math.sqrt(0), 0)
ftest('sqrt(1)', math.sqrt(1), 1)
ftest('sqrt(4)', math.sqrt(4), 2)
assertEqual(math.sqrt(INF), INF)
assertRaises(ValueError, math.sqrt, NINF)
assertTrue(math.isnan(math.sqrt(NAN)))

doc="Tan"
assertRaises(TypeError, math.tan)
ftest('tan(0)', math.tan(0), 0)
ftest('tan(pi/4)', math.tan(math.pi/4), 1)
ftest('tan(-pi/4)', math.tan(-math.pi/4), -1)
try:
    assertTrue(math.isnan(math.tan(INF)))
    assertTrue(math.isnan(math.tan(NINF)))
except:
    assertRaises(ValueError, math.tan, INF)
    assertRaises(ValueError, math.tan, NINF)
assertTrue(math.isnan(math.tan(NAN)))

doc="Tanh"
assertRaises(TypeError, math.tanh)
ftest('tanh(0)', math.tanh(0), 0)
ftest('tanh(1)+tanh(-1)', math.tanh(1)+math.tanh(-1), 0)
ftest('tanh(inf)', math.tanh(INF), 1)
ftest('tanh(-inf)', math.tanh(NINF), -1)
assertTrue(math.isnan(math.tanh(NAN)))

doc="TanhSign"
# check that tanh(-0.) == -0. on IEEE 754 systems
assertEqual(math.tanh(-0.), -0.)
assertEqual(math.copysign(1., math.tanh(-0.)),
                 math.copysign(1., -0.))

doc="trunc"
assertEqual(math.trunc(1), 1)
assertEqual(math.trunc(-1), -1)
assertEqual(type(math.trunc(1)), int)
assertEqual(type(math.trunc(1.5)), int)
assertEqual(math.trunc(1.5), 1)
assertEqual(math.trunc(-1.5), -1)
assertEqual(math.trunc(1.999999), 1)
assertEqual(math.trunc(-1.999999), -1)
assertEqual(math.trunc(-0.999999), -0)
assertEqual(math.trunc(-100.999), -100)

class TestTrunc(object):
    def __trunc__(self):
        return 23

class TestNoTrunc(object):
    pass

assertEqual(math.trunc(TestTrunc()), 23)

assertRaises(TypeError, math.trunc)
assertRaises(TypeError, math.trunc, 1, 2)
assertRaises(TypeError, math.trunc, TestNoTrunc())

doc="Isfinite"
assertTrue(math.isfinite(0.0))
assertTrue(math.isfinite(-0.0))
assertTrue(math.isfinite(1.0))
assertTrue(math.isfinite(-1.0))
assertFalse(math.isfinite(float("nan")))
assertFalse(math.isfinite(float("inf")))
assertFalse(math.isfinite(float("-inf")))

doc="Isnan"
assertTrue(math.isnan(float("nan")))
assertTrue(math.isnan(float("inf")* 0.))
assertFalse(math.isnan(float("inf")))
assertFalse(math.isnan(0.))
assertFalse(math.isnan(1.))

doc="Isinf"
assertTrue(math.isinf(float("inf")))
assertTrue(math.isinf(float("-inf")))
# FIXME assertTrue(math.isinf(1E400))
# FIXME assertTrue(math.isinf(-1E400))
assertFalse(math.isinf(float("nan")))
assertFalse(math.isinf(0.))
assertFalse(math.isinf(1.))

doc="exceptions"
try:
    x = math.exp(-1000000000)
except:
    # mathmodule.c is failing to weed out underflows from libm, or
    # we've got an fp format with huge dynamic range
    fail("underflowing exp() should not have raised "
                "an exception")
if x != 0:
    fail("underflowing exp() should have returned 0")

# If this fails, probably using a strict IEEE-754 conforming libm, and x
# is +Inf afterwards.  But Python wants overflows detected by default.
try:
    x = math.exp(1000000000)
except OverflowError:
    pass
else:
    fail("overflowing exp() didn't trigger OverflowError")

# If this fails, it could be a puzzle.  One odd possibility is that
# mathmodule.c's macros are getting confused while comparing
# Inf (HUGE_VAL) to a NaN, and artificially setting errno to ERANGE
# as a result (and so raising OverflowError instead).
try:
    x = math.sqrt(-1.0)
except ValueError:
    pass
else:
    fail("sqrt(-1) didn't raise ValueError")

doc="finished"
