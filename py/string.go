// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// String objects
//
// Note that string objects in Python are arrays of unicode
// characters.  However we are using the native Go string which is
// UTF-8 encoded.  This makes very little difference most of the time,
// but care is needed when indexing, slicing or iterating through
// strings.

package py

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type String string

var StringType = ObjectType.NewType("str",
	`str(object='') -> str
str(bytes_or_buffer[, encoding[, errors]]) -> str

Create a new string object from the given object. If encoding or
errors is specified, then the object must expose a data buffer
that will be decoded using the given encoding and error handler.
Otherwise, returns the result of object.__str__() (if defined)
or repr(object).
encoding defaults to sys.getdefaultencoding().
errors defaults to 'strict'.`, StrNew, nil)

// Escape the py.String
func StringEscape(a String, ascii bool) string {
	s := string(a)
	var out bytes.Buffer
	quote := '\''
	if strings.ContainsRune(s, '\'') && !strings.ContainsRune(s, '"') {
		quote = '"'
	}
	if !ascii {
		out.WriteRune(quote)
	}
	for _, c := range s {
		switch {
		case c < 0x20:
			switch c {
			case '\t':
				out.WriteString(`\t`)
			case '\n':
				out.WriteString(`\n`)
			case '\r':
				out.WriteString(`\r`)
			default:
				fmt.Fprintf(&out, `\x%02x`, c)
			}
		case !ascii && c < 0x7F:
			if c == '\\' || (quote == '\'' && c == '\'') || (quote == '"' && c == '"') {
				out.WriteRune('\\')
			}
			out.WriteRune(c)
		case c < 0x100:
			if ascii || strconv.IsPrint(c) {
				out.WriteRune(c)
			} else {
				fmt.Fprintf(&out, "\\x%02x", c)
			}
		case c < 0x10000:
			if !ascii && strconv.IsPrint(c) {
				out.WriteRune(c)
			} else {
				fmt.Fprintf(&out, "\\u%04x", c)
			}
		default:
			if !ascii && strconv.IsPrint(c) {
				out.WriteRune(c)
			} else {
				fmt.Fprintf(&out, "\\U%08x", c)
			}
		}
	}
	if !ascii {
		out.WriteRune(quote)
	}
	return out.String()
}

// standard golang strings.Fields doesn't have a 'first N' argument
func fieldsN(s string, n int) []string {
	out := []string{}
	cur := []rune{}
	for _, c := range s {
		//until we have covered the first N elements, multiple white-spaces are 'merged'
		if n < 0 || len(out) < n {
			if unicode.IsSpace(c) {
				if len(cur) > 0 {
					out = append(out, string(cur))
					cur = []rune{}
				}
			} else {
				cur = append(cur, c)
			}
			//until we see the next letter, after collecting the first N fields, continue to merge whitespaces
		} else if len(out) == n && len(cur) == 0 {
			if !unicode.IsSpace(c) {
				cur = append(cur, c)
			}
			//now that enough words have been collected, just copy into the last element
		} else {
			cur = append(cur, c)
		}
	}
	if len(cur) > 0 {
		out = append(out, string(cur))
	}
	return out
}

