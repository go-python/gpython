// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/* Math module -- standard C math library functions, pi and e */
package math

// For cpython's tests see
//   Lib/test/test_math.py and
//   Lib/test/math_testcases.txt
//   Lib/test/cmath_testcases.txt

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"

	"github.com/go-python/gpython/py"
)

/* Here are some comments from Tim Peters, extracted from the
   discussion attached to http://bugs.python.org/issue1640.  They
   describe the general aims of the math module with respect to
   special values, IEEE-754 floating-point exceptions, and Python
   exceptions.

These are the "spirit of 754" rules:

1. If the mathematical result is a real number, but of magnitude too
large to approximate by a machine float, overflow is signaled and the
result is an infinity (with the appropriate sign).

2. If the mathematical result is a real number, but of magnitude too
small to approximate by a machine float, underflow is signaled and the
result is a zero (with the appropriate sign).

3. At a singularity (a value x such that the limit of f(y) as y
approaches x exists and is an infinity), "divide by zero" is signaled
and the result is an infinity (with the appropriate sign).  This is
complicated a little by that the left-side and right-side limits may
not be the same; e.g., 1/x approaches +inf or -inf as x approaches 0
from the positive or negative directions.  In that specific case, the
sign of the zero determines the result of 1/0.

4. At a point where a function has no defined result in the extended
reals (i.e., the reals plus an infinity or two), invalid operation is
signaled and a NaN is returned.

And these are what Python has historically /tried/ to do (but not
always successfully, as platform libm behavior varies a lot):

For #1, raise OverflowError.

For #2, return a zero (with the appropriate sign if that happens by
accident ;-)).

For #3 and #4, raise ValueError.  It may have made sense to raise
Python's ZeroDivisionError in #3, but historically that's only been
raised for division by zero and mod by zero.

*/

/*
   In general, on an IEEE-754 platform the aim is to follow the C99
   standard, including Annex 'F', whenever possible.  Where the
   standard recommends raising the 'divide-by-zero' or 'invalid'
   floating-point exceptions, Python should raise a ValueError.  Where
   the standard recommends raising 'overflow', Python should raise an
   OverflowError.  In all other circumstances a value should be
   returned.
*/
var (
	EDOM   = py.ExceptionNewf(py.ValueError, "math domain error")
	ERANGE = py.ExceptionNewf(py.OverflowError, "math range error")
)

// panic if ok is false
func assert(ok bool) {
	if !ok {
		panic("assertion failed")
	}
}

// isFinite is true if x is not Nan or +/-Inf
func isFinite(x float64) bool {
	return !(math.IsInf(x, 0) || math.IsNaN(x))
}

/*
   math_1 is used to wrap a libm function f that takes a float64
   arguments and returns a float64.

   The error reporting follows these rules, which are designed to do
   the right thing on C89/C99 platforms and IEEE 754/non IEEE 754
   platforms.

   - a NaN result from non-NaN inputs causes ValueError to be raised
   - an infinite result from finite inputs causes OverflowError to be
     raised if can_overflow is 1, or raises ValueError if can_overflow
     is 0.
   - if the result is finite and errno == EDOM then ValueError is
     raised
   - if the result is finite and nonzero and errno == ERANGE then
     OverflowError is raised

   The last rule is used to catch overflow on platforms which follow
   C89 but for which HUGE_VAL is not an infinity.

   For the majority of one-argument functions these rules are enough
   to ensure that Python's functions behave as specified in 'Annex F'
   of the C99 standard, with the 'invalid' and 'divide-by-zero'
   floating-point exceptions mapping to Python's ValueError and the
   'overflow' floating-point exception mapping to OverflowError.
   math_1 only works for functions that don't have singularities *and*
   the possibility of overflow; fortunately, that covers everything we
   care about right now.
*/
func math_1_to_whatever(arg py.Object, fn func(float64) float64, can_overflow bool) (float64, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return 0, err
	}
	r := fn(x)
	return checkResult(x, r, can_overflow)
}

// checkResult returns EDOM or ERANGE accordingly - see math_1_to_whatever for rules
func checkResult(x, r float64, can_overflow bool) (float64, error) {
	if math.IsNaN(r) && !math.IsNaN(x) {
		return 0, EDOM /* invalid arg */
	}
	if math.IsInf(r, 0) && isFinite(x) {
		if can_overflow {
			return 0, ERANGE /* overflow */
		} else {
			return 0, EDOM /* singularity */
		}
	}
	return r, nil
}

/* variant of math_1, to be used when the function being wrapped is known to
   set errno properly (that is, errno = EDOM for invalid or divide-by-zero,
   errno = ERANGE for overflow). */
func math_1a(arg py.Object, fn func(float64) float64) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	r := fn(x)
	return py.Float(r), nil
}

/*
   math_2 is used to wrap a libm function f that takes two float64
   arguments and returns a float64.

   The error reporting follows these rules, which are designed to do
   the right thing on C89/C99 platforms and IEEE 754/non IEEE 754
   platforms.

   - a NaN result from non-NaN inputs causes ValueError to be raised
   - an infinite result from finite inputs causes OverflowError to be
     raised.
   - if the result is finite and errno == EDOM then ValueError is
     raised
   - if the result is finite and nonzero and errno == ERANGE then
     OverflowError is raised

   The last rule is used to catch overflow on platforms which follow
   C89 but for which HUGE_VAL is not an infinity.

   For most two-argument functions (copysign, fmod, hypot, atan2)
   these rules are enough to ensure that Python's functions behave as
   specified in 'Annex F' of the C99 standard, with the 'invalid' and
   'divide-by-zero' floating-point exceptions mapping to Python's
   ValueError and the 'overflow' floating-point exception mapping to
   OverflowError.
*/
func math_1(arg py.Object, fn func(float64) float64, can_overflow bool) (py.Object, error) {
	f, err := math_1_to_whatever(arg, fn, can_overflow)
	if err != nil {
		return nil, err
	}
	return py.Float(f), nil
}

