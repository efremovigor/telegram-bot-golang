package config

import (
	"fmt"
	"telegram-bot-golang/env"
)

const sendUrl = "https://api.telegram.org/bot%s/sendMessage"

func GetTelegramUrl() string {
	return fmt.Sprintf(sendUrl, env.GetEnvVariable("TELEGRAM_API_TOKEN"))
}
func GetUrlPrefix() string {
	return "/" + env.GetEnvVariable("TELEGRAM_API_TOKEN") + "/"
}
