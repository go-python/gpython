// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Time module

package time

import (
	"time"

	"github.com/go-python/gpython/py"
)

const time_doc = `time() -> floating point number

Return the current time in seconds since the Epoch.
Fractions of a second may be present if the system clock provides them.`

func time_time(self py.Object) (py.Object, error) {
	return py.Float(time.Now().UnixNano()) / 1e9, nil
}

// func floatclock(_Py_clock_info_t *info) (py.Object, error) {
// 	value := clock()
// 	if value == (clock_t)-1 {
// 		PyErr_SetString(PyExc_RuntimeError, "the processor time used is not available or its value cannot be represented")
// 		return nil
// 	}
// 	if info {
// 		info.implementation = "clock()"
// 		info.resolution = 1.0 / float64(CLOCKS_PER_SEC)
// 		info.monotonic = 1
// 		info.adjustable = 0
// 	}
// 	return py.Float(float64(value) / CLOCKS_PER_SEC)
// }

// func pyclock(info *_Py_clock_info_t) (py.Object, error) {
// 	return floatclock(info)
// }

const clock_doc = `clock() -> floating point number

Return the CPU time or real time since the start of the process or since
the first call to clock().  This has as much precision as the system
records.`

func time_clock(self py.Object) (py.Object, error) {
	return time_time(self)
}

const clock_gettime_doc = `clock_gettime(clk_id) -> floating point number

Return the time of the specified clock clk_id.`

func time_clock_gettime(self py.Object, args py.Tuple) (py.Object, error) {
	// var ret int
	// var clk_id int
	// var tp timespec

	// if !PyArg_ParseTuple(args, "i:clock_gettime", &clk_id) {
	// 	return nil
	// }

	// ret = clock_gettime(clk_id, &tp)
	// if ret != 0 {
	// 	PyErr_SetFromErrno(PyExc_IOError)
	// 	return nil
	// }
	// return py.Float(tp.tv_sec + tp.tv_nsec*1e-9)
	return nil, py.NotImplementedError
}

const clock_settime_doc = `clock_settime(clk_id, time)

Set the time of the specified clock clk_id.`

func time_clock_settime(self py.Object, args py.Tuple) (py.Object, error) {
	// var clk_id int
	// var obj py.Object
	// var tv_sec time_t
	// var tv_nsec int64
	// var tp timespec
	// var ret int

	// if !PyArg_ParseTuple(args, "iO:clock_settime", &clk_id, &obj) {
	// 	return nil
	// }

	// if _PyTime_ObjectToTimespec(obj, &tv_sec, &tv_nsec) == -1 {
	// 	return nil
	// }
	// tp.tv_sec = tv_sec
	// tp.tv_nsec = tv_nsec

	// ret = clock_settime(clk_id, &tp)
	// if ret != 0 {
	// 	PyErr_SetFromErrno(PyExc_IOError)
	// 	return nil
	// }
	// Py_RETURN_NONE
	return nil, py.NotImplementedError
}

const clock_getres_doc = `clock_getres(clk_id) -> floating point number

Return the resolution (precision) of the specified clock clk_id.`

func time_clock_getres(self py.Object, args py.Tuple) (py.Object, error) {
	// var ret int
	// var clk_id int
	// var tp timespec

	// if !PyArg_ParseTuple(args, "i:clock_getres", &clk_id) {
	// 	return nil
	// }

	// ret = clock_getres(clk_id, &tp)
	// if ret != 0 {
	// 	PyErr_SetFromErrno(PyExc_IOError)
	// 	return nil
	// }

	// return py.Float(tp.tv_sec + tp.tv_nsec*1e-9)
	return nil, py.NotImplementedError
}

const sleep_doc = `sleep(seconds)

Delay execution for a given number of seconds.  The argument may be
a floating point number for subsecond precision.`

func time_sleep(self py.Object, args py.Tuple) (py.Object, error) {
	var secsObj py.Object
	err := py.ParseTuple(args, "d:sleep", &secsObj)
	if err != nil {
		return nil, err
	}
	secs := secsObj.(py.Float)
	if secs < 0 {
		return nil, py.ExceptionNewf(py.ValueError, "sleep length must be non-negative")
	}
	time.Sleep(time.Duration(secs * 1e9))
	return py.None, nil
}

// var struct_time_type_fields = []PyStructSequence_Field{
// 	{"tm_year", "year, for example, 1993"},
// 	{"tm_mon", "month of year, range [1, 12]"},
// 	{"tm_mday", "day of month, range [1, 31]"},
// 	{"tm_hour", "hours, range [0, 23]"},
// 	{"tm_min", "minutes, range [0, 59]"},
// 	{"tm_sec", "seconds, range [0, 61])"},
// 	{"tm_wday", "day of week, range [0, 6], Monday is 0"},
// 	{"tm_yday", "day of year, range [1, 366]"},
// 	{"tm_isdst", "1 if summer time is in effect, 0 if not, and -1 if unknown"},
// 	{"tm_zone", "abbreviation of timezone name"},
// 	{"tm_gmtoff", "offset from UTC in seconds"},
// }

