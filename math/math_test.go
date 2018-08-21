package math_test

import (
	"testing"

	"github.com/go-python/gpython/pytest"
)

func TestVm(t *testing.T) {
	pytest.RunTests(t, "tests")
}
