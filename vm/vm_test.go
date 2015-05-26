package vm_test

import (
	"testing"

	"github.com/ncw/gpython/pytest"
)

func TestVm(t *testing.T) {
	pytest.RunTests(t, "tests")
}
