#!/usr/bin/env python3
"""
Read in grammar_test.go, and re-write the tests section
"""

import sys
import ast

inp = [
    ("", "exec"),
    ("pass", "exec"),
    ("()", "eval"),
    ("()", "exec"),
    ("[ ]", "exec"),
    ("True\n", "eval"),
    ("False\n", "eval"),
    ("None\n", "eval"),
    ("...", "eval"),
    ("abc123", "eval"),
    ('"abc"', "eval"),
    ('"abc" """123"""', "eval"),
    ("b'abc'", "eval"),
    ("b'abc' b'''123'''", "eval"),
    ("1234", "eval"),
    ("0x1234", "eval"),
    ("12.34", "eval"),
    ("{ }", "eval"),
]

def dump(source, mode):
    """Dump source after parsing with mode"""
    a = ast.parse(source, mode=mode)
    return ast.dump(a, annotate_fields=True, include_attributes=False)

def escape(x):
    """Encode strings with backslashes for python/go"""
    return x.replace('\\', "\\\\").replace('"', r'\"').replace("\n", r'\n').replace("\t", r'\t')

def main():
    """Read in grammar_test.go, and re-write the tests section"""
    path = "grammar_test.go"
    with open(path) as f:
        grammar_test = f.read()
    lines = grammar_test.split("\n")
    while lines[-1] == "":
        lines = lines[:-1]
    out = []
    in_tests = False
    for line in lines:
        if "START TESTS" in line:
            out.append(line)
            for source, mode in inp:
                out.append('\t\t{"%s", "%s", "%s"},' % (escape(source), mode, escape(dump(source, mode))))
            in_tests = True
        elif "END TESTS" in line:
            in_tests = False
        if not in_tests:
            out.append(line)
    print("Rewriting %s" % path)
    with open(path, "w") as f:
        f.write("\n".join(out))
        f.write("\n")

if __name__ == "__main__":
    main()