func math_1_to_int(arg py.Object, fn func(float64) float64, can_overflow bool) (py.Object, error) {
	f, err := math_1_to_whatever(arg, fn, can_overflow)
	if err != nil {
		return nil, err
	}
	return py.Float(f).M__int__()
}

func math_2(args py.Tuple, fn func(float64, float64) float64, fnname string) (py.Object, error) {
	var ox, oy py.Object
	var x, y, r float64
	err := py.UnpackTuple(args, nil, fnname, 2, 2, &ox, &oy)
	if err != nil {
		return nil, err
	}
	x, err = py.FloatAsFloat64(ox)
	if err != nil {
		return nil, err
	}
	y, err = py.FloatAsFloat64(oy)
	if err != nil {
		return nil, err
	}
	r = fn(x, y)
	if math.IsNaN(r) {
		if !math.IsNaN(x) && !math.IsNaN(y) {
			return nil, EDOM
		}
	} else if math.IsInf(r, 0) {
		if isFinite(x) && isFinite(y) {
			return nil, ERANGE
		}
	}
	return py.Float(r), nil
}

func math_acos(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Acos, false)
}

const math_acos_doc = "acos(x)\n\nReturn the arc cosine (measured in radians) of x."

func math_acosh(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Acosh, false)
}

const math_acosh_doc = "acosh(x)\n\nReturn the inverse hyperbolic cosine of x."

func math_asin(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Asin, false)
}

const math_asin_doc = "asin(x)\n\nReturn the arc sine (measured in radians) of x."

func math_asinh(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Asinh, false)
}

const math_asinh_doc = "asinh(x)\n\nReturn the inverse hyperbolic sine of x."

func math_atan(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Atan, false)
}

const math_atan_doc = "atan(x)\n\nReturn the arc tangent (measured in radians) of x."

func math_atan2(self py.Object, args py.Tuple) (py.Object, error) {
	return math_2(args, math.Atan2, "atan2")
}

const math_atan2_doc = "atan2(y, x)\n\nReturn the arc tangent (measured in radians) of y/x.\nUnlike atan(y/x), the signs of both x and y are considered."

func math_atanh(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Atanh, false)
}

const math_atanh_doc = "atanh(x)\n\nReturn the inverse hyperbolic tangent of x."

func math_ceil(self py.Object, number py.Object) (py.Object, error) {
	if I, ok := number.(py.I__ceil__); ok {
		return I.M__ceil__()
	} else if res, ok, err := py.TypeCall0(number, "__ceil__"); ok {
		return res, err
	}
	return math_1_to_int(number, math.Ceil, false)
}

const math_ceil_doc = `ceil(x)\n\nReturn the ceiling of x as an int.
This is the smallest integral value >= x.`

func math_copysign(self py.Object, args py.Tuple) (py.Object, error) {
	return math_2(args, math.Copysign, "copysign")
}

const math_copysign_doc = "copysign(x, y)\n\nReturn a float with the magnitude (absolute value) of x but the sign \nof y. On platforms that support signed zeros, copysign(1.0, -0.0) nreturns -1.0.\n"

func math_cos(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Cos, false)
}

const math_cos_doc = "cos(x)\n\nReturn the cosine of x (measured in radians)."

func math_cosh(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Cosh, true)
}

const math_cosh_doc = "cosh(x)\n\nReturn the hyperbolic cosine of x."

func math_erf(self py.Object, arg py.Object) (py.Object, error) {
	return math_1a(arg, math.Erf)
}

const math_erf_doc = "erf(x)\n\nError function at x."

func math_erfc(self py.Object, arg py.Object) (py.Object, error) {
	return math_1a(arg, math.Erfc)
}

const math_erfc_doc = "erfc(x)\n\nComplementary error function at x."

func math_exp(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Exp, true)
}

const math_exp_doc = "exp(x)\n\nReturn e raised to the power of x."

func math_expm1(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Expm1, true)
}

const math_expm1_doc = "expm1(x)\n\nReturn exp(x)-1.\nThis function avoids the loss of precision involved in the direct evaluation of exp(x)-1 for small x."

func math_fabs(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Abs, false)
}

const math_fabs_doc = "fabs(x)\n\nReturn the absolute value of the float x."

func math_floor(self py.Object, number py.Object) (py.Object, error) {
	if I, ok := number.(py.I__floor__); ok {
		return I.M__floor__()
	} else if res, ok, err := py.TypeCall0(number, "__floor__"); ok {
		return res, err
	}
	return math_1_to_int(number, math.Floor, false)
}

const math_floor_doc = `floor(x)\n\nReturn the floor of x as an int.
This is the largest integral value <= x.`

func math_gamma(self py.Object, arg py.Object) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	// If x is -ve integer...
	if x <= 0 && x == math.Floor(x) {
		return nil, EDOM
	}
	r := math.Gamma(x)
	res, err := checkResult(x, r, true)
	if err != nil {
		return nil, err
	}
	return py.Float(res), nil
}

