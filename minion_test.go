package minion_test

import (
	"testing"

	"github.com/ustrajunior/minion"
	"github.com/ustrajunior/minion/tst"
)

func TestVersion(t *testing.T) {
	tst.AssertEqual(t, "0.0.1", minion.Version())
}

func TestNew(t *testing.T) {
	m := minion.New(minion.Options{})

	tst.AssertNotNil(t, m)
	tst.AssertNotNil(t, m.Router)
}

func TestClassic(t *testing.T) {
	m := minion.Classic(minion.Options{})

	tst.AssertNotNil(t, m)
	tst.AssertNotNil(t, m.Router)
}

func TestDefaultEnv(t *testing.T) {
	tst.AssertEqual(t, minion.DEV, minion.MinionEnv)
}
