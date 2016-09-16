package minion

import (
	"testing"

	"github.com/ustrajunior/minion/tst"
)

func TestCalculateAbsolutePath(t *testing.T) {
	var router *Router
	var path string

	router = &Router{namespace: "/"}
	path = router.calculateAbsolutePath("/users")
	tst.AssertEqual(t, "/users", path)

	router = &Router{namespace: "/api"}
	path = router.calculateAbsolutePath("/users")
	tst.AssertEqual(t, "/api/users", path)

	router = &Router{namespace: "/api"}
	path = router.calculateAbsolutePath("")
	tst.AssertEqual(t, "/api", path)
}
