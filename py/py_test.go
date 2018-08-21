package py_test

import (
	"testing"

	"github.com/go-python/gpython/pytest"
)

func TestPy(t *testing.T) {
	pytest.RunTests(t, "tests")
}
