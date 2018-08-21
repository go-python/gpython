# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Testcases for functions in math.
#
# Each line takes the form:
#
# <testid> <function> <input_value> -> <output_value> <flags>
#
# where:
#
#   <testid> is a short name identifying the test,
#
#   <function> is the function to be tested (exp, cos, asinh, ...),
#
#   <input_value> is a string representing a floating-point value
#
#   <output_value> is the expected (ideal) output value, again
#     represented as a string.
#
#   <flags> is a list of the floating-point flags required by C99
#
# The possible flags are:
#
#   divide-by-zero : raised when a finite input gives a
#     mathematically infinite result.
#
#   overflow : raised when a finite input gives a finite result that
#     is too large to fit in the usual range of an IEEE 754 double.
#
#   invalid : raised for invalid inputs (e.g., sqrt(-1))
#
#   ignore-sign : indicates that the sign of the result is
#     unspecified; e.g., if the result is given as inf,
#     then both -inf and inf should be accepted as correct.
#
# Flags may appear in any order.
#
# Lines beginning with '--' (like this one) start a comment, and are
# ignored.  Blank lines, or lines containing only whitespace, are also
# ignored.

# Many of the values below were computed with the help of
# version 2.4 of the MPFR library for multiple-precision
# floating-point computations with correct rounding.  All output
# values in this file are (modulo yet-to-be-discovered bugs)
# correctly rounded, provided that each input and output decimal
# floating-point value below is interpreted as a representation of
# the corresponding nearest IEEE 754 double-precision value.  See the
# MPFR homepage at http://www.mpfr.org for more information about the
# MPFR project.

import math
from libtest import *
from libulp import *
doc="testcases"
inf = float("inf")
nan = float("nan")

def tolerance(a, b, e):
    """Return if a-b is within tolerance e"""
    d = a - b
    if d < 0:
        d = -d
    if a != 0:
        e = e * a
        if e < 0:
            e = -e
    return d <= e

def acc_check(what, want, got, rel_err=2e-15, abs_err = 5e-323):
    """Determine whether non-NaN floats a and b are equal to within a
    (small) rounding error.  The default values for rel_err and
    abs_err are chosen to be suitable for platforms where a float is
    represented by an IEEE 754 double.  They allow an error of between
    9 and 19 ulps."""

    # need to special case infinities, since inf - inf gives nan
    if math.isinf(want) and got == want:
        return

    error = got - want

    permitted_error = rel_err * abs(want)
    if abs_err > permitted_error:
        permitted_error = abs_err
    if abs(error) < permitted_error:
        return
    raise AssertionError("%s: want %g, got %g: error = %g; permitted error = %g" % (what, want, got, error, permitted_error))

def t(name, fn, x, want, exc=None):
    global doc
    doc = name
    if exc is None:
        got = fn(x)
        if math.isnan(want) and math.isnan(got):
            return
        if want == inf and got == inf:
            return
        if want == -inf and got == -inf:
            return
        if fn == math.lgamma:
            # we use a weaker accuracy test for lgamma;
            # lgamma only achieves an absolute error of
            # a few multiples of the machine accuracy, in
            # general.
            acc_check(doc, want, got, rel_err = 5e-15, abs_err = 5e-15)
        elif fn == math.erfc:
            # erfc has less-than-ideal accuracy for large
            # arguments (x ~ 25 or so), mainly due to the
            # error involved in computing exp(-x*x).
            #
            # XXX Would be better to weaken this test only
            # for large x, instead of for all x.
            ulps_check(doc, want, got, 2000)
        else:
            ulps_check(doc, want, got, 20)
    else:
        try:
            got = fn(x)
        except exc as e:
            pass
        else:
            assert False, "%s not raised" % exc

#
# erf: error function --
#

t("erf0000", math.erf, 0.0, 0.0)
t("erf0001", math.erf, -0.0, -0.0)
t("erf0002", math.erf, inf, 1.0)
t("erf0003", math.erf, -inf, -1.0)
t("erf0004", math.erf, nan, nan)

# tiny values
t("erf0010", math.erf, 1e-308, 1.1283791670955125e-308)
t("erf0011", math.erf, 5e-324, 4.9406564584124654e-324)
t("erf0012", math.erf, 1e-10, 1.1283791670955126e-10)

# small integers
t("erf0020", math.erf, 1, 0.84270079294971489)
t("erf0021", math.erf, 2, 0.99532226501895271)
t("erf0022", math.erf, 3, 0.99997790950300136)
t("erf0023", math.erf, 4, 0.99999998458274209)
t("erf0024", math.erf, 5, 0.99999999999846256)
t("erf0025", math.erf, 6, 1.0)

t("erf0030", math.erf, -1, -0.84270079294971489)
t("erf0031", math.erf, -2, -0.99532226501895271)
t("erf0032", math.erf, -3, -0.99997790950300136)
t("erf0033", math.erf, -4, -0.99999998458274209)
t("erf0034", math.erf, -5, -0.99999999999846256)
t("erf0035", math.erf, -6, -1.0)

