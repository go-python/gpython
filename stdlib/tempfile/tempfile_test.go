// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tempfile_test

import (
	"testing"

	"github.com/go-python/gpython/pytest"
)

func TestTempfile(t *testing.T) {
	pytest.RunScript(t, "./testdata/test.py")
}
