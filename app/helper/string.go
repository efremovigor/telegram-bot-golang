package helper

func IsEmpty(s string) bool {
	if len([]rune(s)) > 0 {
		return false
	}
	return true
}
func Len(s string) int {
	return len([]rune(s))
}