const math_gamma_doc = "gamma(x)\n\nGamma function at x."

func math_lgamma(self py.Object, arg py.Object) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	// If x is -ve integer...
	if x <= 0 && x == math.Floor(x) {
		return nil, EDOM
	}
	r, _ := math.Lgamma(x)
	res, err := checkResult(x, r, true)
	if err != nil {
		return nil, err
	}
	return py.Float(res), nil
}

const math_lgamma_doc = "lgamma(x)\n\nNatural logarithm of absolute value of Gamma function at x."

func math_log1p(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Log1p, false)
}

const math_log1p_doc = "log1p(x)\n\nReturn the natural logarithm of 1+x (base e).\nThe result is computed in a way which is accurate for x near zero."

func math_sin(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Sin, false)
}

const math_sin_doc = "sin(x)\n\nReturn the sine of x (measured in radians)."

func math_sinh(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Sinh, true)
}

const math_sinh_doc = "sinh(x)\n\nReturn the hyperbolic sine of x."

func math_sqrt(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Sqrt, false)
}

const math_sqrt_doc = "sqrt(x)\n\nReturn the square root of x."

func math_tan(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Tan, false)
}

const math_tan_doc = "tan(x)\n\nReturn the tangent of x (measured in radians)."

func math_tanh(self py.Object, arg py.Object) (py.Object, error) {
	return math_1(arg, math.Tanh, false)
}

const math_tanh_doc = "tanh(x)\n\nReturn the hyperbolic tangent of x."

/* Precision summation function as msum() by Raymond Hettinger in
   <http://aspn.activestate.com/ASPN/Cookbook/Python/Recipe/393090>,
   enhanced with the exact partials sum and roundoff from Mark
   Dickinson's post at <http://bugs.python.org/file10357/msum4.py>.
   See those links for more details, proofs and other references.

   Note 1: IEEE 754R floating point semantics are assumed,
   but the current implementation does not re-establish special
   value semantics across iterations (i.e. handling -Inf + Inf).

   Note 2:  No provision is made for intermediate overflow handling;
   therefore, sum([1e+308, 1e-308, 1e+308]) returns 1e+308 while
   sum([1e+308, 1e+308, 1e-308]) raises an OverflowError due to the
   overflow of the first partial sum.

   Note 3: The intermediate values lo, yr, and hi are declared volatile so
   aggressive compilers won't algebraically reduce lo to always be exactly 0.0.
   Also, the volatile declaration forces the values to be stored in memory as
   regular float64s instead of extended long precision (80-bit) values.  This
   prevents float64 rounding because any addition or subtraction of two float64s
   can be resolved exactly into float64-sized hi and lo values.  As long as the
   hi value gets forced into a float64 before yr and lo are computed, the extra
   bits in downstream extended precision operations (x87 for example) will be
   exactly zero and therefore can be losslessly stored back into a float64,
   thereby preventing float64 rounding.

   Note 4: A similar implementation is in Modules/cmathmodule.c.
   Be sure to update both when making changes.

   Note 5: The signature of math.fsum() differs from builtins.sum()
   because the start argument doesn't make sense in the context of
   accurate summation.  Since the partials table is collapsed before
   returning a result, sum(seq2, start=sum(seq1)) may not equal the
   accurate result returned by sum(itertools.chain(seq1, seq2)).
*/

/* Full precision summation of a sequence of floats.

   def msum(iterable):
       partials = []  # sorted, non-overlapping partial sums
       for x in iterable:
           i = 0
           for y in partials:
               if abs(x) < abs(y):
                   x, y = y, x
               hi = x + y
               lo = y - (hi - x)
               if lo:
                   partials[i] = lo
                   i += 1
               x = hi
           partials[i:] = [x]
       return sum_exact(partials)

   Rounded x+y stored in hi with the roundoff stored in lo.  Together hi+lo
   are exactly equal to x+y.  The inner loop applies hi/lo summation to each
   partial so that the list of partial sums remains exact.

   Sum_exact() adds the partial sums exactly and correctly rounds the final
   result (using the round-half-to-even rule).  The items in partials remain
   non-zero, non-special, non-overlapping and strictly increasing in
   magnitude, but possibly not all having the same sign.

   Depends on IEEE 754 arithmetic guarantees and half-even rounding.
*/
func math_fsum(self py.Object, seq py.Object) (py.Object, error) {
	const NUM_PARTIALS = 32 /* initial partials array size, on stack */
	var i, j int
	var x, y float64
	p := make([]float64, 0, NUM_PARTIALS)
	special_sum := 0.0
	inf_sum := 0.0
	var hi, yr, lo float64

	iter, err := py.Iter(seq)
	if err != nil {
		return nil, err
	}

	for { /* for x in iterable */
		item, err := py.Next(iter)
		if err != nil {
			if py.IsException(py.StopIteration, err) {
				break
			}
			return nil, err
		}
		x, err = py.FloatAsFloat64(item)
		if err != nil {
			return nil, err
		}

		xsave := x
		for i, j = 0, 0; j < len(p); j++ { /* for y in partials */
			y = p[j]
			if math.Abs(x) < math.Abs(y) {
				x, y = y, x
			}
			hi = x + y
			yr = hi - x
			lo = y - yr
			if lo != 0.0 {
				p[i] = lo
				i++
			}
			x = hi
		}

		p = p[:i] /* ps[i:] = [x] */
		if x != 0.0 {
			if !isFinite(x) {
				/* a nonfinite x could arise either as
				   a result of intermediate overflow, or
				   as a result of a nan or inf in the
				   summands */
				if isFinite(xsave) {
					return nil, py.ExceptionNewf(py.OverflowError, "intermediate overflow in fsum")
				}
				if math.IsInf(xsave, 0) {
					inf_sum += xsave
				}
				special_sum += xsave
				/* reset partials */
				p = p[:0]
			} else {
				p = append(p, x)
			}
		}
	}

	if special_sum != 0.0 {
		if math.IsNaN(inf_sum) {
			return nil, py.ExceptionNewf(py.ValueError, "-inf + inf in fsum")
		} else {
			return py.Float(special_sum), nil
		}
	}

	hi = 0.0
	if len(p) > 0 {
		hi = p[len(p)-1]
		p = p[:len(p)-1]
		/* sum_exact(ps, hi) from the top, stop when the sum becomes
		   inexact. */
		for len(p) > 0 {
			x = hi
			y = p[len(p)-1]
			p = p[:len(p)-1]
			// assert(math.Abs(y) < math.Abs(x))
			hi = x + y
			yr = hi - x
			lo = y - yr
			if lo != 0.0 {
				break
			}
		}
		/* Make half-even rounding work across multiple partials.
		   Needed so that sum([1e-16, 1, 1e16]) will round-up the last
		   digit to two instead of down to zero (the 1e-16 makes the 1
		   slightly closer to two).  With a potential 1 ULP rounding
		   error fixed-up, math.fsum() can guarantee commutativity. */
		if len(p) > 0 && ((lo < 0.0 && p[len(p)-1] < 0.0) ||
			(lo > 0.0 && p[len(p)-1] > 0.0)) {
			y = lo * 2.0
			x = hi + y
			yr = x - hi
			if y == yr {
				hi = x
			}
		}
	}
	return py.Float(hi), nil
}

