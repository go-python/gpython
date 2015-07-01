from libtest import assertRaises

doc="float"
assert float("1.5") == 1.5
assert float(" 1.5") == 1.5
assert float(" 1.5 ") == 1.5
assert float("1.5 ") == 1.5
assert float("-1E9") == -1E9
assert float("1E400") == float("inf")
assert float(" -1E400") == float("-inf")
assertRaises(ValueError, float, "1 E200")

doc="finished"