# huge values should all go to +/-1, depending on sign
t("erf0040", math.erf, -40, -1.0)
t("erf0041", math.erf, 1e16, 1.0)
t("erf0042", math.erf, -1e150, -1.0)
t("erf0043", math.erf, 1.7e308, 1.0)

# Issue 8986: inputs x with exp(-x*x) near the underflow threshold
# incorrectly signalled overflow on some platforms.
t("erf0100", math.erf, 26.2, 1.0)
t("erf0101", math.erf, 26.4, 1.0)
t("erf0102", math.erf, 26.6, 1.0)
t("erf0103", math.erf, 26.8, 1.0)
t("erf0104", math.erf, 27.0, 1.0)
t("erf0105", math.erf, 27.2, 1.0)
t("erf0106", math.erf, 27.4, 1.0)
t("erf0107", math.erf, 27.6, 1.0)

t("erf0110", math.erf, -26.2, -1.0)
t("erf0111", math.erf, -26.4, -1.0)
t("erf0112", math.erf, -26.6, -1.0)
t("erf0113", math.erf, -26.8, -1.0)
t("erf0114", math.erf, -27.0, -1.0)
t("erf0115", math.erf, -27.2, -1.0)
t("erf0116", math.erf, -27.4, -1.0)
t("erf0117", math.erf, -27.6, -1.0)

#
# erfc: complementary error function --
#

t("erfc0000", math.erfc, 0.0, 1.0)
t("erfc0001", math.erfc, -0.0, 1.0)
t("erfc0002", math.erfc, inf, 0.0)
t("erfc0003", math.erfc, -inf, 2.0)
t("erfc0004", math.erfc, nan, nan)

# tiny values
t("erfc0010", math.erfc, 1e-308, 1.0)
t("erfc0011", math.erfc, 5e-324, 1.0)
t("erfc0012", math.erfc, 1e-10, 0.99999999988716204)

# small integers
t("erfc0020", math.erfc, 1, 0.15729920705028513)
t("erfc0021", math.erfc, 2, 0.0046777349810472662)
t("erfc0022", math.erfc, 3, 2.2090496998585441e-05)
t("erfc0023", math.erfc, 4, 1.541725790028002e-08)
t("erfc0024", math.erfc, 5, 1.5374597944280349e-12)
t("erfc0025", math.erfc, 6, 2.1519736712498913e-17)

t("erfc0030", math.erfc, -1, 1.8427007929497148)
t("erfc0031", math.erfc, -2, 1.9953222650189528)
t("erfc0032", math.erfc, -3, 1.9999779095030015)
t("erfc0033", math.erfc, -4, 1.9999999845827421)
t("erfc0034", math.erfc, -5, 1.9999999999984626)
t("erfc0035", math.erfc, -6, 2.0)

# as x -> infinity, erfc(x) behaves like exp(-x*x)/x/sqrt(pi)
t("erfc0040", math.erfc, 20, 5.3958656116079012e-176)
t("erfc0041", math.erfc, 25, 8.3001725711965228e-274)
# FIXME(underflows to 0) t("erfc0042", math.erfc, 27, 5.2370464393526292e-319)
t("erfc0043", math.erfc, 28, 0.0)

# huge values
t("erfc0050", math.erfc, -40, 2.0)
t("erfc0051", math.erfc, 1e16, 0.0)
t("erfc0052", math.erfc, -1e150, 2.0)
t("erfc0053", math.erfc, 1.7e308, 0.0)

# Issue 8986: inputs x with exp(-x*x) near the underflow threshold
# incorrectly signalled overflow on some platforms.
t("erfc0100", math.erfc, 26.2, 1.6432507924389461e-300)
t("erfc0101", math.erfc, 26.4, 4.4017768588035426e-305)
t("erfc0102", math.erfc, 26.6, 1.0885125885442269e-309)
# FIXME(underflows to 0) t("erfc0103", math.erfc, 26.8, 2.4849621571966629e-314)
# FIXME(underflows to 0) t("erfc0104", math.erfc, 27.0, 5.2370464393526292e-319)
# FIXME(underflows to 0) t("erfc0105", math.erfc, 27.2, 9.8813129168249309e-324)
t("erfc0106", math.erfc, 27.4, 0.0)
t("erfc0107", math.erfc, 27.6, 0.0)

t("erfc0110", math.erfc, -26.2, 2.0)
t("erfc0111", math.erfc, -26.4, 2.0)
t("erfc0112", math.erfc, -26.6, 2.0)
t("erfc0113", math.erfc, -26.8, 2.0)
t("erfc0114", math.erfc, -27.0, 2.0)
t("erfc0115", math.erfc, -27.2, 2.0)
t("erfc0116", math.erfc, -27.4, 2.0)
t("erfc0117", math.erfc, -27.6, 2.0)