func init() {
	StringType.Dict["endswith"] = MustNewMethod("endswith", func(self Object, args Tuple) (Object, error) {
		selfStr := string(self.(String))
		suffix := []string{}
		if len(args) > 0 {
			if s, ok := args[0].(String); ok {
				suffix = append(suffix, string(s))
			} else if s, ok := args[0].(Tuple); ok {
				for _, t := range s {
					if v, ok := t.(String); ok {
						suffix = append(suffix, string(v))
					}
				}
			} else {
				return nil, ExceptionNewf(TypeError, "endswith first arg must be str, unicode, or tuple, not %s", args[0].Type())
			}
		} else {
			return nil, ExceptionNewf(TypeError, "endswith() takes at least 1 argument (0 given)")
		}
		for _, s := range suffix {
			if strings.HasSuffix(selfStr, s) {
				return Bool(true), nil
			}
		}
		return Bool(false), nil
	}, 0, "endswith(suffix[, start[, end]]) -> bool")

	StringType.Dict["count"] = MustNewMethod("count", func(self Object, args Tuple) (Object, error) {
		return self.(String).Count(args)
	}, 0, `count(sub[, start[, end]]) -> int
Return the number of non-overlapping occurrences of substring sub in
string S[start:end].  Optional arguments start and end are
interpreted as in slice notation.`)

	StringType.Dict["find"] = MustNewMethod("find", func(self Object, args Tuple) (Object, error) {
		return self.(String).find(args)
	}, 0, `find(...)
S.find(sub[, start[, end]]) -> int

Return the lowest index in S where substring sub is found,
such that sub is contained within S[start:end].  Optional
arguments start and end are interpreted as in slice notation.

Return -1 on failure.`)

	StringType.Dict["replace"] = MustNewMethod("replace", func(self Object, args Tuple) (Object, error) {
		return self.(String).Replace(args)
	}, 0, `replace(self, old, new, count=-1) -> return a copy with all occurrences of substring old replaced by new.

  count
    Maximum number of occurrences to replace.
    -1 (the default value) means replace all occurrences.

If the optional argument count is given, only the first count occurrences are
replaced.`)

	StringType.Dict["split"] = MustNewMethod("split", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
		return self.(String).Split(args, kwargs)
	}, 0, "split(sub) -> split string with sub.")

	StringType.Dict["startswith"] = MustNewMethod("startswith", func(self Object, args Tuple) (Object, error) {
		selfStr := string(self.(String))
		prefix := []string{}
		if len(args) > 0 {
			if s, ok := args[0].(String); ok {
				prefix = append(prefix, string(s))
			} else if s, ok := args[0].(Tuple); ok {
				for _, t := range s {
					if v, ok := t.(String); ok {
						prefix = append(prefix, string(v))
					}
				}
			} else {
				return nil, ExceptionNewf(TypeError, "startswith first arg must be str, unicode, or tuple, not %s", args[0].Type())
			}
		} else {
			return nil, ExceptionNewf(TypeError, "startswith() takes at least 1 argument (0 given)")
		}
		if len(args) > 1 {
			if s, ok := args[1].(Int); ok {
				selfStr = selfStr[s:]
			}
		}

		for _, s := range prefix {
			if strings.HasPrefix(selfStr, s) {
				return Bool(true), nil
			}
		}
		return Bool(false), nil
	}, 0, "startswith(prefix[, start[, end]]) -> bool")

	StringType.Dict["strip"] = MustNewMethod("strip", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
		return self.(String).Strip(args)
	}, 0, "strip(chars) -> replace chars from begining and end of string")

	StringType.Dict["rstrip"] = MustNewMethod("rstrip", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
		return self.(String).RStrip(args)
	}, 0, "rstrip(chars) -> replace chars from end of string")

	StringType.Dict["lstrip"] = MustNewMethod("lstrip", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
		return self.(String).LStrip(args)
	}, 0, "lstrip(chars) -> replace chars from begining of string")

	StringType.Dict["upper"] = MustNewMethod("upper", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
		return self.(String).Upper()
	}, 0, "upper() -> a copy of the string converted to uppercase")

	StringType.Dict["lower"] = MustNewMethod("lower", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
		return self.(String).Lower()
	}, 0, "lower() -> a copy of the string converted to lowercase")

	StringType.Dict["join"] = MustNewMethod("join", func(self Object, args Tuple) (Object, error) {
		return self.(String).Join(args)
	}, 0, "join(iterable) -> return a string which is the concatenation of the strings in iterable")
}

// Type of this object
func (s String) Type() *Type {
	return StringType
}

// StrNew
func StrNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	var (
		sObj     Object = String("")
		encoding Object
		errors   Object
	)
	// FIXME ignoring encoding and errors
	err := ParseTupleAndKeywords(args, kwargs, "|OOO:str", []string{"bytes_or_buffer", "encoding", "errors"}, &sObj, &encoding, &errors)
	if err != nil {
		return nil, err
	}
	// FIXME ignoring encoding
	// FIXME ignoring buffer protocol
	return Str(sObj)
}

// Intern s possibly returning a reference to an already interned string
func (s String) Intern() String {
	// fmt.Printf("FIXME interning %q\n", s)
	return s
}

func (a String) M__str__() (Object, error) {
	return a, nil
}

func (a String) M__repr__() (Object, error) {
	out := StringEscape(a, false)
	return String(out), nil
}

func (s String) M__bool__() (Object, error) {
	return NewBool(len(s) > 0), nil
}

// len returns length of the string in unicode characters
func (s String) len() int {
	return utf8.RuneCountInString(string(s))
}

func (s String) M__len__() (Object, error) {
	return Int(s.len()), nil
}

func (a String) M__add__(other Object) (Object, error) {
	if b, ok := other.(String); ok {
		return a + b, nil
	}
	return NotImplemented, nil
}

func (a String) M__radd__(other Object) (Object, error) {
	if b, ok := other.(String); ok {
		return b + a, nil
	}
	return NotImplemented, nil
}

