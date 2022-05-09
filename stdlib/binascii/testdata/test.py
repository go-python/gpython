# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import binascii

print("globals:")
for name in ("Error", "Incomplete"):
    v = getattr(binascii, name)
    print("\nbinascii.%s:\n%s" % (name,repr(v)))

def assertEqual(x, y):
    assert x == y, "got: %s, want: %s" % (repr(x), repr(y))

## base64
assertEqual(binascii.b2a_base64(b'hello world!'), b'aGVsbG8gd29ybGQh\n')
assertEqual(binascii.b2a_base64(b'hello\nworld!'), b'aGVsbG8Kd29ybGQh\n')
assertEqual(binascii.b2a_base64(b'hello world!', newline=False), b'aGVsbG8gd29ybGQh')
assertEqual(binascii.b2a_base64(b'hello\nworld!', newline=False), b'aGVsbG8Kd29ybGQh')
assertEqual(binascii.a2b_base64("aGVsbG8gd29ybGQh\n"), b'hello world!')

try:
    binascii.b2a_base64("string")
    print("expected an exception")
except TypeError as e:
    print("expected an exception:", e)
    pass

## crc32
assertEqual(binascii.crc32(b'hello world!'), 62177901)
assertEqual(binascii.crc32(b'hello world!', 0), 62177901)
assertEqual(binascii.crc32(b'hello world!', 42), 4055036404)

## hex
assertEqual(binascii.b2a_hex(b'hello world!'), b'68656c6c6f20776f726c6421')
assertEqual(binascii.a2b_hex(b'68656c6c6f20776f726c6421'), b'hello world!')
assertEqual(binascii.hexlify(b'hello world!'), b'68656c6c6f20776f726c6421')
assertEqual(binascii.unhexlify(b'68656c6c6f20776f726c6421'), b'hello world!')

try:
    binascii.a2b_hex(b'123')
    print("expected an exception")
except binascii.Error as e:
    print("expected an exception:",e)
    pass

try:
    binascii.a2b_hex(b'hell')
    print("expected an exception")
except binascii.Error as e:
    print("expected an exception:",e)
    pass

print("done.")