#
# lgamma: log of absolute value of the gamma function --
#

# special values
t("lgam0000", math.lgamma, 0.0, inf, ValueError)
t("lgam0001", math.lgamma, -0.0, inf, ValueError)
t("lgam0002", math.lgamma, inf, inf)
# FIXME(ValueError) t("lgam0003", math.lgamma, -inf, inf)
t("lgam0004", math.lgamma, nan, nan)

# negative integers
t("lgam0010", math.lgamma, -1, inf, ValueError)
t("lgam0011", math.lgamma, -2, inf, ValueError)
t("lgam0012", math.lgamma, -1e16, inf, ValueError)
t("lgam0013", math.lgamma, -1e300, inf, ValueError)
t("lgam0014", math.lgamma, -1.79e308, inf, ValueError)

# small positive integers give factorials
t("lgam0020", math.lgamma, 1, 0.0)
t("lgam0021", math.lgamma, 2, 0.0)
t("lgam0022", math.lgamma, 3, 0.69314718055994529)
t("lgam0023", math.lgamma, 4, 1.791759469228055)
t("lgam0024", math.lgamma, 5, 3.1780538303479458)
t("lgam0025", math.lgamma, 6, 4.7874917427820458)

# half integers
t("lgam0030", math.lgamma, 0.5, 0.57236494292470008)
t("lgam0031", math.lgamma, 1.5, -0.12078223763524522)
t("lgam0032", math.lgamma, 2.5, 0.28468287047291918)
t("lgam0033", math.lgamma, 3.5, 1.2009736023470743)
t("lgam0034", math.lgamma, -0.5, 1.2655121234846454)
t("lgam0035", math.lgamma, -1.5, 0.86004701537648098)
t("lgam0036", math.lgamma, -2.5, -0.056243716497674054)
t("lgam0037", math.lgamma, -3.5, -1.309006684993042)

# values near 0
t("lgam0040", math.lgamma, 0.1, 2.252712651734206)
t("lgam0041", math.lgamma, 0.01, 4.5994798780420219)
t("lgam0042", math.lgamma, 1e-8, 18.420680738180209)
t("lgam0043", math.lgamma, 1e-16, 36.841361487904734)
t("lgam0044", math.lgamma, 1e-30, 69.077552789821368)
t("lgam0045", math.lgamma, 1e-160, 368.41361487904732)
# FIXME(inaccurate) t("lgam0046", math.lgamma, 1e-308, 709.19620864216608)
# FIXME(inaccurate) t("lgam0047", math.lgamma, 5.6e-309, 709.77602713741896)
# FIXME(inaccurate) t("lgam0048", math.lgamma, 5.5e-309, 709.79404564292167)
# FIXME(inaccurate) t("lgam0049", math.lgamma, 1e-309, 711.49879373516012)
# FIXME(inaccurate) t("lgam0050", math.lgamma, 1e-323, 743.74692474082133)
# FIXME(inaccurate) t("lgam0051", math.lgamma, 5e-324, 744.44007192138122)
t("lgam0060", math.lgamma, -0.1, 2.3689613327287886)
t("lgam0061", math.lgamma, -0.01, 4.6110249927528013)
t("lgam0062", math.lgamma, -1e-8, 18.420680749724522)
t("lgam0063", math.lgamma, -1e-16, 36.841361487904734)
t("lgam0064", math.lgamma, -1e-30, 69.077552789821368)
t("lgam0065", math.lgamma, -1e-160, 368.41361487904732)
# FIXME(inaccurate) t("lgam0066", math.lgamma, -1e-308, 709.19620864216608)
# FIXME(inaccurate) t("lgam0067", math.lgamma, -5.6e-309, 709.77602713741896)
# FIXME(inaccurate) t("lgam0068", math.lgamma, -5.5e-309, 709.79404564292167)
# FIXME(inaccurate) t("lgam0069", math.lgamma, -1e-309, 711.49879373516012)
# FIXME(inaccurate) t("lgam0070", math.lgamma, -1e-323, 743.74692474082133)
# FIXME(inaccurate) t("lgam0071", math.lgamma, -5e-324, 744.44007192138122)

# values near negative integers
t("lgam0080", math.lgamma, -0.99999999999999989, 36.736800569677101)
t("lgam0081", math.lgamma, -1.0000000000000002, 36.043653389117154)
t("lgam0082", math.lgamma, -1.9999999999999998, 35.350506208557213)
t("lgam0083", math.lgamma, -2.0000000000000004, 34.657359027997266)
t("lgam0084", math.lgamma, -100.00000000000001, -331.85460524980607)
t("lgam0085", math.lgamma, -99.999999999999986, -331.85460524980596)

