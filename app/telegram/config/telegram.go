package config

import (
	"fmt"
	"telegram-bot-golang/env"
)

const sendUrl = "https://api.telegram.org/bot%s/%s"

func GetTelegramUrl(method string) string {
	return fmt.Sprintf(sendUrl, env.GetEnvVariable("TELEGRAM_API_TOKEN"), method)
}
func GetUrlPrefix() string {
	return "/" + env.GetEnvVariable("TELEGRAM_API_TOKEN") + "/"
}