func (a String) M__iadd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a String) M__mul__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			b = 0
		}
		var out bytes.Buffer
		for i := 0; i < int(b); i++ {
			out.WriteString(string(a))
		}
		return String(out.String()), nil
	}
	return NotImplemented, nil
}

func (a String) M__rmul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a String) M__imul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

// Convert an Object to an String
//
// Returns ok as to whether the conversion worked or not
func convertToString(other Object) (String, bool) {
	switch b := other.(type) {
	case String:
		return b, true
	}
	return "", false
}

// Rich comparison

func (a String) M__lt__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(a < b), nil
	}
	return NotImplemented, nil
}

func (a String) M__le__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(a <= b), nil
	}
	return NotImplemented, nil
}

func (a String) M__eq__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(a == b), nil
	}
	return NotImplemented, nil
}

func (a String) M__ne__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(a != b), nil
	}
	return NotImplemented, nil
}

func (a String) M__gt__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(a > b), nil
	}
	return NotImplemented, nil
}

func (a String) M__ge__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(a >= b), nil
	}
	return NotImplemented, nil
}

// % operator

/*
4.7.2. printf-style String Formatting

Note The formatting operations described here exhibit a variety of
quirks that lead to a number of common errors (such as failing to
display tuples and dictionaries correctly). Using the newer
str.format() interface helps avoid these errors, and also provides a
generally more powerful, flexible and extensible approach to
formatting text.

String objects have one unique built-in operation: the % operator
(modulo). This is also known as the string formatting or interpolation
operator. Given format % values (where format is a string), %
conversion specifications in format are replaced with zero or more
elements of values. The effect is similar to using the sprintf() in
the C language.

If format requires a single argument, values may be a single non-tuple
object. [5] Otherwise, values must be a tuple with exactly the number
of items specified by the format string, or a single mapping object
(for example, a dictionary).

A conversion specifier contains two or more characters and has the
following components, which must occur in this order:

The '%' character, which marks the start of the specifier.

Mapping key (optional), consisting of a parenthesised sequence of
characters (for example, (somename)).

Conversion flags (optional), which affect the result of some
conversion types.

Minimum field width (optional). If specified as an '*' (asterisk), the
actual width is read from the next element of the tuple in values, and
the object to convert comes after the minimum field width and optional
precision.

Precision (optional), given as a '.' (dot) followed by the
precision. If specified as '*' (an asterisk), the actual precision is
read from the next element of the tuple in values, and the value to
convert comes after the precision.

Length modifier (optional).

Conversion type.

When the right argument is a dictionary (or other mapping type), then
the formats in the string must include a parenthesised mapping key
into that dictionary inserted immediately after the '%' character. The
mapping key selects the value to be formatted from the mapping. For
example:

>>>
>>> print('%(language)s has %(number)03d quote types.' %
...       {'language': "Python", "number": 2})
Python has 002 quote types.

In this case no * specifiers may occur in a format (since they require
a sequential parameter list).

The conversion flag characters are:

Flag	Meaning
'#'	The value conversion will use the “alternate form” (where defined below).
'0'	The conversion will be zero padded for numeric values.
'-'	The converted value is left adjusted (overrides the '0' conversion if both are given).
' '	(a space) A blank should be left before a positive number (or empty string) produced by a signed conversion.
'+'	A sign character ('+' or '-') will precede the conversion (overrides a “space” flag).

A length modifier (h, l, or L) may be present, but is ignored as it is
not necessary for Python – so e.g. %ld is identical to %d.

The conversion types are:

Conversion	Meaning	Notes
'd'	Signed integer decimal.
'i'	Signed integer decimal.
'o'	Signed octal value.	(1)
'u'	Obsolete type – it is identical to 'd'.	(7)
'x'	Signed hexadecimal (lowercase).	(2)
'X'	Signed hexadecimal (uppercase).	(2)
'e'	Floating point exponential format (lowercase).	(3)
'E'	Floating point exponential format (uppercase).	(3)
'f'	Floating point decimal format.	(3)
'F'	Floating point decimal format.	(3)
'g'	Floating point format. Uses lowercase exponential format if exponent is less than -4 or not less than precision, decimal format otherwise.	(4)
'G'	Floating point format. Uses uppercase exponential format if exponent is less than -4 or not less than precision, decimal format otherwise.	(4)
'c'	Single character (accepts integer or single character string).
'r'	String (converts any Python object using repr()).	(5)
's'	String (converts any Python object using str()).	(5)
'a'	String (converts any Python object using ascii()).	(5)
'%'	No argument is converted, results in a '%' character in the result.
Notes:

The alternate form causes a leading zero ('0') to be inserted between
left-hand padding and the formatting of the number if the leading
character of the result is not already a zero.

The alternate form causes a leading '0x' or '0X' (depending on whether
the 'x' or 'X' format was used) to be inserted between left-hand
padding and the formatting of the number if the leading character of
the result is not already a zero.

The alternate form causes the result to always contain a decimal
point, even if no digits follow it.

The precision determines the number of digits after the decimal point
and defaults to 6.

The alternate form causes the result to always contain a decimal
point, and trailing zeroes are not removed as they would otherwise be.

The precision determines the number of significant digits before and
after the decimal point and defaults to 6.

If precision is N, the output is truncated to N characters.

See PEP 237.  Since Python strings have an explicit length, %s
conversions do not assume that '\0' is the end of the string.

Changed in version 3.1: %f conversions for numbers whose absolute
value is over 1e50 are no longer replaced by %g conversions.
*/
func (a String) M__mod__(other Object) (Object, error) {
	var values Tuple
	switch b := other.(type) {
	case Tuple:
		values = b
	default:
		values = Tuple{other}
	}
	// FIXME not a full implementation ;-)
	params := make([]interface{}, len(values))
	for i := range values {
		params[i] = values[i]
	}
	s := string(a)
	s = strings.Replace(s, "%s", "%v", -1)
	s = strings.Replace(s, "%r", "%#v", -1)
	return String(fmt.Sprintf(s, params...)), nil
}

