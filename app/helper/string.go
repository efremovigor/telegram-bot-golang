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

func IsEn(text string) bool {
	ens := []rune{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', 'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'z', 'x', 'c', 'v', 'b', 'n', 'm'}
	var count int
	for _, letter := range []rune(text) {
		for _, en := range ens {
			if en == letter {
				count++
				continue
			}
		}
	}
	if count > Len(text)/2 {
		return true
	}
	return false
}