const math_fsum_doc = `fsum(iterable)

Return an accurate floating point sum of values in the iterable.
Assumes IEEE-754 floating point arithmetic.`

/* Return the smallest integer k such that n < 2**k, or 0 if n == 0.
 * Equivalent to floor(lg(x))+1.  Also equivalent to: bitwidth_of_type -
 * count_leading_zero_bits(x)
 */

/* XXX: This routine does more or less the same thing as
 * bits_in_digit() in Objects/longobject.c.  Someday it would be nice to
 * consolidate them.  On BSD, there's a library function called fls()
 * that we could use, and GCC provides __builtin_clz().
 */
func bit_length(n int64) int64 {
	var len int64 = 0
	for n != 0 {
		len++
		n >>= 1
	}
	return len
}

func count_set_bits(n int64) int64 {
	var count int64 = 0
	for n != 0 {
		count++
		n &= n - 1 /* clear least significant bit */
	}
	return count
}

/* Divide-and-conquer factorial algorithm
 *
 * Based on the formula and psuedo-code provided at:
 * http://www.luschny.de/math/factorial/binarysplitfact.html
 *
 * Faster algorithms exist, but they're more complicated and depend on
 * a fast prime factorization algorithm.
 *
 * Notes on the algorithm
 * ----------------------
 *
 * factorial(n) is written in the form 2**k * m, with m odd.  k and m are
 * computed separately, and then combined using a left shift.
 *
 * The function factorial_odd_part computes the odd part m (i.e., the greatest
 * odd divisor) of factorial(n), using the formula:
 *
 *   factorial_odd_part(n) =
 *
 *        product_{i >= 0} product_{0 < j <= n / 2**i, j odd} j
 *
 * Example: factorial_odd_part(20) =
 *
 *        (1) *
 *        (1) *
 *        (1 * 3 * 5) *
 *        (1 * 3 * 5 * 7 * 9)
 *        (1 * 3 * 5 * 7 * 9 * 11 * 13 * 15 * 17 * 19)
 *
 * Here i goes from large to small: the first term corresponds to i=4 (any
 * larger i gives an empty product), and the last term corresponds to i=0.
 * Each term can be computed from the last by multiplying by the extra odd
 * numbers required: e.g., to get from the penultimate term to the last one,
 * we multiply by (11 * 13 * 15 * 17 * 19).
 *
 * To see a hint of why this formula works, here are the same numbers as above
 * but with the even parts (i.e., the appropriate powers of 2) included.  For
 * each subterm in the product for i, we multiply that subterm by 2**i:
 *
 *   factorial(20) =
 *
 *        (16) *
 *        (8) *
 *        (4 * 12 * 20) *
 *        (2 * 6 * 10 * 14 * 18) *
 *        (1 * 3 * 5 * 7 * 9 * 11 * 13 * 15 * 17 * 19)
 *
 * The factorial_partial_product function computes the product of all odd j in
 * range(start, stop) for given start and stop.  It's used to compute the
 * partial products like (11 * 13 * 15 * 17 * 19) in the example above.  It
 * operates recursively, repeatedly splitting the range into two roughly equal
 * pieces until the subranges are small enough to be computed using only C
 * integer arithmetic.
 *
 * The two-valuation k (i.e., the exponent of the largest power of 2 dividing
 * the factorial) is computed independently in the main math_factorial
 * function.  By standard results, its value is:
 *
 *    two_valuation = n//2 + n//4 + n//8 + ....
 *
 * It can be shown (e.g., by complete induction on n) that two_valuation is
 * equal to n - count_set_bits(n), where count_set_bits(n) gives the number of
 * '1'-bits in the binary expansion of n.
 */