func (a String) M__rmod__(other Object) (Object, error) {
	switch b := other.(type) {
	case String:
		return b.M__mod__(a)
	}
	return NotImplemented, nil
}

func (a String) M__imod__(other Object) (Object, error) {
	return a.M__mod__(other)
}

// Returns position in string of n-th character
//
// returns end of string if not found
func (s String) pos(n int) int {
	characterNumber := 0
	for i := range s {
		if characterNumber == n {
			return i
		}
		characterNumber++
	}
	return len(s)
}

// slice returns the slice of this string using character positions
//
// length should be the length of the string in unicode characters
func (s String) slice(start, stop, length int) String {
	if start >= stop {
		return String("")
	}
	if length == len(s) {
		return s[start:stop] // ascii only
	}
	if start <= 0 && stop >= length {
		return s
	}
	startI := s.pos(start)
	stopI := s[startI:].pos(stop-start) + startI
	return s[startI:stopI]
}

func (s String) M__getitem__(key Object) (Object, error) {
	length := s.len()
	asciiOnly := length == len(s)
	if slice, ok := key.(*Slice); ok {
		start, stop, step, slicelength, err := slice.GetIndices(length)
		if err != nil {
			return nil, err
		}
		if step == 1 {
			// Return a subslice since strings are immutable
			return s.slice(start, stop, length), nil
		}
		if asciiOnly {
			newString := make([]byte, slicelength)
			for i, j := start, 0; j < slicelength; i, j = i+step, j+1 {
				newString[j] = s[i]
			}
			return String(newString), nil
		}
		// Unpack the string into a []rune to do this for speed
		runeString := []rune(string(s))
		newString := make([]rune, slicelength)
		for i, j := start, 0; j < slicelength; i, j = i+step, j+1 {
			newString[j] = runeString[i]
		}
		return String(newString), nil
	}
	i, err := IndexIntCheck(key, length)
	if err != nil {
		return nil, err
	}
	if asciiOnly {
		return s[i : i+1], nil
	}
	s = s[s.pos(i):]
	_, runeSize := utf8.DecodeRuneInString(string(s))
	return s[:runeSize], nil
}

func (s String) M__contains__(item Object) (Object, error) {
	needle, ok := item.(String)
	if !ok {
		return nil, ExceptionNewf(TypeError, "'in <string>' requires string as left operand, not %s", item.Type().Name)
	}
	return NewBool(strings.Contains(string(s), string(needle))), nil
}

func (s String) Count(args Tuple) (Object, error) {
	var (
		pysub Object
		pybeg Object = Int(0)
		pyend Object = Int(s.len())
		pyfmt        = "s|ii:count"
	)
	err := ParseTuple(args, pyfmt, &pysub, &pybeg, &pyend)
	if err != nil {
		return nil, err
	}

	var (
		beg  = int(pybeg.(Int))
		end  = int(pyend.(Int))
		size = s.len()
	)
	if beg > size {
		beg = size
	}
	if end < 0 {
		end = size
	}
	if end > size {
		end = size
	}

	var (
		str = string(s.slice(beg, end, s.len()))
		sub = string(pysub.(String))
	)
	return Int(strings.Count(str, sub)), nil
}