// var struct_time_type_desc = PyStructSequence_Desc{
// 	"time.struct_time",
// 	`The time value as returned by gmtime(), localtime(), and strptime(), and
//  accepted by asctime(), mktime() and strftime().  May be considered as a
//  sequence of 9 integers.

//  Note that several fields' values are not the same as those defined by
//  the C language standard for struct tm.  For example, the value of the
//  field tm_year is the actual year, not year - 1900.  See individual
//  fields' descriptions for details.`,
// 	struct_time_type_fields,
// 	9,
// }

// var initialized int
// var StructTimeType PyTypeObject

// func tmtotuple(p *tm) py.Object {
// 	v := PyStructSequence_New(&StructTimeType)
// 	if v == nil {
// 		return nil
// 	}

// 	// #define SET(i,val) PyStructSequence_SET_ITEM(v, i, PyLong_FromLong((long) val))

// 	SET(0, p.tm_year+1900)
// 	SET(1, p.tm_mon+1) /* Want January == 1 */
// 	SET(2, p.tm_mday)
// 	SET(3, p.tm_hour)
// 	SET(4, p.tm_min)
// 	SET(5, p.tm_sec)
// 	SET(6, (p.tm_wday+6)%7) /* Want Monday == 0 */
// 	SET(7, p.tm_yday+1)     /* Want January, 1 == 1 */
// 	SET(8, p.tm_isdst)
// 	PyStructSequence_SET_ITEM(v, 9,
// 		PyUnicode_DecodeLocale(p.tm_zone, "surrogateescape"))
// 	SET(10, p.tm_gmtoff)
// 	if PyErr_Occurred() {
// 		Py_XDECREF(v)
// 		return nil
// 	}

// 	return v
// }

// /* Parse arg tuple that can contain an optional float-or-None value;
//    format needs to be "|O:name".
//    Returns non-zero on success (parallels PyArg_ParseTuple).
// */
// func parse_time_t_args(args py.Tuple, format string, pwhen *time_t) int {
// 	var ot py.Object
// 	var whent time_t

// 	if !PyArg_ParseTuple(args, format, &ot) {
// 		return 0
// 	}
// 	if ot == nil || ot == Py_None {
// 		whent = time(nil)
// 	} else {
// 		if _PyTime_ObjectToTime_t(ot, &whent) == -1 {
// 			return 0
// 		}
// 	}
// 	*pwhen = whent
// 	return 1
// }

const gmtime_doc = `gmtime([seconds]) -> (tm_year, tm_mon, tm_mday, tm_hour, tm_min,
                       tm_sec, tm_wday, tm_yday, tm_isdst)

Convert seconds since the Epoch to a time tuple expressing UTC (a.k.a.
GMT).  When 'seconds' is not passed in, convert the current time instead.

If the platform supports the tm_gmtoff and tm_zone, they are available as
attributes only.`

func time_gmtime(self py.Object, args py.Tuple) (py.Object, error) {
	// var when time_t
	// var buf tm
	// var local *tm

	// if !parse_time_t_args(args, "|O:gmtime", &when) {
	// 	return nil
	// }

	// errno = 0
	// local = gmtime(&when)
	// if local == nil {
	// 	if errno == 0 {
	// 		errno = EINVAL
	// 	}
	// 	return PyErr_SetFromErrno(PyExc_OSError)
	// }
	// buf = *local
	// return tmtotuple(&buf)
	return nil, py.NotImplementedError
}

// func pylocaltime(timep *time_t, result *tm) int {
// 	var local *tm

// 	assert(timep != nil)
// 	local = localtime(timep)
// 	if local == nil {
// 		/* unconvertible time */
// 		if errno == 0 {
// 			errno = EINVAL
// 		}
// 		PyErr_SetFromErrno(PyExc_OSError)
// 		return -1
// 	}
// 	*result = *local
// 	return 0
// }

const localtime_doc = `localtime([seconds]) -> (tm_year,tm_mon,tm_mday,tm_hour,tm_min,
                          tm_sec,tm_wday,tm_yday,tm_isdst)

Convert seconds since the Epoch to a time tuple expressing local time.
When 'seconds' is not passed in, convert the current time instead.`

func time_localtime(self py.Object, args py.Tuple) (py.Object, error) {
	// var when time_t
	// var buf tm

	// if !parse_time_t_args(args, "|O:localtime", &when) {
	// 	return nil
	// }
	// if pylocaltime(&when, &buf) == -1 {
	// 	return nil
	// }
	// return tmtotuple(&buf)
	return nil, py.NotImplementedError
}

/* Convert 9-item tuple to tm structure.  Return 1 on success, set
 * an exception and return 0 on error.
 */

// func gettmarg(args py.Tuple, p *tm) int {
// 	var y int

// 	// FIXME memset(p, '\0', sizeof(struct tm));

// 	if !PyTuple_Check(args) {
// 		PyErr_SetString(PyExc_TypeError,
// 			"Tuple or struct_time argument required")
// 		return 0
// 	}

