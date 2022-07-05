package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

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
	rus := []rune{'й', 'ц', 'у', 'к', 'е', 'н', 'г', 'ш', 'щ', 'з', 'х', 'ъ', 'ф', 'ы', 'в', 'а', 'п', 'р', 'о', 'л', 'д', 'ж', 'э', 'я', 'ч', 'с', 'м', 'и', 'т', 'ь', 'б', 'ю'}
	var count int
	for _, letter := range []rune(text) {
		for _, en := range ens {
			if en == letter {
				count++
				continue
			}
		}
		for _, ru := range rus {
			if ru == letter {
				count--
				continue
			}
		}
	}
	if count > 0 {
		return true
	}
	return false
}

func ToJson(object interface{}) string {
	var out string
	if requestTelegramInJson, err := json.Marshal(object); err == nil {
		out = string(requestTelegramInJson)
	} else {
		fmt.Println(err)
	}
	return out
}

func ParseJson(r io.Reader, object interface{}) (string, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	if err = json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buf))).Decode(object); err != nil {
		fmt.Println("could not decode microsoft response", err)
	}
	b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
	return string(b), err
}
