#!/usr/bin/env python3.4

# List comprehensions
A = [1,2,3,4]
B = [ 2*a for a in A ]
assert tuple(B) == tuple([2,4,6,8])
B = [ 2*a for a in A if a != 2]
assert tuple(B) == tuple([2,6,8])

# Generator expressions
A = (1,2,3,4)
B = ( 2*a for a in A )
assert tuple(B) == (2,4,6,8)
B = [ 2*a for a in A if a != 2]
assert tuple(B) == (2,6,8)

# Set comprehensions
A = {1,2,3,4}
B = { 2*a for a in A }
assert B == {2,4,6,8}
B = { 2*a for a in A if a != 2}
assert B == {2,6,8}

# Dict comprehensions
A = {"a":1, "b":2, "c":3}
B = { k:k for k in ("a","b","c") }
assert B["b"] == "b"

# End with this
finished = True