# large inputs
t("lgam0100", math.lgamma, 170, 701.43726380873704)
t("lgam0101", math.lgamma, 171, 706.57306224578736)
t("lgam0102", math.lgamma, 171.624, 709.78077443669895)
t("lgam0103", math.lgamma, 171.625, 709.78591682948365)
t("lgam0104", math.lgamma, 172, 711.71472580228999)
t("lgam0105", math.lgamma, 2000, 13198.923448054265)
t("lgam0106", math.lgamma, 2.55998332785163e305, 1.7976931348623099e+308)
t("lgam0107", math.lgamma, 2.55998332785164e305, inf, OverflowError)
t("lgam0108", math.lgamma, 1.7e308, inf, OverflowError)

# inputs for which gamma(x) is tiny
t("lgam0120", math.lgamma, -100.5, -364.90096830942736)
t("lgam0121", math.lgamma, -160.5, -656.88005261126432)
t("lgam0122", math.lgamma, -170.5, -707.99843314507882)
t("lgam0123", math.lgamma, -171.5, -713.14301641168481)
t("lgam0124", math.lgamma, -176.5, -738.95247590846486)
t("lgam0125", math.lgamma, -177.5, -744.13144651738037)
t("lgam0126", math.lgamma, -178.5, -749.3160351186001)

t("lgam0130", math.lgamma, -1000.5, -5914.4377011168517)
t("lgam0131", math.lgamma, -30000.5, -279278.6629959144)
# FIXME t("lgam0132", math.lgamma, -4503599627370495.5, -1.5782258434492883e+17)

# results close to 0:  positive argument ...
t("lgam0150", math.lgamma, 0.99999999999999989, 6.4083812134800075e-17)
t("lgam0151", math.lgamma, 1.0000000000000002, -1.2816762426960008e-16)
t("lgam0152", math.lgamma, 1.9999999999999998, -9.3876980655431170e-17)
t("lgam0153", math.lgamma, 2.0000000000000004, 1.8775396131086244e-16)

# ... and negative argument
# these are very inaccurate in python3
t("lgam0160", math.lgamma, -2.7476826467, -5.2477408147689136e-11)
t("lgam0161", math.lgamma, -2.457024738, 3.3464637541912932e-10)


#
# gamma: Gamma function --
#

# special values
t("gam0000", math.gamma, 0.0, inf, ValueError)
t("gam0001", math.gamma, -0.0, -inf, ValueError)
t("gam0002", math.gamma, inf, inf)
t("gam0003", math.gamma, -inf, nan, ValueError)
t("gam0004", math.gamma, nan, nan)

# negative integers inputs are invalid
t("gam0010", math.gamma, -1, nan, ValueError)
t("gam0011", math.gamma, -2, nan, ValueError)
t("gam0012", math.gamma, -1e16, nan, ValueError)
t("gam0013", math.gamma, -1e300, nan, ValueError)

# small positive integers give factorials
t("gam0020", math.gamma, 1, 1)
t("gam0021", math.gamma, 2, 1)
t("gam0022", math.gamma, 3, 2)
t("gam0023", math.gamma, 4, 6)
t("gam0024", math.gamma, 5, 24)
t("gam0025", math.gamma, 6, 120)

# half integers
t("gam0030", math.gamma, 0.5, 1.7724538509055161)
t("gam0031", math.gamma, 1.5, 0.88622692545275805)
t("gam0032", math.gamma, 2.5, 1.3293403881791370)
t("gam0033", math.gamma, 3.5, 3.3233509704478426)
t("gam0034", math.gamma, -0.5, -3.5449077018110322)
t("gam0035", math.gamma, -1.5, 2.3632718012073548)
t("gam0036", math.gamma, -2.5, -0.94530872048294190)
t("gam0037", math.gamma, -3.5, 0.27008820585226911)

# values near 0
t("gam0040", math.gamma, 0.1, 9.5135076986687306)
t("gam0041", math.gamma, 0.01, 99.432585119150602)
t("gam0042", math.gamma, 1e-8, 99999999.422784343)
t("gam0043", math.gamma, 1e-16, 10000000000000000)
t("gam0044", math.gamma, 1e-30, 9.9999999999999988e+29)
t("gam0045", math.gamma, 1e-160, 1.0000000000000000e+160)
t("gam0046", math.gamma, 1e-308, 1.0000000000000000e+308)
t("gam0047", math.gamma, 5.6e-309, 1.7857142857142848e+308)
t("gam0048", math.gamma, 5.5e-309, inf, OverflowError)
t("gam0049", math.gamma, 1e-309, inf, OverflowError)
t("gam0050", math.gamma, 1e-323, inf, OverflowError)
t("gam0051", math.gamma, 5e-324, inf, OverflowError)
t("gam0060", math.gamma, -0.1, -10.686287021193193)
t("gam0061", math.gamma, -0.01, -100.58719796441078)
t("gam0062", math.gamma, -1e-8, -100000000.57721567)
t("gam0063", math.gamma, -1e-16, -10000000000000000)
t("gam0064", math.gamma, -1e-30, -9.9999999999999988e+29)
t("gam0065", math.gamma, -1e-160, -1.0000000000000000e+160)
t("gam0066", math.gamma, -1e-308, -1.0000000000000000e+308)
t("gam0067", math.gamma, -5.6e-309, -1.7857142857142848e+308)
t("gam0068", math.gamma, -5.5e-309, -inf, OverflowError)
t("gam0069", math.gamma, -1e-309, -inf, OverflowError)
t("gam0070", math.gamma, -1e-323, -inf, OverflowError)
t("gam0071", math.gamma, -5e-324, -inf, OverflowError)