// 	if !PyArg_ParseTuple(args, "iiiiiiiii",
// 		&y, &p.tm_mon, &p.tm_mday,
// 		&p.tm_hour, &p.tm_min, &p.tm_sec,
// 		&p.tm_wday, &p.tm_yday, &p.tm_isdst) {
// 		return 0
// 	}
// 	p.tm_year = y - 1900
// 	p.tm_mon--
// 	p.tm_wday = (p.tm_wday + 1) % 7
// 	p.tm_yday--
// 	// if (Py_TYPE(args) == &StructTimeType) {
// 	// 	    item := PyTuple_GET_ITEM(args, 9);
// 	//     p.tm_zone = item == Py_None ? nil : _PyUnicode_AsString(item);
// 	//     item = PyTuple_GET_ITEM(args, 10);
// 	//     p.tm_gmtoff = item == Py_None ? 0 : PyLong_AsLong(item);
// 	//     if (PyErr_Occurred()) {
// 	//         return 0;
// 	//     }
// 	// }
// 	return 1
// }

/* Check values of the struct tm fields before it is passed to strftime() and
 * asctime().  Return 1 if all values are valid, otherwise set an exception
 * and returns 0.
 */

// func checktm(buf *tm) int {
// 	/* Checks added to make sure strftime() and asctime() does not crash Python by
// 	   indexing blindly into some array for a textual representation
// 	   by some bad index (fixes bug #897625 and #6608).

// 	   Also support values of zero from Python code for arguments in which
// 	   that is out of range by forcing that value to the lowest value that
// 	   is valid (fixed bug #1520914).

// 	   Valid ranges based on what is allowed in struct tm:

// 	   - tm_year: [0, max(int)] (1)
// 	   - tm_mon: [0, 11] (2)
// 	   - tm_mday: [1, 31]
// 	   - tm_hour: [0, 23]
// 	   - tm_min: [0, 59]
// 	   - tm_sec: [0, 60]
// 	   - tm_wday: [0, 6] (1)
// 	   - tm_yday: [0, 365] (2)
// 	   - tm_isdst: [-max(int), max(int)]

// 	   (1) gettmarg() handles bounds-checking.
// 	   (2) Python's acceptable range is one greater than the range in C,
// 	   thus need to check against automatic decrement by gettmarg().
// 	*/
// 	if buf.tm_mon == -1 {
// 		buf.tm_mon = 0
// 	} else if buf.tm_mon < 0 || buf.tm_mon > 11 {
// 		PyErr_SetString(PyExc_ValueError, "month out of range")
// 		return 0
// 	}
// 	if buf.tm_mday == 0 {
// 		buf.tm_mday = 1
// 	} else if buf.tm_mday < 0 || buf.tm_mday > 31 {
// 		PyErr_SetString(PyExc_ValueError, "day of month out of range")
// 		return 0
// 	}
// 	if buf.tm_hour < 0 || buf.tm_hour > 23 {
// 		PyErr_SetString(PyExc_ValueError, "hour out of range")
// 		return 0
// 	}
// 	if buf.tm_min < 0 || buf.tm_min > 59 {
// 		PyErr_SetString(PyExc_ValueError, "minute out of range")
// 		return 0
// 	}
// 	if buf.tm_sec < 0 || buf.tm_sec > 61 {
// 		PyErr_SetString(PyExc_ValueError, "seconds out of range")
// 		return 0
// 	}
// 	/* tm_wday does not need checking of its upper-bound since taking
// 	   ``% 7`` in gettmarg() automatically restricts the range. */
// 	if buf.tm_wday < 0 {
// 		PyErr_SetString(PyExc_ValueError, "day of week out of range")
// 		return 0
// 	}
// 	if buf.tm_yday == -1 {
// 		buf.tm_yday = 0
// 	} else if buf.tm_yday < 0 || buf.tm_yday > 365 {
// 		PyErr_SetString(PyExc_ValueError, "day of year out of range")
// 		return 0
// 	}
// 	return 1
// }

const strftime_doc = `strftime(format[, tuple]) -> string

Convert a time tuple to a string according to a format specification.
See the library reference manual for formatting codes. When the time tuple
is not present, current time as returned by localtime() is used.`