/* factorial_partial_product: Compute product(range(start, stop, 2)) using
 * divide and conquer.  Assumes start and stop are odd and stop > start.
 * max_bits must be >= bit_length(stop - 2). */
func factorial_partial_product(start int64, stop int64, max_bits int64) (py.Object, error) {
	var midpoint, num_operands int64
	var left, right, result py.Object
	var err error

	/* If the return value will fit an int64, then we can
	 * multiply in a tight, fast loop where each multiply is O(1).
	 * Compute an upper bound on the number of bits required to store
	 * the answer.
	 *
	 * Storing some integer z requires floor(lg(z))+1 bits, which is
	 * conveniently the value returned by bit_length(z).  The
	 * product x*y will require at most
	 * bit_length(x) + bit_length(y) bits to store, based
	 * on the idea that lg product = lg x + lg y.
	 *
	 * We know that stop - 2 is the largest number to be multiplied.  From
	 * there, we have: bit_length(answer) <= num_operands *
	 * bit_length(stop - 2)
	 */

	num_operands = (stop - start) / 2
	/* The "num_operands <= 63 check guards against the
	 * unlikely case of an overflow in num_operands * max_bits. */
	if num_operands <= 63 && num_operands*max_bits <= 63 {
		var j, total int64
		for total, j = start, start+2; j < stop; j += 2 {
			total *= j
		}
		return py.Int(total), nil
	}

	/* find midpoint of range(start, stop), rounded up to next odd number. */
	midpoint = (start + num_operands) | 1
	left, err = factorial_partial_product(start, midpoint, bit_length(midpoint-2))
	if err != nil {
		return nil, err
	}
	right, err = factorial_partial_product(midpoint, stop, max_bits)
	if err != nil {
		return nil, err
	}
	result, err = py.Mul(left, right)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/* factorial_odd_part:  compute the odd part of factorial(n). */
func factorial_odd_part(n int64) (py.Object, error) {
	var v, lower, upper int64
	var partial, tmp, inner, outer py.Object
	var err error

	inner = py.Int(1)
	outer = inner

	upper = 3
	for i := bit_length(n) - 2; i >= 0; i-- {
		v = n >> uint(i)
		if v <= 2 {
			continue
		}
		lower = upper
		/* (v + 1) | 1 = least odd integer strictly larger than n / 2**i */
		upper = (v + 1) | 1
		/* Here inner is the product of all odd integers j in the range (0,
		   n/2**(i+1)].  The factorial_partial_product call below gives the
		   product of all odd integers j in the range (n/2**(i+1), n/2**i]. */
		partial, err = factorial_partial_product(lower, upper, bit_length(upper-2))
		if err != nil {
			return nil, err
		}
		/* inner *= partial */
		tmp, err = py.Mul(inner, partial)
		if err != nil {
			return nil, err
		}
		inner = tmp
		/* Now inner is the product of all odd integers j in the range (0,
		   n/2**i], giving the inner product in the formula above. */

		/* outer *= inner; */
		tmp, err = py.Mul(outer, inner)
		if err != nil {
			return nil, err
		}
		outer = tmp
	}
	return outer, nil
}

/* Lookup table for small factorial values */
var smallFactorials = []py.Int{
	1, 1, 2, 6, 24, 120, 720, 5040, 40320,
	362880, 3628800, 39916800, 479001600,
	6227020800, 87178291200, 1307674368000,
	20922789888000, 355687428096000, 6402373705728000,
	121645100408832000, 2432902008176640000,
}

func math_factorial(self py.Object, arg py.Object) (py.Object, error) {
	var odd_part, two_valuation py.Object
	var x int64
	var err error

	if dxObj, err := py.FloatCheck(arg); err == nil {
		dx := float64(dxObj)
		if !isFinite(dx) || dx != math.Floor(dx) {
			return nil, py.ExceptionNewf(py.ValueError, "factorial() only accepts integral values")
		}
		arg = py.Int(dx)
	}
	x, err = py.MakeGoInt64(arg)
	if err != nil {
		return nil, err
	}
	if x < 0 {
		return nil, py.ExceptionNewf(py.ValueError, "factorial() not defined for negative values")
	}

	/* use lookup table if x is small */
	if x < int64(len(smallFactorials)) {
		return smallFactorials[x], nil
	}

	/* else express in the form odd_part * 2**two_valuation, and compute as
	   odd_part << two_valuation. */
	odd_part, err = factorial_odd_part(x)
	if err != nil {
		return nil, err
	}
	two_valuation = py.Int(x - count_set_bits(x))
	return py.Lshift(odd_part, two_valuation)
}

const math_factorial_doc = `factorial(x) -> Integral

Find x!. Raise a ValueError if x is negative or non-integral.`

func math_trunc(self py.Object, number py.Object) (py.Object, error) {
	if I, ok := number.(py.I__trunc__); ok {
		return I.M__trunc__()
	} else if res, ok, err := py.TypeCall0(number, "__trunc__"); ok {
		return res, err
	}
	return math_1_to_int(number, math.Trunc, false)
}

const math_trunc_doc = `trunc(x:Real) -> Integral

Truncates x to the nearest Integral toward 0. Uses the __trunc__ magic method.`

func math_frexp(self py.Object, arg py.Object) (py.Object, error) {
	var exp int
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	/* deal with special cases directly, to sidestep platform differences */
	if math.IsNaN(x) || math.IsInf(x, 0) || x == 0 {
		exp = 0
	} else {
		x, exp = math.Frexp(x)
	}
	return py.Tuple{py.Float(x), py.Int(exp)}, nil
}

const math_frexp_doc = `frexp(x)

Return the mantissa and exponent of x, as pair (m, e).
m is a float and e is an int, such that x = m * 2.**e.
If x is 0, m and e are both 0.  Else 0.5 <= abs(m) < 1.0.`

func math_ldexp(self py.Object, args py.Tuple) (py.Object, error) {
	var x, r float64
	var xObj py.Object
	var expObj py.Object
	var exp int
	err := py.UnpackTuple(args, nil, "ldexp", 2, 2, &xObj, &expObj)
	if err != nil {
		return nil, err
	}
	x, err = py.FloatAsFloat64(xObj)
	if err != nil {
		return nil, err
	}
	exp, err = py.MakeGoInt(expObj)
	if err != nil {
		// on overflow, replace exponent with either LONG_MAX
		// or LONG_MIN, depending on the sign.
		expInt, err := py.MakeInt(expObj)
		if err != nil {
			return nil, err
		}
		lt, err := py.Lt(expInt, py.Int(0))
		if err != nil {
			return nil, err
		}
		if lt == py.True {
			exp = py.GoIntMin
		} else {
			exp = py.GoIntMax
		}
	}
	if x == 0. || !isFinite(x) {
		/* NaNs, zeros and infinities are returned unchanged */
		r = x
	} else if exp > math.MaxInt16 {
		/* overflow */
		// r = math.Copysign(math.Inf(1), x)
		return nil, ERANGE
	} else if exp < math.MinInt16 {
		/* underflow to +-0 */
		r = math.Copysign(0., x)
	} else {
		r = math.Ldexp(x, exp)
		if math.IsInf(r, 0) {
			return nil, ERANGE
		}
	}
	return py.Float(r), nil
}

const math_ldexp_doc = `ldexp(x, i)

Return x * (2**i).`

func math_modf(self py.Object, arg py.Object) (py.Object, error) {
	var y float64
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	/* some platforms don't do the right thing for NaNs and
	   infinities, so we take care of special cases directly. */
	if !isFinite(x) {
		if math.IsInf(x, 0) {
			return py.Tuple{py.Float(math.Copysign(0., x)), py.Float(x)}, nil
		} else if math.IsNaN(x) {
			return py.Tuple{py.Float(x), py.Float(x)}, nil
		}
	}

	y, x = math.Modf(x)
	return py.Tuple{py.Float(x), py.Float(y)}, nil
}

const math_modf_doc = `modf(x)

Return the fractional and integer parts of x.  Both results carry the sign
of x and are floats.`

/* A decent logarithm is easy to compute even for huge ints, but libm can't
   do that by itself -- loghelper can.  func is log or log10, and name is
   "log" or "log10".  Note that overflow of the result isn't possible: an int
   can contain no more than INT_MAX * SHIFT bits, so has value certainly less
   than 2**(2**64 * 2**16) == 2**2**80, and log2 of that is 2**80, which is
   small enough to fit in an IEEE single.  log and log10 are even smaller.
   However, intermediate overflow is possible for an int if the number of bits
   in that int is larger than PY_SSIZE_T_MAX. */
func loghelper(arg py.Object, fn func(float64) float64, fnname string) (py.Object, error) {
	/* If it is int, do it ourselves. */
	if xBig, err := py.BigIntCheck(arg); err == nil {
		var x, result float64
		var e int

		/* Negative or zero inputs give a ValueError. */
		if (*big.Int)(xBig).Sign() <= 0 {
			return nil, EDOM
		}

		xf, err := xBig.Float()
		x = float64(xf)
		if err != nil {
			/* Here the conversion to float64 overflowed, but it's possible
			   to compute the log anyway. */
			x, e = xBig.Frexp()
			/* Value is ~= x * 2**e, so the log ~= log(x) + log(2) * e. */
			result = fn(x) + fn(2.0)*float64(e)
		} else {
			/* Successfully converted x to a float64. */
			result = fn(x)
		}
		return py.Float(result), nil
	}
	/* Else let libm handle it by itself. */
	return math_1(arg, fn, false)
}

func math_log(self py.Object, args py.Tuple) (py.Object, error) {
	var arg py.Object
	var base py.Object = py.Float(math.E)

	err := py.UnpackTuple(args, nil, "log", 1, 2, &arg, &base)
	if err != nil {
		return nil, err
	}

	num, err := loghelper(arg, math.Log, "log")
	if err != nil {
		return nil, err
	}

	den, err := loghelper(base, math.Log, "log")
	if err != nil {
		return nil, err
	}

	return py.TrueDiv(num, den)
}

const math_log_doc = `log(x[, base])

Return the logarithm of x to the given base.
If the base not specified, returns the natural logarithm (base e) of x.`

func math_log2(self py.Object, arg py.Object) (py.Object, error) {
	return loghelper(arg, math.Log2, "log2")
}

const math_log2_doc = `log2(x)
Return the base 2 logarithm of x.`

func math_log10(self py.Object, arg py.Object) (py.Object, error) {
	return loghelper(arg, math.Log10, "log10")
}

const math_log10_doc = `log10(x)
Return the base 10 logarithm of x.`

func math_fmod(self py.Object, args py.Tuple) (py.Object, error) {
	var ox, oy py.Object
	var r, x, y float64
	err := py.UnpackTuple(args, nil, "fmod", 2, 2, &ox, &oy)
	if err != nil {
		return nil, err
	}
	x, err = py.FloatAsFloat64(ox)
	if err != nil {
		return nil, err
	}
	y, err = py.FloatAsFloat64(oy)
	if err != nil {
		return nil, err
	}
	/* fmod(x, +/-Inf) returns x for finite x. */
	if math.IsInf(y, 0) && isFinite(x) {
		return py.Float(x), nil
	}
	r = math.Mod(x, y)
	if math.IsNaN(r) {
		if !math.IsNaN(x) && !math.IsNaN(y) {
			return nil, EDOM
		}
	}
	return py.Float(r), nil
}

const math_fmod_doc = `fmod(x, y)

Return fmod(x, y), according to platform C.  x % y may differ.`

func math_hypot(self py.Object, args py.Tuple) (py.Object, error) {
	var ox, oy py.Object
	var r, x, y float64
	err := py.UnpackTuple(args, nil, "hypot", 2, 2, &ox, &oy)
	if err != nil {
		return nil, err
	}
	x, err = py.FloatAsFloat64(ox)
	if err != nil {
		return nil, err
	}
	y, err = py.FloatAsFloat64(oy)
	if err != nil {
		return nil, err
	}
	/* hypot(x, +/-Inf) returns Inf, even if x is a NaN. */
	if math.IsInf(x, 0) {
		return py.Float(math.Abs(x)), nil
	}
	if math.IsInf(y, 0) {
		return py.Float(math.Abs(y)), nil
	}
	r = math.Hypot(x, y)
	if math.IsNaN(r) {
		if !math.IsNaN(x) && !math.IsNaN(y) {
			return nil, EDOM
		}
	} else if math.IsInf(r, 0) {
		if isFinite(x) && isFinite(y) {
			return nil, ERANGE
		}
	}
	return py.Float(r), nil
}

const math_hypot_doc = `hypot(x, y)
Return the Euclidean distance, sqrt(x*x + y*y).`

/* pow can't use math_2, but needs its own wrapper: the problem is
   that an infinite result can arise either as a result of overflow
   (in which case OverflowError should be raised) or as a result of
   e.g. 0.**-5. (for which ValueError needs to be raised.)
*/
func math_pow(self py.Object, args py.Tuple) (py.Object, error) {
	var ox, oy py.Object
	var r, x, y float64

	err := py.UnpackTuple(args, nil, "pow", 2, 2, &ox, &oy)
	if err != nil {
		return nil, err
	}
	x, err = py.FloatAsFloat64(ox)
	if err != nil {
		return nil, err
	}
	y, err = py.FloatAsFloat64(oy)
	if err != nil {
		return nil, err
	}

	/* deal directly with IEEE specials, to cope with problems on various
	   platforms whose semantics don't exactly match C99 */
	r = 0. /* silence compiler warning */
	if !isFinite(x) || !isFinite(y) {
		if math.IsNaN(x) {
			if y == 0. {
				r = 1.
			} else {
				r = x
			} /* NaN**0 = 1 */
		} else if math.IsNaN(y) {
			if x == 1. {
				r = 1.
			} else {
				r = y
			} /* 1**NaN = 1 */
		} else if math.IsInf(x, 0) {
			odd_y := isFinite(y) && math.Mod(math.Abs(y), 2.0) == 1.0
			if y > 0. {
				if odd_y {
					r = x
				} else {
					r = math.Abs(x)
				}
			} else if y == 0. {
				r = 1.
			} else { /* y < 0. */
				if odd_y {
					r = math.Copysign(0., x)
				} else {
					r = 0.
				}
			}
		} else if math.IsInf(y, 0) {
			if math.Abs(x) == 1.0 {
				r = 1.
			} else if y > 0. && math.Abs(x) > 1.0 {
				r = y
			} else if y < 0. && math.Abs(x) < 1.0 {
				r = -y       /* result is +inf */
				if x == 0. { /* 0**-inf: divide-by-zero */
					return nil, EDOM
				}
			} else {
				r = 0.
			}
		}
	} else {
		// Go returns Inf rather than NaN for -ve, so pick this off early
		if x == 0 && y < 0 {
			return nil, EDOM
		}
		/* let libm handle finite**finite */
		r = math.Pow(x, y)
		/* a NaN result should arise only from (-ve)**(finite
		   non-integer); in this case we want to raise ValueError. */
		if !isFinite(r) {
			if math.IsNaN(r) {
				return nil, EDOM
			} else if math.IsInf(r, 0) {
				/*
				   an infinite result here arises either from:
				   (A) (+/-0.)**negative (-> divide-by-zero)
				   (B) overflow of x**y with x and y finite
				*/
				if x != 0. {
					return nil, ERANGE
				}
			}
		}
	}
	return py.Float(r), nil
}

const math_pow_doc = `pow(x, y)

Return x**y (x to the power of y).`

const (
	degToRad = math.Pi / 180.0
	radToDeg = 180.0 / math.Pi
)

func math_degrees(self py.Object, arg py.Object) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	return py.Float(x * radToDeg), nil
}

