doc="str"
assert str(()) == "()"
assert str((1,2,3)) == "(1, 2, 3)"
assert str((1,(2,3),4)) == "(1, (2, 3), 4)"
assert str(("1",(2.5,17,()))) == "('1', (2.5, 17, ()))"

doc="repr"
assert repr(()) == "()"
assert repr((1,2,3)) == "(1, 2, 3)"
assert repr((1,(2,3),4)) == "(1, (2, 3), 4)"
assert repr(("1",(2.5,17,()))) == "('1', (2.5, 17, ()))"

doc="finished"