func time_strftime(self py.Object, args py.Tuple) (py.Object, error) {
	// var tup py.Object
	// var buf tm
	// var fmt *time_char
	// var format py.Object
	// var format_arg py.Object
	// var fmtlen, buflen int
	// var outbuf *time_char
	// var i int
	// var ret py.Object

	// // memset((void *) &buf, '\0', sizeof(buf));

	// /* Will always expect a unicode string to be passed as format.
	//    Given that there's no str type anymore in py3k this seems safe.
	// */
	// if !PyArg_ParseTuple(args, "U|O:strftime", &format_arg, &tup) {
	// 	return nil
	// }

	// if tup == nil {
	// 	tt := time(nil)
	// 	if pylocaltime(&tt, &buf) == -1 {
	// 		return nil
	// 	}
	// } else if !gettmarg(tup, &buf) || !checktm(&buf) {
	// 	return nil
	// }

	// /* Normalize tm_isdst just in case someone foolishly implements %Z
	//    based on the assumption that tm_isdst falls within the range of
	//    [-1, 1] */
	// if buf.tm_isdst < -1 {
	// 	buf.tm_isdst = -1
	// } else if buf.tm_isdst > 1 {
	// 	buf.tm_isdst = 1
	// }

	// /* Convert the unicode string to an ascii one */
	// format = PyUnicode_EncodeLocale(format_arg, "surrogateescape")
	// if format == nil {
	// 	return nil
	// }
	// fmt = PyBytes_AS_STRING(format)

	// fmtlen = time_strlen(fmt)

	// /* I hate these functions that presume you know how big the output
	//  * will be ahead of time...
	//  */
	// for i = 1024; ; i += i {
	// 	outbuf = PyMem_Malloc(i * sizeof(time_char))
	// 	if outbuf == nil {
	// 		PyErr_NoMemory()
	// 		break
	// 	}
	// 	buflen = format_time(outbuf, i, fmt, &buf)
	// 	if buflen > 0 || i >= 256*fmtlen {
	// 		/* If the buffer is 256 times as long as the format,
	// 		   it's probably not failing for lack of room!
	// 		   More likely, the format yields an empty result,
	// 		   e.g. an empty format, or %Z when the timezone
	// 		   is unknown. */
	// 		ret = PyUnicode_DecodeLocaleAndSize(outbuf, buflen,
	// 			"surrogateescape")
	// 		PyMem_Free(outbuf)
	// 		break
	// 	}
	// 	PyMem_Free(outbuf)
	// }
	// Py_DECREF(format)
	// return ret
	return nil, py.NotImplementedError
}

const strptime_doc = `strptime(string, format) -> struct_time

Parse a string to a time tuple according to a format specification.
See the library reference manual for formatting codes (same as strftime()).`

func time_strptime(self py.Object, args py.Tuple) (py.Object, error) {
	// strptime_module := PyImport_ImportModuleNoBlock("_strptime")

	// if !strptime_module {
	// 	return nil
	// }
	// strptime_result = _py.Object_CallMethodId(strptime_module,
	// 	"_strptime_time", "O", args)
	// Py_DECREF(strptime_module)
	// return strptime_result
	return nil, py.NotImplementedError
}

// func _asctime(timeptr *tm) py.Object {
// 	/* Inspired by Open Group reference implementation available at
// 	 * http://pubs.opengroup.org/onlinepubs/009695399/functions/asctime.html */
// 	var wday_name = [7]string{
// 		"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat",
// 	}
// 	var mon_name = [12]string{
// 		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
// 		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
// 	}
// 	return PyUnicode_FromFormat(
// 		"%s %s%3d %.2d:%.2d:%.2d %d",
// 		wday_name[timeptr.tm_wday],
// 		mon_name[timeptr.tm_mon],
// 		timeptr.tm_mday, timeptr.tm_hour,
// 		timeptr.tm_min, timeptr.tm_sec,
// 		1900+timeptr.tm_year)
// }

const asctime_doc = `asctime([tuple]) -> string

Convert a time tuple to a string, e.g. 'Sat Jun 06 16:26:11 1998'.
When the time tuple is not present, current time as returned by localtime()
is used.`

func time_asctime(self py.Object, args py.Tuple) (py.Object, error) {
	// var tup py.Object
	// var buf tm

	// if !PyArg_UnpackTuple(args, "asctime", 0, 1, &tup) {
	// 	return nil
	// }
	// if tup == nil {
	// 	tt := time(nil)
	// 	if pylocaltime(&tt, &buf) == -1 {
	// 		return nil
	// 	}

	// } else if !gettmarg(tup, &buf) || !checktm(&buf) {
	// 	return nil
	// }
	// return _asctime(&buf)
	return nil, py.NotImplementedError
}

const ctime_doc = `ctime(seconds) -> string

Convert a time in seconds since the Epoch to a string in local time.
This is equivalent to asctime(localtime(seconds)). When the time tuple is
not present, current time as returned by localtime() is used.`

func time_ctime(self py.Object, args py.Tuple) (py.Object, error) {
	// var tt int
	// var tm buf
	// if !parse_time_t_args(args, "|O:ctime", &tt) {
	// 	return nil
	// }
	// if pylocaltime(&tt, &buf) == -1 {
	// 	return nil
	// }
	// return _asctime(&buf)
	return nil, py.NotImplementedError
}

const mktime_doc = `mktime(tuple) -> floating point number

Convert a time tuple in local time to seconds since the Epoch.
Note that mktime(gmtime(0)) will not generally return zero for most
time zones; instead the returned value will either be equal to that
of the timezone or altzone attributes on the time module.`

