package minion

import (
	"testing"

	"github.com/ustrajunior/minion/tst"
)

func TestLastChar(t *testing.T) {
	var c uint8

	c = lastChar("minion")
	tst.AssertEqual(t, uint8('n'), c)

	c = lastChar("/users/")
	tst.AssertEqual(t, uint8('/'), c)

	c = lastChar("")
	tst.AssertEqual(t, uint8(0), c)
}
