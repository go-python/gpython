# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

try:
    from math import to_ulps
except ImportError:
    import struct
    def to_ulps(x):
        """Convert a non-NaN float x to an integer, in such a way that
        adjacent floats are converted to adjacent integers.  Then
        abs(ulps(x) - ulps(y)) gives the difference in ulps between two
        floats.

        The results from this function will only make sense on platforms
        where C doubles are represented in IEEE 754 binary64 format.

        """
        n = struct.unpack('<q', struct.pack('<d', x))[0]
        if n < 0:
            n = -(n+2**63)
        return n

def ulps_check(what, want, got, ulps=20):
    """Given non-NaN floats `want` and `got`,
    check that they're equal to within the given number of ulps.

    Returns None on success and an error message on failure."""
    ulps_error = to_ulps(got) - to_ulps(want)
    if abs(ulps_error) <= ulps:
        return None
    raise AssertionError("%s: want %g got %g: error = %d ulps; permitted error = %s ulps" % (what, want, got, ulps_error, ulps))