func time_mktime(self, tup py.Object) (py.Object, error) {
	// var buf tm
	// var tt int
	// if !gettmarg(tup, &buf) {
	// 	return nil
	// }
	// buf.tm_wday = -1 /* sentinel; original value ignored */
	// tt = mktime(&buf)
	// /* Return value of -1 does not necessarily mean an error, but tm_wday
	//  * cannot remain set to -1 if mktime succeeded. */
	// if tt == (time_t)(-1) && buf.tm_wday == -1 {
	// 	PyErr_SetString(PyExc_OverflowError,
	// 		"mktime argument out of range")
	// 	return nil
	// }
	// return py.Float(float64(tt))
	return nil, py.NotImplementedError
}

const tzset_doc = `tzset()

Initialize, or reinitialize, the local timezone to the value stored in
os.environ['TZ']. The TZ environment variable should be specified in
standard Unix timezone format as documented in the tzset man page
(eg. 'US/Eastern', 'Europe/Amsterdam'). Unknown timezones will silently
fall back to UTC. If the TZ environment variable is not set, the local
timezone is set to the systems best guess of wallclock time.
Changing the TZ environment variable without calling tzset *may* change
the local timezone used by methods such as localtime, but this behaviour
should not be relied on.`

func time_tzset(self py.Object) (py.Object, error) {
	// py.Object * m

	// m = PyImport_ImportModuleNoBlock("time")
	// if m == nil {
	// 	return nil
	// }

	// tzset()

	// /* Reset timezone, altzone, daylight and tzname */
	// PyInit_timezone(m)
	// Py_DECREF(m)

	// Py_INCREF(Py_None)
	// return Py_None
	return nil, py.NotImplementedError
}

const monotonic_doc = `monotonic() -> float

Monotonic clock, cannot go backward.`

// func pymonotonic(_Py_clock_info_t *info) (py.Object, error) {
// 	const clockid_t clk_id = CLOCK_HIGHRES
// 	const char *function = "clock_gettime(CLOCK_HIGHRES)"
// 	// const clockid_t clk_id = CLOCK_MONOTONIC;
// 	// const char *function = "clock_gettime(CLOCK_MONOTONIC)";

// 	if clock_gettime(clk_id, &tp) != 0 {
// 		PyErr_SetFromErrno(PyExc_OSError)
// 		return nil
// 	}

// 	if info {
// 		var res timespec
// 		info.monotonic = 1
// 		info.implementation = function
// 		info.adjustable = 0
// 		if clock_getres(clk_id, &res) == 0 {
// 			info.resolution = res.tv_sec + res.tv_nsec*1e-9
// 		} else {
// 			info.resolution = 1e-9
// 		}
// 	}
// 	return py.Float(tp.tv_sec + tp.tv_nsec*1e-9)
// }

func time_monotonic(self py.Object) (py.Object, error) {
	// return pymonotonic(nil)
	return nil, py.NotImplementedError
}

// func perf_counter(_Py_clock_info_t *info) py.Object {
// 	use_monotonic := true

// 	if use_monotonic {
// 		res = pymonotonic(info)
// 		if res != nil {
// 			return res
// 		}
// 		use_monotonic = false
// 		PyErr_Clear()
// 	}

// 	return floattime(info)
// }

const perf_counter_doc = `perf_counter() -> float

Performance counter for benchmarking.`

func time_perf_counter(self py.Object) (py.Object, error) {
	// return perf_counter(nil)
	return nil, py.NotImplementedError
}

// func py_process_time(_Py_clock_info_t *info) py.Object {

// #if defined(HAVE_SYS_RESOURCE_H)
//     struct rusage ru;
// #endif
// #ifdef HAVE_TIMES
//     struct tms t;
//      long ticks_per_second = -1;
// #endif

// #if defined(HAVE_CLOCK_GETTIME) \
//     && (defined(CLOCK_PROCESS_CPUTIME_ID) || defined(CLOCK_PROF))
//     struct timespec tp;
// #ifdef CLOCK_PROF
//     const clockid_t clk_id = CLOCK_PROF;
//     const char *function = "clock_gettime(CLOCK_PROF)";
// #else
//     const clockid_t clk_id = CLOCK_PROCESS_CPUTIME_ID;
//     const char *function = "clock_gettime(CLOCK_PROCESS_CPUTIME_ID)";
// #endif

//     if (clock_gettime(clk_id, &tp) == 0) {
//         if (info) {
//             struct timespec res;
//             info.implementation = function;
//             info.monotonic = 1;
//             info.adjustable = 0;
//             if (clock_getres(clk_id, &res) == 0) {
//                 info.resolution = res.tv_sec + res.tv_nsec * 1e-9;
//             } else {
//                 info.resolution = 1e-9;
//             }
//         }
//         return py.Float(tp.tv_sec + tp.tv_nsec * 1e-9);
//     }
// #endif

// #if defined(HAVE_SYS_RESOURCE_H)
//     if (getrusage(RUSAGE_SELF, &ru) == 0) {
//         float64 total;
//         total = ru.ru_utime.tv_sec + ru.ru_utime.tv_usec * 1e-6;
//         total += ru.ru_stime.tv_sec + ru.ru_stime.tv_usec * 1e-6;
//         if (info) {
//             info.implementation = "getrusage(RUSAGE_SELF)";
//             info.monotonic = 1;
//             info.adjustable = 0;
//             info.resolution = 1e-6;
//         }
//         return py.Float(total);
//     }
// #endif