const math_degrees_doc = `degrees(x)

Convert angle x from radians to degrees.`

func math_radians(self py.Object, arg py.Object) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	return py.Float(x * degToRad), nil
}

const math_radians_doc = `radians(x)

Convert angle x from degrees to radians.`

func math_isfinite(self py.Object, arg py.Object) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	return py.Bool(isFinite(x)), nil
}

const math_isfinite_doc = `isfinite(x) -> bool

Return True if x is neither an infinity nor a NaN, and False otherwise.`

func math_isnan(self py.Object, arg py.Object) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	return py.Bool(math.IsNaN(x)), nil
}

const math_isnan_doc = `isnan(x) -> bool

Return True if x is a NaN (not a number), and False otherwise.`

func math_isinf(self py.Object, arg py.Object) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	return py.Bool(math.IsInf(x, 0)), nil
}

const math_isinf_doc = `isinf(x) -> bool

Return True if x is a positive or negative infinity, and False otherwise.`

func math_to_ulps(self py.Object, arg py.Object) (py.Object, error) {
	x, err := py.FloatAsFloat64(arg)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.LittleEndian, x)
	if err != nil {
		return py.ExceptionNewf(py.ValueError, "to_ulps: binary.Write failed: %v", err), nil
	}
	var n int64
	err = binary.Read(&buf, binary.LittleEndian, &n)
	if err != nil {
		return py.ExceptionNewf(py.ValueError, "to_ulps: binary.Read failed: %v", err), nil
	}
	if n < 0 {
		n = -n + -(1 << 63)
	}
	return py.Int(n), nil
}