func (s String) find(args Tuple) (Object, error) {
	var (
		pysub Object
		pybeg Object = Int(0)
		pyend Object = Int(s.len())
		pyfmt        = "s|ii:find"
	)
	err := ParseTuple(args, pyfmt, &pysub, &pybeg, &pyend)
	if err != nil {
		return nil, err
	}

	var (
		beg  = int(pybeg.(Int))
		end  = int(pyend.(Int))
		size = s.len()
	)
	if beg > size {
		beg = size
	}
	if end < 0 {
		end = size
	}
	if end > size {
		end = size
	}

	var (
		off = s.slice(0, beg, s.len()).len()
		str = string(s.slice(beg, end, s.len()))
		sub = string(pysub.(String))
		idx = strings.Index(str, sub)
	)
	if idx < 0 {
		return Int(idx), nil
	}
	return Int(off + String(str[:idx]).len()), nil
}

func (s String) Split(args Tuple, kwargs StringDict) (Object, error) {
	var (
		pyval Object = None
		pymax Object = Int(-2)
		pyfmt        = "|Oi:split"
		kwlst        = []string{"sep", "maxsplit"}
	)
	err := ParseTupleAndKeywords(args, kwargs, pyfmt, kwlst, &pyval, &pymax)
	if err != nil {
		return nil, err
	}

	var (
		max = pymax.(Int)
		vs  []string
	)
	switch v := pyval.(type) {
	case String:
		vs = strings.SplitN(string(s), string(v), int(max)+1)
	case NoneType:
		vs = fieldsN(string(s), int(max))
	default:
		return nil, ExceptionNewf(TypeError, "Can't convert '%s' object to str implicitly", pyval.Type())
	}
	o := List{}
	for _, j := range vs {
		o.Items = append(o.Items, String(j))
	}
	return &o, nil
}

func (s String) Replace(args Tuple) (Object, error) {
	var (
		pyold Object = None
		pynew Object = None
		pycnt Object = Int(-1)
	)
	err := ParseTuple(args, "ss|i:replace", &pyold, &pynew, &pycnt)
	if err != nil {
		return nil, err
	}

	var (
		old = string(pyold.(String))
		new = string(pynew.(String))
		cnt = int(pycnt.(Int))
	)

	return String(strings.Replace(string(s), old, new, cnt)), nil
}

func stripFunc(args Tuple) (func(rune) bool, error) {
	var (
		pyval Object = None
	)
	err := ParseTuple(args, "|s", &pyval)
	if err != nil {
		return nil, err
	}
	f := unicode.IsSpace
	switch v := pyval.(type) {
	case String:
		chars := []rune(string(v))
		f = func(s rune) bool {
			for _, i := range chars {
				if s == i {
					return true
				}
			}
			return false
		}
	}
	return f, nil
}

func (s String) Strip(args Tuple) (Object, error) {
	f, err := stripFunc(args)
	if err != nil {
		return nil, err
	}
	return String(strings.TrimFunc(string(s), f)), nil
}

func (s String) LStrip(args Tuple) (Object, error) {
	f, err := stripFunc(args)
	if err != nil {
		return nil, err
	}
	return String(strings.TrimLeftFunc(string(s), f)), nil
}

func (s String) RStrip(args Tuple) (Object, error) {
	f, err := stripFunc(args)
	if err != nil {
		return nil, err
	}
	return String(strings.TrimRightFunc(string(s), f)), nil
}

func (s String) Upper() (Object, error) {
	return String(strings.ToUpper(string(s))), nil
}

func (s String) Lower() (Object, error) {
	return String(strings.ToLower(string(s))), nil
}

func (s String) Join(args Tuple) (Object, error) {
	if len(args) != 1 {
		return nil, ExceptionNewf(TypeError, "join() takes exactly one argument (%d given)", len(args))
	}
	var parts []string
	iterable, err := Iter(args[0])
	if err != nil {
		return nil, err
	}
	item, err := Next(iterable)
	for err == nil {
		str, ok := item.(String)
		if !ok {
			return nil, ExceptionNewf(TypeError, "sequence item %d: expected str instance, %s found", len(parts), item.Type().Name)
		}
		parts = append(parts, string(str))
		item, err = Next(iterable)
	}
	if err != StopIteration {
		return nil, err
	}
	return String(strings.Join(parts, string(s))), nil
}

// Check stringerface is satisfied
var (
	_ richComparison     = String("")
	_ sequenceArithmetic = String("")
	_ I__mod__           = String("")
	_ I__rmod__          = String("")
	_ I__imod__          = String("")
	_ I__len__           = String("")
	_ I__bool__          = String("")
	_ I__getitem__       = String("")
	_ I__contains__      = String("")
)