// #ifdef HAVE_TIMES
//     if (times(&t) != (clock_t)-1) {
//         float64 total;

//         if (ticks_per_second == -1) {
// #if defined(HAVE_SYSCONF) && defined(_SC_CLK_TCK)
//             ticks_per_second = sysconf(_SC_CLK_TCK);
//             if (ticks_per_second < 1) {
//                 ticks_per_second = -1;
//             }
// #elif defined(HZ)
//             ticks_per_second = HZ;
// #else
//             ticks_per_second = 60; /* magic fallback value; may be bogus */
// #endif
//         }

//         if (ticks_per_second != -1) {
//             total = (float64)t.tms_utime / ticks_per_second;
//             total += (float64)t.tms_stime / ticks_per_second;
//             if (info) {
//                 info.implementation = "times()";
//                 info.monotonic = 1;
//                 info.adjustable = 0;
//                 info.resolution = 1.0 / ticks_per_second;
//             }
//             return py.Float(total);
//         }
//     }
// #endif

//     return floatclock(info);
// }

const process_time_doc = `process_time() . float

Process time for profiling: sum of the kernel and user-space CPU time.`

func time_process_time(self py.Object) (py.Object, error) {
	// return py_process_time(nil)
	return nil, py.NotImplementedError
}

const get_clock_info_doc = `get_clock_info(name: str) -> dict

Get information of the specified clock.`

func time_get_clock_info(self py.Object, args py.Tuple) (py.Object, error) {
	// 	char * name
	// 	var info _Py_clock_info_t
	// 	var obj, dict, ns py.Object

	// 	if !PyArg_ParseTuple(args, "s:get_clock_info", &name) {
	// 		return nil
	// 	}

	// 	info.implementation = ""
	// 	info.monotonic = 0
	// 	info.adjustable = 0
	// 	info.resolution = 1.0

	// 	if strcmp(name, "time") == 0 {
	// 		obj = floattime(&info)
	// 	} else if strcmp(name, "clock") == 0 {
	// 		obj = pyclock(&info)
	// 	} else if strcmp(name, "monotonic") == 0 {
	// 		obj = pymonotonic(&info)
	// 	} else if strcmp(name, "perf_counter") == 0 {
	// 		obj = perf_counter(&info)
	// 	} else if strcmp(name, "process_time") == 0 {
	// 		obj = py_process_time(&info)
	// 	} else {
	// 		PyErr_SetString(PyExc_ValueError, "unknown clock")
	// 		return nil
	// 	}
	// 	if obj == nil {
	// 		return nil
	// 	}
	// 	Py_DECREF(obj)

	// 	dict = PyDict_New()
	// 	if dict == nil {
	// 		return nil
	// 	}

	// 	assert(info.implementation != nil)
	// 	obj = PyUnicode_FromString(info.implementation)
	// 	if obj == nil {
	// 		goto error
	// 	}
	// 	if PyDict_SetItemString(dict, "implementation", obj) == -1 {
	// 		goto error
	// 	}
	// 	Py_CLEAR(obj)

	// 	assert(info.monotonic != -1)
	// 	obj = PyBool_FromLong(info.monotonic)
	// 	if obj == nil {
	// 		goto error
	// 	}
	// 	if PyDict_SetItemString(dict, "monotonic", obj) == -1 {
	// 		goto error
	// 	}
	// 	Py_CLEAR(obj)

	// 	assert(info.adjustable != -1)
	// 	obj = PyBool_FromLong(info.adjustable)
	// 	if obj == nil {
	// 		goto error
	// 	}
	// 	if PyDict_SetItemString(dict, "adjustable", obj) == -1 {
	// 		goto error
	// 	}
	// 	Py_CLEAR(obj)

	// 	assert(info.resolution > 0.0)
	// 	assert(info.resolution <= 1.0)
	// 	obj = py.Float(info.resolution)
	// 	if obj == nil {
	// 		goto error
	// 	}
	// 	if PyDict_SetItemString(dict, "resolution", obj) == -1 {
	// 		goto error
	// 	}
	// 	Py_CLEAR(obj)

	// 	ns = _PyNamespace_New(dict)
	// 	Py_DECREF(dict)
	// 	return ns

	// error:
	// 	Py_DECREF(dict)
	// 	Py_XDECREF(obj)
	// 	return nil
	return nil, py.NotImplementedError
}

