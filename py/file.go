// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// File object
//
// FIXME cpython 3.3 has a compicated heirachy of types to implement
// this which we do not emulate yet

package py

import (
	"io"
	"io/ioutil"
	"os"
)

var FileType = NewType("file", `represents an open file`)
var errClosed = ExceptionNewf(ValueError, "I/O operation on closed file.")

func init() {
	FileType.Dict["write"] = MustNewMethod("write", func(self Object, value Object) (Object, error) {
		return self.(*File).Write(value)
	}, 0, "write(arg) -> writes the contents of arg to the file, returning the number of characters written.")

	FileType.Dict["read"] = MustNewMethod("read", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
		return self.(*File).Read(args, kwargs)
	}, 0, "read([size]) -> read at most size bytes, returned as a string.\n\nIf the size argument is negative or omitted, read until EOF is reached.\nNotice that when in non-blocking mode, less data than what was requested\nmay be returned, even if no size parameter was given.")
	FileType.Dict["close"] = MustNewMethod("close", func(self Object) (Object, error) {
		return self.(*File).Close()
	}, 0, "close() -> None or (perhaps) an integer.  Close the file.\n\nSets data attribute .closed to True.  A closed file cannot be used for\nfurther I/O operations.  close() may be called more than once without\nerror.  Some kinds of file objects (for example, opened by popen())\nmay return an exit status upon closing.")
	FileType.Dict["flush"] = MustNewMethod("flush", func(self Object) (Object, error) {
		return self.(*File).Flush()
	}, 0, "flush() -> Flush the write buffers of the stream if applicable. This does nothing for read-only and non-blocking streams.")
}

type FileMode int

const (
	FileRead   FileMode = 0x01
	FileWrite  FileMode = 0x02
	FileText   FileMode = 0x4000
	FileBinary FileMode = 0x8000

	FileReadWrite = FileRead + FileWrite
)

type File struct {
	*os.File
	FileMode
}

// Type of this object
func (o *File) Type() *Type {
	return FileType
}

func (o *File) Can(mode FileMode) bool {
	return o.FileMode&mode == mode
}

func (o *File) Write(value Object) (Object, error) {
	var b []byte

	switch v := value.(type) {
	// FIXME Bytearray
	case Bytes:
		b = v

	case String:
		b = []byte(v)

	default:
		return nil, ExceptionNewf(TypeError, "expected a string or other character buffer object")
	}

	n, err := o.File.Write(b)
	if err != nil && err.(*os.PathError).Err == os.ErrClosed {
		return nil, errClosed
	}
	return Int(n), err
}

func (o *File) readResult(b []byte) (Object, error) {
	if o.Can(FileBinary) {
		if b != nil {
			return Bytes(b), nil
		}

		return Bytes{}, nil
	}

	if b != nil {
		return String(b), nil
	}

	return String(""), nil
}

func (o *File) Read(args Tuple, kwargs StringDict) (Object, error) {
	var arg Object = None

	err := UnpackTuple(args, kwargs, "read", 0, 1, &arg)
	if err != nil {
		return nil, err
	}

	var r io.Reader = o.File

	switch pyN, ok := arg.(Int); {
	case arg == None:
		// read all

	case ok:
		// number of bytes to read
		// 0: read nothing
		// < 0: read all
		// > 0: read n
		n, _ := pyN.GoInt64()
		if n == 0 {
			return o.readResult(nil)
		}
		if n > 0 {
			r = io.LimitReader(r, n)
		}

	default:
		// invalid type
		return nil, ExceptionNewf(TypeError, "read() argument 1 must be int, not %s", arg.Type().Name)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		if err == io.EOF {
			return o.readResult(nil)
		}
		if perr, ok := err.(*os.PathError); ok && perr.Err == os.ErrClosed {
			return nil, errClosed
		}

		return nil, err
	}

	return o.readResult(b)
}

func (o *File) Close() (Object, error) {
	_ = o.File.Close()
	return None, nil
}

func (o *File) Flush() (Object, error) {
	err := o.File.Sync()
	if perr, ok := err.(*os.PathError); ok && perr.Err == os.ErrClosed {
		return nil, errClosed
	}

	return None, nil
}

func (o *File) M__enter__() (Object, error) {
	return o, nil
}

func (o *File) M__exit__(exc_type, exc_value, traceback Object) (Object, error) {
	return o.Close()
}

func OpenFile(filename, mode string, buffering int) (Object, error) {
	var fileMode FileMode
	var truncate bool
	var exclusive bool

	for _, m := range mode {
		switch m {
		case 'r':
			if fileMode&FileReadWrite != 0 {
				return nil, ExceptionNewf(ValueError, "must have exactly one of create/read/write/append mode")
			}
			fileMode |= FileRead

		case 'w':
			if fileMode&FileReadWrite != 0 {
				return nil, ExceptionNewf(ValueError, "must have exactly one of create/read/write/append mode")
			}
			fileMode |= FileWrite
			truncate = true

		case 'x':
			if fileMode&FileReadWrite != 0 {
				return nil, ExceptionNewf(ValueError, "must have exactly one of create/read/write/append mode")
			}
			fileMode |= FileWrite
			exclusive = true

		case 'a':
			if fileMode&FileReadWrite != 0 {
				return nil, ExceptionNewf(ValueError, "must have exactly one of create/read/write/append mode")
			}
			fileMode |= FileWrite
			truncate = false

		case '+':
			if fileMode&FileReadWrite == 0 {
				return nil, ExceptionNewf(ValueError, "Must have exactly one of create/read/write/append mode and at most one plus")
			}

			truncate = (fileMode & FileWrite) != 0
			fileMode |= FileReadWrite

		case 'b':
			if fileMode&FileReadWrite == 0 {
				return nil, ExceptionNewf(ValueError, "Must have exactly one of create/read/write/append mode and at most one plus")
			}

			if fileMode&FileText != 0 {
				return nil, ExceptionNewf(ValueError, "can't have text and binary mode at once")
			}

			fileMode |= FileBinary

		case 't':
			if fileMode&FileReadWrite == 0 {
				return nil, ExceptionNewf(ValueError, "Must have exactly one of create/read/write/append mode and at most one plus")
			}

			if fileMode&FileBinary != 0 {
				return nil, ExceptionNewf(ValueError, "can't have text and binary mode at once")
			}

			fileMode |= FileText
		}
	}

	var fmode int

	switch fileMode & FileReadWrite {
	case FileReadWrite:
		fmode = os.O_RDWR

	case FileRead:
		fmode = os.O_RDONLY

	case FileWrite:
		fmode = os.O_WRONLY
	}

	if exclusive {
		fmode |= os.O_EXCL
	}

	if truncate {
		fmode |= os.O_CREATE | os.O_TRUNC
	} else {
		fmode |= os.O_APPEND
	}

	f, err := os.OpenFile(filename, fmode, 0666)
	if err != nil {
		switch {
		case os.IsExist(err):
			return nil, ExceptionNewf(FileExistsError, err.Error())

		case os.IsNotExist(err):
			return nil, ExceptionNewf(FileNotFoundError, err.Error())
		}

		return nil, ExceptionNewf(OSError, err.Error())
	}

	if finfo, err := f.Stat(); err == nil {
		if finfo.IsDir() {
			f.Close()
			return nil, ExceptionNewf(IsADirectoryError, "Is a directory: '%s'", filename)
		}
	}

	return &File{f, fileMode}, nil
}

// Check interface is satisfied
var _ I__enter__ = (*File)(nil)
var _ I__exit__ = (*File)(nil)
