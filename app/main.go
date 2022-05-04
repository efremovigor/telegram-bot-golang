package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"telegram-bot-golang/config"
	telegram "telegram-bot-golang/telegram"
	telegramConfig "telegram-bot-golang/telegram/config"
)

func sayPolo(chatID int, msg string) error {
	reqBody := &telegram.SendMessageReqBody{
		ChatID: chatID,
		Text:   "Иди на хуй со своим:" + msg,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(telegramConfig.GetTelegramUrl(), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func main() {
	e := echo.New()
	e.POST(telegramConfig.GetUrlPrefix(), func(c echo.Context) error {
		body := &telegram.WebhookReqBody{}

		if err := json.NewDecoder(c.Request().Body).Decode(body); err != nil {
			fmt.Println("could not decode request body", err)
			return err
		}
		fmt.Println(body)

		if err := sayPolo(body.Message.Chat.ID, body.Message.Text); err != nil {
			fmt.Println("error in sending reply:", err)
			return err
		}

		return c.JSON(http.StatusOK, "")
	})

	e.Logger.Fatal(e.StartTLS(":443", config.GetCertPath(), config.GetCertKeyPath()))
}