func PyInit_timezone(m py.Object) {
	/* This code moved from PyInit_time wholesale to allow calling it from
	   time_tzset. In the future, some parts of it can be moved back
	   (for platforms that don't HAVE_WORKING_TZSET, when we know what they
	   are), and the extraneous calls to tzset(3) should be removed.
	   I haven't done this yet, as I don't want to change this code as
	   little as possible when introducing the time.tzset and time.tzsetwall
	   methods. This should simply be a method of doing the following once,
	   at the top of this function and removing the call to tzset() from
	   time_tzset():

	       #ifdef HAVE_TZSET
	       tzset()
	       #endif

	   And I'm lazy and hate C so nyer.
	*/
	// #if defined(HAVE_TZNAME) && !defined(__GLIBC__) && !defined(__CYGWIN__)
	//     py.Object otz0, *otz1;
	//     tzset();
	// #ifdef PYOS_OS2
	//     PyModule_AddIntConstant(m, "timezone", _timezone);
	// #else /* !PYOS_OS2 */
	//     PyModule_AddIntConstant(m, "timezone", timezone);
	// #endif /* PYOS_OS2 */
	// #ifdef HAVE_ALTZONE
	//     PyModule_AddIntConstant(m, "altzone", altzone);
	// #else
	// #ifdef PYOS_OS2
	//     PyModule_AddIntConstant(m, "altzone", _timezone-3600);
	// #else /* !PYOS_OS2 */
	//     PyModule_AddIntConstant(m, "altzone", timezone-3600);
	// #endif /* PYOS_OS2 */
	// #endif
	//     PyModule_AddIntConstant(m, "daylight", daylight);
	//     otz0 = PyUnicode_DecodeLocale(tzname[0], "surrogateescape");
	//     otz1 = PyUnicode_DecodeLocale(tzname[1], "surrogateescape");
	//     PyModule_AddObject(m, "tzname", Py_BuildValue("(NN)", otz0, otz1));
	// #else /* !HAVE_TZNAME || __GLIBC__ || __CYGWIN__*/
	// #ifdef HAVE_STRUCT_TM_TM_ZONE
	//     {
	// #define YEAR ((time_t)((365 * 24 + 6) * 3600))
	//         time_t t;
	//         struct tm *p;
	//         long janzone, julyzone;
	//         char janname[10], julyname[10];
	//         t = (time((time_t *)0) / YEAR) * YEAR;
	//         p = localtime(&t);
	//         janzone = -p.tm_gmtoff;
	//         strncpy(janname, p.tm_zone ? p.tm_zone : "   ", 9);
	//         janname[9] = '\0';
	//         t += YEAR/2;
	//         p = localtime(&t);
	//         julyzone = -p.tm_gmtoff;
	//         strncpy(julyname, p.tm_zone ? p.tm_zone : "   ", 9);
	//         julyname[9] = '\0';

	//         if( janzone < julyzone ) {
	//             /* DST is reversed in the southern hemisphere */
	//             PyModule_AddIntConstant(m, "timezone", julyzone);
	//             PyModule_AddIntConstant(m, "altzone", janzone);
	//             PyModule_AddIntConstant(m, "daylight",
	//                                     janzone != julyzone);
	//             PyModule_AddObject(m, "tzname",
	//                                Py_BuildValue("(zz)",
	//                                              julyname, janname));
	//         } else {
	//             PyModule_AddIntConstant(m, "timezone", janzone);
	//             PyModule_AddIntConstant(m, "altzone", julyzone);
	//             PyModule_AddIntConstant(m, "daylight",
	//                                     janzone != julyzone);
	//             PyModule_AddObject(m, "tzname",
	//                                Py_BuildValue("(zz)",
	//                                              janname, julyname));
	//         }
	//     }
	// #else
	// #endif /* HAVE_STRUCT_TM_TM_ZONE */
	// #ifdef __CYGWIN__
	//     tzset();
	//     PyModule_AddIntConstant(m, "timezone", _timezone);
	//     PyModule_AddIntConstant(m, "altzone", _timezone-3600);
	//     PyModule_AddIntConstant(m, "daylight", _daylight);
	//     PyModule_AddObject(m, "tzname",
	//                        Py_BuildValue("(zz)", _tzname[0], _tzname[1]));
	// #endif /* __CYGWIN__ */
	// #endif /* !HAVE_TZNAME || __GLIBC__ || __CYGWIN__*/

	// #if defined(HAVE_CLOCK_GETTIME)
	//     PyModule_AddIntMacro(m, CLOCK_REALTIME);
	// #ifdef CLOCK_MONOTONIC
	//     PyModule_AddIntMacro(m, CLOCK_MONOTONIC);
	// #endif
	// #ifdef CLOCK_MONOTONIC_RAW
	//     PyModule_AddIntMacro(m, CLOCK_MONOTONIC_RAW);
	// #endif
	// #ifdef CLOCK_HIGHRES
	//     PyModule_AddIntMacro(m, CLOCK_HIGHRES);
	// #endif
	// #ifdef CLOCK_PROCESS_CPUTIME_ID
	//     PyModule_AddIntMacro(m, CLOCK_PROCESS_CPUTIME_ID);
	// #endif
	// #ifdef CLOCK_THREAD_CPUTIME_ID
	//     PyModule_AddIntMacro(m, CLOCK_THREAD_CPUTIME_ID);
	// #endif
	// #endif /* HAVE_CLOCK_GETTIME */
}

