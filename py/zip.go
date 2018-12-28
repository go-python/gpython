// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// A python Zip object
type Zip struct {
	itTuple Tuple
	size    int
}

// A python ZipIterator iterator
type ZipIterator struct {
	zip Zip
}

var ZipType = NewTypeX("zip", `zip(iter1 [,iter2 [...]]) --> zip object

Return a zip object whose .__next__() method returns a tuple where
the i-th element comes from the i-th iterable argument.  The .__next__()
method continues until the shortest iterable in the argument sequence
is exhausted and then it raises StopIteration.`,
	ZipTypeNew, nil)

// Type of this object
func (z *Zip) Type() *Type {
	return ZipType
}

// ZipTypeNew
func ZipTypeNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	tupleSize := len(args)
	itTuple := make(Tuple, tupleSize)
	for i := 0; i < tupleSize; i++ {
		item := args[i]
		iter, err := Iter(item)
		if err != nil {
			return nil, ExceptionNewf(TypeError, "zip argument #%d must support iteration", i+1)
		}
		itTuple[i] = iter
	}

	return &Zip{itTuple: itTuple, size: tupleSize}, nil
}

// Zip iterator
func (z *Zip) M__iter__() (Object, error) {
	return z, nil
}

func (z *Zip) M__next__() (Object, error) {
	result := make(Tuple, z.size)
	for i := 0; i < z.size; i++ {
		value, err := Next(z.itTuple[i])
		if err != nil {
			return nil, err
		}
		result[i] = value
	}
	return result, nil
}

// Check interface is satisfied
var _ I__iter__ = (*Zip)(nil)
var _ I__next__ = (*Zip)(nil)
