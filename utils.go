package minion

// AllRoutes a shortcut to be used on the UnauthenticatedRoutes when
// you want to have all your routes without jwt verification
const AllRoutes = "^.*$"

func lastChar(str string) (lc uint8) {
	size := len(str)
	if size == 0 {
		return lc
	}
	return str[size-1]
}
