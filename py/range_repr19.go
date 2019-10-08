// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.10
// Range object

package py

import "bytes"

func (r *Range) repr() (Object, error) {
	var b bytes.Buffer
	b.WriteString("range(")
	start, err := ReprAsString(r.Start)
	if err != nil {
		return nil, err
	}
	stop, err := ReprAsString(r.Stop)
	if err != nil {
		return nil, err
	}
	b.WriteString(start)
	b.WriteString(", ")
	b.WriteString(stop)

	if r.Step != 1 {
		step, err := ReprAsString(r.Step)
		if err != nil {
			return nil, err
		}
		b.WriteString(", ")
		b.WriteString(step)
	}
	b.WriteString(")")

	return String(b.String()), nil
}