# values near negative integers
t("gam0080", math.gamma, -0.99999999999999989, -9007199254740992.0)
t("gam0081", math.gamma, -1.0000000000000002, 4503599627370495.5)
t("gam0082", math.gamma, -1.9999999999999998, 2251799813685248.5)
t("gam0083", math.gamma, -2.0000000000000004, -1125899906842623.5)
t("gam0084", math.gamma, -100.00000000000001, -7.5400833348831090e-145)
t("gam0085", math.gamma, -99.999999999999986, 7.5400833348840962e-145)

# large inputs
t("gam0100", math.gamma, 170, 4.2690680090047051e+304)
t("gam0101", math.gamma, 171, 7.2574156153079990e+306)
# FIXME(overflows) t("gam0102", math.gamma, 171.624, 1.7942117599248104e+308)
t("gam0103", math.gamma, 171.625, inf, OverflowError)
t("gam0104", math.gamma, 172, inf, OverflowError)
t("gam0105", math.gamma, 2000, inf, OverflowError)
t("gam0106", math.gamma, 1.7e308, inf, OverflowError)

# inputs for which gamma(x) is tiny
t("gam0120", math.gamma, -100.5, -3.3536908198076787e-159)
t("gam0121", math.gamma, -160.5, -5.2555464470078293e-286)
t("gam0122", math.gamma, -170.5, -3.3127395215386074e-308)
# Reported as https://github.com/golang/go/issues/11441
# FIXME(overflows) t("gam0123", math.gamma, -171.5, 1.9316265431711902e-310)
# FIXME(overflows) t("gam0124", math.gamma, -176.5, -1.1956388629358166e-321)
# FIXME(overflows) t("gam0125", math.gamma, -177.5, 4.9406564584124654e-324)
# FIXME(overflows) t("gam0126", math.gamma, -178.5, -0.0)
# FIXME(overflows) t("gam0127", math.gamma, -179.5, 0.0)
# FIXME(overflows) t("gam0128", math.gamma, -201.0001, 0.0)
# FIXME(overflows) t("gam0129", math.gamma, -202.9999, -0.0)
# FIXME(overflows) t("gam0130", math.gamma, -1000.5, -0.0)
# FIXME(overflows) t("gam0131", math.gamma, -1000000000.3, -0.0)
# FIXME(overflows) t("gam0132", math.gamma, -4503599627370495.5, 0.0)

# inputs that cause problems for the standard reflection formula,
# thanks to loss of accuracy in 1-x
t("gam0140", math.gamma, -63.349078729022985, 4.1777971677761880e-88)
t("gam0141", math.gamma, -127.45117632943295, 1.1831110896236810e-214)


#
# log1p: log(1 + x), without precision loss for small x --
#

# special values
t("log1p0000", math.log1p, 0.0, 0.0)
t("log1p0001", math.log1p, -0.0, -0.0)
t("log1p0002", math.log1p, inf, inf)
t("log1p0003", math.log1p, -inf, nan, ValueError)
t("log1p0004", math.log1p, nan, nan)

# singularity at -1.0
t("log1p0010", math.log1p, -1.0, -inf, ValueError)
t("log1p0011", math.log1p, -0.9999999999999999, -36.736800569677101)

# finite values < 1.0 are invalid
t("log1p0020", math.log1p, -1.0000000000000002, nan, ValueError)
t("log1p0021", math.log1p, -1.1, nan, ValueError)
t("log1p0022", math.log1p, -2.0, nan, ValueError)
t("log1p0023", math.log1p, -1e300, nan, ValueError)

# tiny x: log1p(x) ~ x
t("log1p0110", math.log1p, 5e-324, 5e-324)
t("log1p0111", math.log1p, 1e-320, 1e-320)
t("log1p0112", math.log1p, 1e-300, 1e-300)
t("log1p0113", math.log1p, 1e-150, 1e-150)
t("log1p0114", math.log1p, 1e-20, 1e-20)

t("log1p0120", math.log1p, -5e-324, -5e-324)
t("log1p0121", math.log1p, -1e-320, -1e-320)
t("log1p0122", math.log1p, -1e-300, -1e-300)
t("log1p0123", math.log1p, -1e-150, -1e-150)
t("log1p0124", math.log1p, -1e-20, -1e-20)

