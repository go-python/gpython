#!/usr/bin/env python3.4

class C:
    attr1 = 42
    def __init__(self):
        self.attr2 = 43
        self.attr3 = 44
c = C()

# Test LOAD_ATTR
assert c.attr1 == 42
assert C.attr1 == 42
assert c.attr2 == 43
assert c.attr3 == 44

# Test DELETE_ATTR
del c.attr3

# FIXME - exception handling broken
# ok = False
# try:
#     c.attr3
# except AttributeError:
#     ok = True
# assert ok

# Test STORE_ATTR
c.attr1 = 100
c.attr2 = 101
c.attr3 = 102
assert c.attr1 == 100
assert C.attr1 == 42
assert c.attr2 == 101
assert c.attr3 == 102

# End with this
finished = True
