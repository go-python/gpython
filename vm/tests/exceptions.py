# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Test exceptions

doc="except"
ok = False
try:
    raise ValueError
except ValueError:
    ok = True
assert ok

doc="except as"
ok = False
try:
    raise ValueError
except ValueError as e:
    ok = True
assert ok

doc="except (a, b) as"
ok = False
try:
    raise ValueError
except (IOError, ValueError) as e:
    ok = True
assert ok

doc="raise exception instance"
ok = False
try:
    raise ValueError("Potato")
except (IOError, ValueError) as e:
    ok = True
assert ok

doc="exception hierarchy"
# FIXME doesn't work because IsSubtype is broken ValueError.IsSubtype(Exception) == false
# ok = False
# try:
#     raise ValueError("potato")
# except Exception:
#     ok = True
# assert ok

doc="exception match"
ok = False
try:
    raise ValueError
except IOError:
    assert False, "Not expecting IO Error"
except ValueError:
    ok = True
assert ok

doc="no exception"
ok = False
try:
    pass
except ValueError:
    assert False, "Not expecting ValueError"
else:
    ok = True
assert ok

doc="nested"
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

doc="nested #2"
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

doc="re-raise"
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

doc="try/finally"
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

doc="try/finally #2"
ok1 = False
ok2 = False
try:
    try:
        raise ValueError()
    finally:
        ok1 = True
except ValueError:
    ok2 = True
else:
    assert False, "Expecting ValueError (outer)"
assert ok1 and ok2

doc="internal exception"
ok = False
try:
    print(1/0)
except ZeroDivisionError:
    ok = True
assert ok

doc = "raise in else"
ok = False
try:
    try:
        pass
    except NameError as e:
        pass
    else:
        raise ValueError
except ValueError:
    ok = True
assert ok, "ValueError not raised"

doc = "finished"
