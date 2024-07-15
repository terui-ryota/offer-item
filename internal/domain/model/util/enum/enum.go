package enum

func Validate(value int, validList []int) bool {
	for _, v := range validList {
		if value == v {
			return true
		}
	}
	return false
}