const math_to_ulps_doc = `to_ulps(x) -> int

Convert a non-NaN float x to an integer, the unit of least precision,
in such a way that adjacent floats are converted to adjacent integers.
Then abs(ulps(x) - ulps(y)) gives the difference in ulps between two
floats.

The results from this function will only make sense on platforms where
float64 are represented in IEEE 754 binary64 format.`

const math_doc = `This module is always available.  It provides access to the
mathematical functions defined by the C standard.`

// Initialise the module
func init() {
	methods := []*py.Method{
		py.MustNewMethod("acos", math_acos, 0, math_acos_doc),
		py.MustNewMethod("acosh", math_acosh, 0, math_acosh_doc),
		py.MustNewMethod("asin", math_asin, 0, math_asin_doc),
		py.MustNewMethod("asinh", math_asinh, 0, math_asinh_doc),
		py.MustNewMethod("atan", math_atan, 0, math_atan_doc),
		py.MustNewMethod("atan2", math_atan2, 0, math_atan2_doc),
		py.MustNewMethod("atanh", math_atanh, 0, math_atanh_doc),
		py.MustNewMethod("ceil", math_ceil, 0, math_ceil_doc),
		py.MustNewMethod("copysign", math_copysign, 0, math_copysign_doc),
		py.MustNewMethod("cos", math_cos, 0, math_cos_doc),
		py.MustNewMethod("cosh", math_cosh, 0, math_cosh_doc),
		py.MustNewMethod("degrees", math_degrees, 0, math_degrees_doc),
		py.MustNewMethod("erf", math_erf, 0, math_erf_doc),
		py.MustNewMethod("erfc", math_erfc, 0, math_erfc_doc),
		py.MustNewMethod("exp", math_exp, 0, math_exp_doc),
		py.MustNewMethod("expm1", math_expm1, 0, math_expm1_doc),
		py.MustNewMethod("fabs", math_fabs, 0, math_fabs_doc),
		py.MustNewMethod("factorial", math_factorial, 0, math_factorial_doc),
		py.MustNewMethod("floor", math_floor, 0, math_floor_doc),
		py.MustNewMethod("fmod", math_fmod, 0, math_fmod_doc),
		py.MustNewMethod("frexp", math_frexp, 0, math_frexp_doc),
		py.MustNewMethod("fsum", math_fsum, 0, math_fsum_doc),
		py.MustNewMethod("gamma", math_gamma, 0, math_gamma_doc),
		py.MustNewMethod("hypot", math_hypot, 0, math_hypot_doc),
		py.MustNewMethod("isfinite", math_isfinite, 0, math_isfinite_doc),
		py.MustNewMethod("isinf", math_isinf, 0, math_isinf_doc),
		py.MustNewMethod("isnan", math_isnan, 0, math_isnan_doc),
		py.MustNewMethod("ldexp", math_ldexp, 0, math_ldexp_doc),
		py.MustNewMethod("lgamma", math_lgamma, 0, math_lgamma_doc),
		py.MustNewMethod("log", math_log, 0, math_log_doc),
		py.MustNewMethod("log1p", math_log1p, 0, math_log1p_doc),
		py.MustNewMethod("log10", math_log10, 0, math_log10_doc),
		py.MustNewMethod("log2", math_log2, 0, math_log2_doc),
		py.MustNewMethod("modf", math_modf, 0, math_modf_doc),
		py.MustNewMethod("pow", math_pow, 0, math_pow_doc),
		py.MustNewMethod("radians", math_radians, 0, math_radians_doc),
		py.MustNewMethod("sin", math_sin, 0, math_sin_doc),
		py.MustNewMethod("sinh", math_sinh, 0, math_sinh_doc),
		py.MustNewMethod("sqrt", math_sqrt, 0, math_sqrt_doc),
		py.MustNewMethod("tan", math_tan, 0, math_tan_doc),
		py.MustNewMethod("tanh", math_tanh, 0, math_tanh_doc),
		py.MustNewMethod("trunc", math_trunc, 0, math_trunc_doc),
		py.MustNewMethod("to_ulps", math_to_ulps, 0, math_to_ulps_doc),
	}

	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name:  "math",
			Doc:   math_doc,
			Flags: py.ShareModule,
		},
		Methods: methods,
		Globals: py.StringDict{
			"pi": py.Float(math.Pi),
			"e":  py.Float(math.E),
		},
	})
	
}
