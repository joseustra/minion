package minion

func lastChar(str string) (lc uint8) {
	size := len(str)
	if size == 0 {
		return lc
	}
	return str[size-1]
}