// Initialise the module
func init() {
	methods := []*py.Method{
		py.MustNewMethod("time", time_time, 0, time_doc),
		py.MustNewMethod("clock", time_clock, 0, clock_doc),
		py.MustNewMethod("clock_gettime", time_clock_gettime, 0, clock_gettime_doc),
		py.MustNewMethod("clock_settime", time_clock_settime, 0, clock_settime_doc),
		py.MustNewMethod("clock_getres", time_clock_getres, 0, clock_getres_doc),
		py.MustNewMethod("sleep", time_sleep, 0, sleep_doc),
		py.MustNewMethod("gmtime", time_gmtime, 0, gmtime_doc),
		py.MustNewMethod("localtime", time_localtime, 0, localtime_doc),
		py.MustNewMethod("asctime", time_asctime, 0, asctime_doc),
		py.MustNewMethod("ctime", time_ctime, 0, ctime_doc),
		py.MustNewMethod("mktime", time_mktime, 0, mktime_doc),
		py.MustNewMethod("strftime", time_strftime, 0, strftime_doc),
		py.MustNewMethod("strptime", time_strptime, 0, strptime_doc),
		py.MustNewMethod("tzset", time_tzset, 0, tzset_doc),
		py.MustNewMethod("monotonic", time_monotonic, 0, monotonic_doc),
		py.MustNewMethod("process_time", time_process_time, 0, process_time_doc),
		py.MustNewMethod("perf_counter", time_perf_counter, 0, perf_counter_doc),
		py.MustNewMethod("get_clock_info", time_get_clock_info, 0, get_clock_info_doc),
	}
	globals := py.StringDict{
		//"version": py.Int(MARSHAL_VERSION),
	}
	py.NewModule("time", module_doc, methods, globals)

}

const module_doc = `This module provides various functions to manipulate time values.

There are two standard representations of time.  One is the number
of seconds since the Epoch, in UTC (a.k.a. GMT).  It may be an integer
or a floating point number (to represent fractions of seconds).
The Epoch is system-defined; on Unix, it is generally January 1st, 1970.
The actual value can be retrieved by calling gmtime(0).

The other representation is a tuple of 9 integers giving local time.
The tuple items are:
  year (including century, e.g. 1998)
  month (1-12)
  day (1-31)
  hours (0-23)
  minutes (0-59)
  seconds (0-59)
  weekday (0-6, Monday is 0)
  Julian day (day in the year, 1-366)
  DST (Daylight Savings Time) flag (-1, 0 or 1)
If the DST flag is 0, the time is given in the regular time zone;
if it is 1, the time is given in the DST time zone;
if it is -1, mktime() should guess based on the date and time.

Variables:

timezone -- difference in seconds between UTC and local standard time
altzone -- difference in  seconds between UTC and local DST time
daylight -- whether local time should reflect DST
tzname -- tuple of (standard time zone name, DST time zone name)

Functions:

time() -- return current time in seconds since the Epoch as a float
clock() -- return CPU time since process start as a float
sleep() -- delay for a number of seconds given as a float
gmtime() -- convert seconds since Epoch to UTC tuple
localtime() -- convert seconds since Epoch to local time tuple
asctime() -- convert time tuple to string
ctime() -- convert time in seconds to string
mktime() -- convert local time tuple to seconds since Epoch
strftime() -- convert time tuple to string according to format specification
strptime() -- parse string to time tuple according to format specification
tzset() -- change the local timezone`

// func PyInit_time() {
//     py.Object m;
//     m = PyModule_Create(&timemodule);
//     if (m == nil) {
//         return nil;
//     }

//     /* Set, or reset, module variables like time.timezone */
//     PyInit_timezone(m);

//     if (!initialized) {
//         PyStructSequence_InitType(&StructTimeType,
//                                   &struct_time_type_desc);

//     }
//     Py_INCREF(&StructTimeType);
//     PyModule_AddIntConstant(m, "_STRUCT_TM_ITEMS", 11);
//     PyModule_AddObject(m, "struct_time", (py.Object*) &StructTimeType);
//     initialized = 1;
//     return m;
// }

// func floattime(_Py_clock_info_t *info) py.Object {
// 	var t _PyTime_timeval
// 	var tp timespec
// 	var ret int

// 	/* _PyTime_gettimeofday() does not use clock_gettime()
// 	   because it would require to link Python to the rt (real-time)
// 	   library, at least on Linux */
// 	ret = clock_gettime(CLOCK_REALTIME, &tp)
// 	if ret == 0 {
// 		if info {
// 			var res timespec
// 			info.implementation = "clock_gettime(CLOCK_REALTIME)"
// 			info.monotonic = 0
// 			info.adjustable = 1
// 			if clock_getres(CLOCK_REALTIME, &res) == 0 {
// 				info.resolution = res.tv_sec + res.tv_nsec*1e-9
// 			} else {
// 				info.resolution = 1e-9
// 			}
// 		}
// 		return py.Float(tp.tv_sec + tp.tv_nsec*1e-9)
// 	}
// 	_PyTime_gettimeofday_info(&t, info)
// 	return py.Float(float64(t.tv_sec) + float64(t.tv_usec)*1e-6)
// }
