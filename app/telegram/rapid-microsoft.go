package telegram

import "fmt"

func GetBlockWithRapidInfo(word string, translate string) string {
	return fmt.Sprintf(
		DecodeForTelegram("✅ *Rapid-microsoft*: ")+"*%s*\n", DecodeForTelegram(word)) + "\n" +
		GetFieldIfCan(translate, "💡 Translate")
}
