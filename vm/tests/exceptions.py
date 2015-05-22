#!/usr/bin/env python3.4

# Test exceptions

# Straight forward
ok = False
try:
    raise ValueError
except ValueError:
    ok = True
assert ok

ok = False
try:
    raise ValueError
except ValueError as e:
    ok = True
assert ok

ok = False
try:
    raise ValueError
except ValueError as e:
    ok = True
assert ok

ok = False
try:
    raise ValueError
except (IOError, ValueError) as e:
    ok = True
assert ok

ok = False
try:
    raise ValueError("Potato")
except (IOError, ValueError) as e:
    ok = True
assert ok

# hierarchy
# FIXME doesn't work because IsSubtype is broken ValueError.IsSubtype(Exception) == false
# ok = False
# try:
#     raise ValueError("potato")
# except Exception:
#     ok = True
# assert ok

ok = False
try:
    raise ValueError
except IOError:
    assert False, "Not expecting IO Error"
except ValueError:
    ok = True
assert ok

# no exception
ok = False
try:
    pass
except ValueError:
    assert False, "Not expecting ValueError"
else:
    ok = True
assert ok

# nested
ok = False
try:
    try:
        raise ValueError("potato")
    except IOError as e:
        assert False, "Not expecting IOError"
    else:
        assert False, "Expecting ValueError"
except ValueError:
    ok = True
else:
    assert False, "Expecting ValueError (outer)"
assert ok

ok1 = False
ok2 = False
try:
    try:
        raise IOError("potato")
    except IOError as e:
        ok1 = True
    else:
        assert False, "Expecting ValueError"
except ValueError:
    assert False, "Expecting IOError"
except IOError:
    assert False, "Expecting IOError"
else:
    ok2 = True
assert ok

# re-raise
ok1 = False
ok2 = False
try:
    try:
        raise ValueError("potato")
    except ValueError as e:
        ok2 = True
        raise
    else:
        assert False, "Expecting ValueError (inner)"
except ValueError:
    ok1 = True
else:
    assert False, "Expecting ValueError (outer)"
assert ok1 and ok2

# try/finally
ok1 = False
ok2 = False
ok3 = False
try:
    try:
        ok1 = True
    finally:
        ok2 = True
except ValueError:
    assert False, "Not expecting ValueError (outer)"
else:
    ok3 = True
assert ok1 and ok2 and ok3

# FIXME
# ok1 = False
# ok2 = False
# try:
#     try:
#         raise ValueError()
#     finally:
#         ok1 = True
# except ValueError:
#     ok2 = True
# else:
#     assert False, "Expecting ValueError (outer)"
# assert ok1 and ok2

# FIXME - exeption not being caught
# ok = False
# try:
#     print(1/0)
# except ZeroDivisionError:
#     ok = True
# assert ok

# End with this
finished = True