# some (mostly) random small and moderate-sized values
t("log1p0200", math.log1p, -0.89156889782277482, -2.2216403106762863)
t("log1p0201", math.log1p, -0.23858496047770464, -0.27257668276980057)
t("log1p0202", math.log1p, -0.011641726191307515, -0.011710021654495657)
t("log1p0203", math.log1p, -0.0090126398571693817, -0.0090534993825007650)
t("log1p0204", math.log1p, -0.00023442805985712781, -0.00023445554240995693)
t("log1p0205", math.log1p, -1.5672870980936349e-5, -1.5672993801662046e-5)
t("log1p0206", math.log1p, -7.9650013274825295e-6, -7.9650330482740401e-6)
t("log1p0207", math.log1p, -2.5202948343227410e-7, -2.5202951519170971e-7)
t("log1p0208", math.log1p, -8.2446372820745855e-11, -8.2446372824144559e-11)
t("log1p0209", math.log1p, -8.1663670046490789e-12, -8.1663670046824230e-12)
t("log1p0210", math.log1p, 7.0351735084656292e-18, 7.0351735084656292e-18)
t("log1p0211", math.log1p, 5.2732161907375226e-12, 5.2732161907236188e-12)
t("log1p0212", math.log1p, 1.0000000000000000e-10, 9.9999999995000007e-11)
t("log1p0213", math.log1p, 2.1401273266000197e-9, 2.1401273243099470e-9)
t("log1p0214", math.log1p, 1.2668914653979560e-8, 1.2668914573728861e-8)
t("log1p0215", math.log1p, 1.6250007816299069e-6, 1.6249994613175672e-6)
t("log1p0216", math.log1p, 8.3740495645839399e-6, 8.3740145024266269e-6)
t("log1p0217", math.log1p, 3.0000000000000001e-5, 2.9999550008999799e-5)
t("log1p0218", math.log1p, 0.0070000000000000001, 0.0069756137364252423)
t("log1p0219", math.log1p, 0.013026235315053002, 0.012942123564008787)
t("log1p0220", math.log1p, 0.013497160797236184, 0.013406885521915038)
t("log1p0221", math.log1p, 0.027625599078135284, 0.027250897463483054)
t("log1p0222", math.log1p, 0.14179687245544870, 0.13260322540908789)

# large values
t("log1p0300", math.log1p, 1.7976931348623157e+308, 709.78271289338397)
t("log1p0301", math.log1p, 1.0000000000000001e+300, 690.77552789821368)
t("log1p0302", math.log1p, 1.0000000000000001e+70, 161.18095650958321)
t("log1p0303", math.log1p, 10000000000.000000, 23.025850930040455)

# other values transferred from testLog1p in test_math
t("log1p0400", math.log1p, -0.63212055882855767, -1.0000000000000000)
t("log1p0401", math.log1p, 1.7182818284590451, 1.0000000000000000)
t("log1p0402", math.log1p, 1.0000000000000000, 0.69314718055994529)
t("log1p0403", math.log1p, 1.2379400392853803e+27, 62.383246250395075)


#
# expm1: exp(x) - 1, without precision loss for small x --
#

# special values
t("expm10000", math.expm1, 0.0, 0.0)
t("expm10001", math.expm1, -0.0, -0.0)
t("expm10002", math.expm1, inf, inf)
t("expm10003", math.expm1, -inf, -1.0)
t("expm10004", math.expm1, nan, nan)

# expm1(x) ~ x for tiny x
t("expm10010", math.expm1, 5e-324, 5e-324)
t("expm10011", math.expm1, 1e-320, 1e-320)
t("expm10012", math.expm1, 1e-300, 1e-300)
t("expm10013", math.expm1, 1e-150, 1e-150)
t("expm10014", math.expm1, 1e-20, 1e-20)

t("expm10020", math.expm1, -5e-324, -5e-324)
t("expm10021", math.expm1, -1e-320, -1e-320)
t("expm10022", math.expm1, -1e-300, -1e-300)
t("expm10023", math.expm1, -1e-150, -1e-150)
t("expm10024", math.expm1, -1e-20, -1e-20)

# moderate sized values, where direct evaluation runs into trouble
t("expm10100", math.expm1, 1e-10, 1.0000000000500000e-10)
t("expm10101", math.expm1, -9.9999999999999995e-08, -9.9999995000000163e-8)
t("expm10102", math.expm1, 3.0000000000000001e-05, 3.0000450004500034e-5)
t("expm10103", math.expm1, -0.0070000000000000001, -0.0069755570667648951)
t("expm10104", math.expm1, -0.071499208740094633, -0.069002985744820250)
t("expm10105", math.expm1, -0.063296004180116799, -0.061334416373633009)
t("expm10106", math.expm1, 0.02390954035597756, 0.024197665143819942)
t("expm10107", math.expm1, 0.085637352649044901, 0.089411184580357767)
t("expm10108", math.expm1, 0.5966174947411006, 0.81596588596501485)
t("expm10109", math.expm1, 0.30247206212075139, 0.35319987035848677)
t("expm10110", math.expm1, 0.74574727375889516, 1.1080161116737459)
t("expm10111", math.expm1, 0.97767512926555711, 1.6582689207372185)
t("expm10112", math.expm1, 0.8450154566787712, 1.3280137976535897)
t("expm10113", math.expm1, -0.13979260323125264, -0.13046144381396060)
t("expm10114", math.expm1, -0.52899322039643271, -0.41080213643695923)
t("expm10115", math.expm1, -0.74083261478900631, -0.52328317124797097)
t("expm10116", math.expm1, -0.93847766984546055, -0.60877704724085946)
t("expm10117", math.expm1, 10.0, 22025.465794806718)
t("expm10118", math.expm1, 27.0, 532048240600.79865)
t("expm10119", math.expm1, 123, 2.6195173187490626e+53)
t("expm10120", math.expm1, -12.0, -0.99999385578764666)
t("expm10121", math.expm1, -35.100000000000001, -0.99999999999999944)

