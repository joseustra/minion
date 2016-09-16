package minion

import (
	"testing"

	"github.com/ustrajunior/minion/tst"
)

func TestCalculateAbsolutePath(t *testing.T) {
	var router *Router
	var path string

	router = &Router{absolutePath: "/"}
	path = router.calculateAbsolutePath("/users")
	tst.AssertEqual(t, "/users", path)

	router = &Router{absolutePath: "/api"}
	path = router.calculateAbsolutePath("/users")
	tst.AssertEqual(t, "/api/users", path)

	router = &Router{absolutePath: "/api"}
	path = router.calculateAbsolutePath("")
	tst.AssertEqual(t, "/api", path)
}
