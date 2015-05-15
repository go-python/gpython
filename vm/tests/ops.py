#!/usr/bin/env python3.4

_2 = 2
_10 = 10
_11 = 11
_100 = 100
a=[0,1,2,3,4]

# Unary Ops
assert +_2 == 2
assert -_2 == -2
assert not _2 is False
assert ~_2 == -3

# Binary Ops
assert _2**_10 == 1024
assert _2*_10 == 20
assert _10//_2 == 5
assert _11//_2 == 5
assert _10/_2 == 5.0
assert _11/_2 == 5.5
assert _10 % _2 == 0
assert _11 % _2 == 1
assert _10 + _2 == 12
assert _10 - _2 == 8
assert a[1] == 1
assert a[4] == 4
assert _2 << _10 == 2048
assert _100 >> 2 == 25
assert _10 & _2 == 2
assert _100 | _2 == 102
assert _10 ^ _2 == 8

# Inplace Ops
a = _2
a **= _10 
assert a == 1024
a = _2
a *= _10 
assert a == 20
a = _10
a //= _2 
assert a == 5
a = _11
a //= _2 
assert a == 5
a = _10
a /= _2 
assert a == 5.0
a = _11
a /= _2 
assert a == 5.5
a = _10
a %= _2 
assert a == 0
a = _11
a %= _2 
assert a == 1
a = _10
a += _2 
assert a == 12
a = _10
a -= _2 
assert a == 8
a = _2
a <<= _10 
assert a == 2048
a = _100
a >>= 2 
assert a == 25
a = _10
a &= _2 
assert a == 2
a = _100
a |= _2 
assert a == 102
a = _10
a ^= _2 
assert a == 8

# Comparison
assert _2 < _10
assert _2 <= _10
assert _2 <= _2
assert _2 == _2
assert _2 != _10
assert _10 > _2
assert _10 >= _2
assert _2 >= _2
# FIXME in
# FIXME not in
assert True is True
assert True is not False
# FIXME EXC_MATCH