# extreme negative values
t("expm10201", math.expm1, -37.0, -0.99999999999999989)
t("expm10200", math.expm1, -38.0, -1.0)
# FIXME(overflows) t("expm10210", math.expm1, -710.0, -1.0)
# the formula expm1(x) = 2 * sinh(x/2) * exp(x/2) doesn't work so
# well when exp(x/2) is subnormal or underflows to zero; check we're
# not using it!
# Reported as https://github.com/golang/go/issues/11442
# FIXME(overflows) t("expm10211", math.expm1, -1420.0, -1.0)
# FIXME(overflows) t("expm10212", math.expm1, -1450.0, -1.0)
# FIXME(overflows) t("expm10213", math.expm1, -1500.0, -1.0)
# FIXME(overflows) t("expm10214", math.expm1, -1e50, -1.0)
# FIXME(overflows) t("expm10215", math.expm1, -1.79e308, -1.0)

# extreme positive values
# FIXME(fails on 32 bit) t("expm10300", math.expm1, 300, 1.9424263952412558e+130)
# FIXME(fails on 32 bit) t("expm10301", math.expm1, 700, 1.0142320547350045e+304)
# the next test (expm10302) is disabled because it causes failure on
# OS X 10.4/Intel: apparently all values over 709.78 produce an
# overflow on that platform.  See issue #7575.
# expm10302 expm1 709.78271289328393 -> 1.7976931346824240e+308
t("expm10303", math.expm1, 709.78271289348402, inf, OverflowError)
t("expm10304", math.expm1, 1000, inf, OverflowError)
t("expm10305", math.expm1, 1e50, inf, OverflowError)
t("expm10306", math.expm1, 1.79e308, inf, OverflowError)

# weaker version of expm10302
# FIXME(fails on 32 bit) t("expm10307", math.expm1, 709.5, 1.3549863193146328e+308)

#
# log2: log to base 2 --
#

# special values
t("log20000", math.log2, 0.0, -inf, ValueError)
t("log20001", math.log2, -0.0, -inf, ValueError)
t("log20002", math.log2, inf, inf)
t("log20003", math.log2, -inf, nan, ValueError)
t("log20004", math.log2, nan, nan)

# exact value at 1.0
t("log20010", math.log2, 1.0, 0.0)

# negatives
t("log20020", math.log2, -5e-324, nan, ValueError)
t("log20021", math.log2, -1.0, nan, ValueError)
t("log20022", math.log2, -1.7e-308, nan, ValueError)

# exact values at powers of 2
t("log20100", math.log2, 2.0, 1.0)
t("log20101", math.log2, 4.0, 2.0)
t("log20102", math.log2, 8.0, 3.0)
t("log20103", math.log2, 16.0, 4.0)
t("log20104", math.log2, 32.0, 5.0)
t("log20105", math.log2, 64.0, 6.0)
t("log20106", math.log2, 128.0, 7.0)
t("log20107", math.log2, 256.0, 8.0)
t("log20108", math.log2, 512.0, 9.0)
t("log20109", math.log2, 1024.0, 10.0)
t("log20110", math.log2, 2048.0, 11.0)

t("log20200", math.log2, 0.5, -1.0)
t("log20201", math.log2, 0.25, -2.0)
t("log20202", math.log2, 0.125, -3.0)
t("log20203", math.log2, 0.0625, -4.0)

# values close to 1.0
# FIXME(inaccurate) t("log20300", math.log2, 1.0000000000000002, 3.2034265038149171e-16)
# FIXME(inaccurate) t("log20301", math.log2, 1.0000000001, 1.4426951601859516e-10)
# FIXME(inaccurate) t("log20302", math.log2, 1.00001, 1.4426878274712997e-5)

t("log20310", math.log2, 0.9999999999999999, -1.6017132519074588e-16)
t("log20311", math.log2, 0.9999999999, -1.4426951603302210e-10)
t("log20312", math.log2, 0.99999, -1.4427022544056922e-5)

