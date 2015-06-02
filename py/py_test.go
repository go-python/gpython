package py_test

import (
	"testing"

	"github.com/ncw/gpython/pytest"
)

func TestPy(t *testing.T) {
	pytest.RunTests(t, "tests")
}