# tiny values
t("log20400", math.log2, 5e-324, -1074.0)
t("log20401", math.log2, 1e-323, -1073.0)
t("log20402", math.log2, 1.5e-323, -1072.4150374992789)
t("log20403", math.log2, 2e-323, -1072.0)

t("log20410", math.log2, 1e-308, -1023.1538532253076)
t("log20411", math.log2, 2.2250738585072014e-308, -1022.0)
t("log20412", math.log2, 4.4501477170144028e-308, -1021.0)
t("log20413", math.log2, 1e-307, -1019.8319251304202)

# huge values
t("log20500", math.log2, 1.7976931348623157e+308, 1024.0)
t("log20501", math.log2, 1.7e+308, 1023.9193879716706)
t("log20502", math.log2, 8.9884656743115795e+307, 1023.0)

# selection of random values
t("log20600", math.log2, -7.2174324841039838e+289, nan, ValueError)
t("log20601", math.log2, -2.861319734089617e+265, nan, ValueError)
t("log20602", math.log2, -4.3507646894008962e+257, nan, ValueError)
t("log20603", math.log2, -6.6717265307520224e+234, nan, ValueError)
t("log20604", math.log2, -3.9118023786619294e+229, nan, ValueError)
t("log20605", math.log2, -1.5478221302505161e+206, nan, ValueError)
t("log20606", math.log2, -1.4380485131364602e+200, nan, ValueError)
t("log20607", math.log2, -3.7235198730382645e+185, nan, ValueError)
t("log20608", math.log2, -1.0472242235095724e+184, nan, ValueError)
t("log20609", math.log2, -5.0141781956163884e+160, nan, ValueError)
t("log20610", math.log2, -2.1157958031160324e+124, nan, ValueError)
t("log20611", math.log2, -7.9677558612567718e+90, nan, ValueError)
t("log20612", math.log2, -5.5553906194063732e+45, nan, ValueError)
t("log20613", math.log2, -16573900952607.953, nan, ValueError)
t("log20614", math.log2, -37198371019.888618, nan, ValueError)
t("log20615", math.log2, -6.0727115121422674e-32, nan, ValueError)
t("log20616", math.log2, -2.5406841656526057e-38, nan, ValueError)
t("log20617", math.log2, -4.9056766703267657e-43, nan, ValueError)
t("log20618", math.log2, -2.1646786075228305e-71, nan, ValueError)
t("log20619", math.log2, -2.470826790488573e-78, nan, ValueError)
t("log20620", math.log2, -3.8661709303489064e-165, nan, ValueError)
t("log20621", math.log2, -1.0516496976649986e-182, nan, ValueError)
t("log20622", math.log2, -1.5935458614317996e-255, nan, ValueError)
t("log20623", math.log2, -2.8750977267336654e-293, nan, ValueError)
t("log20624", math.log2, -7.6079466794732585e-296, nan, ValueError)
t("log20625", math.log2, 3.2073253539988545e-307, -1018.1505544209213)
t("log20626", math.log2, 1.674937885472249e-244, -809.80634755783126)
t("log20627", math.log2, 1.0911259044931283e-214, -710.76679472274213)
t("log20628", math.log2, 2.0275372624809709e-154, -510.55719818383272)
t("log20629", math.log2, 7.3926087369631841e-115, -379.13564735312292)
t("log20630", math.log2, 1.3480198206342423e-86, -285.25497445094436)
t("log20631", math.log2, 8.9927384655719947e-83, -272.55127136401637)
t("log20632", math.log2, 3.1452398713597487e-60, -197.66251564496875)
t("log20633", math.log2, 7.0706573215457351e-55, -179.88420087782217)
t("log20634", math.log2, 3.1258285390731669e-49, -161.13023800505653)
t("log20635", math.log2, 8.2253046627829942e-41, -133.15898277355879)
t("log20636", math.log2, 7.8691367397519897e+49, 165.75068202732419)
t("log20637", math.log2, 2.9920561983925013e+64, 214.18453534573757)
t("log20638", math.log2, 4.7827254553946841e+77, 258.04629628445673)
t("log20639", math.log2, 3.1903566496481868e+105, 350.47616767491166)
t("log20640", math.log2, 5.6195082449502419e+113, 377.86831861008250)
t("log20641", math.log2, 9.9625658250651047e+125, 418.55752921228753)
t("log20642", math.log2, 2.7358945220961532e+145, 483.13158636923413)
t("log20643", math.log2, 2.785842387926931e+174, 579.49360214860280)
t("log20644", math.log2, 2.4169172507252751e+193, 642.40529039289652)
t("log20645", math.log2, 3.1689091206395632e+205, 682.65924573798395)
t("log20646", math.log2, 2.535995592365391e+208, 692.30359597460460)
t("log20647", math.log2, 6.2011236566089916e+233, 776.64177576730913)
t("log20648", math.log2, 2.1843274820677632e+253, 841.57499717289647)
t("log20649", math.log2, 8.7493931063474791e+297, 989.74182713073981)

doc="finished"
