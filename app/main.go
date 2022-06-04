package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"telegram-bot-golang/command"
	"telegram-bot-golang/config"
	"telegram-bot-golang/telegram"
	telegramConfig "telegram-bot-golang/telegram/config"
)

func reply(body telegram.WebhookReqBody) error {

	fromTelegram, err := json.Marshal(body)
	if err != nil {
		return err
	}
	fmt.Println("from telegram json:" + string(fromTelegram))

	switch body.GetChatText() {
	case command.StartCommand:
		command.SayHello(body)
	case command.HelpCommand:
		command.Help(body)
	case command.RuEnCommand:
		command.ChangeTranslateTransition(command.RuEnCommand, body)
	case command.EnRuCommand:
		command.ChangeTranslateTransition(command.EnRuCommand, body)
	case command.GetAllTopCommand:
		command.GetTop10(body)
	case command.GetMyTopCommand:
		command.GetTop10ForUser(body)
	default:
		command.General(body)
	}

	return nil
}

func main() {
	e := echo.New()
	e.POST(telegramConfig.GetUrlPrefix(), func(c echo.Context) error {
		body := &telegram.WebhookReqBody{}

		buf, _ := ioutil.ReadAll(c.Request().Body)
		b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("raw json from telegram:" + string(b))
		fmt.Println("----")

		if err := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buf))).Decode(body); err != nil {
			fmt.Println("could not decode request body", err)
			return err
		}

		if err := reply(*body); err != nil {
			fmt.Println("error in sending reply:", err)
			return err
		}

		return c.JSON(http.StatusOK, "")
	})

	//e.GET("/dictionary/:query", func(c echo.Context) error {
	//	query := c.Param("query")
	//	cambridgeInfo := cambridge.Get(query)
	//	if cambridgeInfo.IsValid() {
	//		statistic.Consider(query, 1)
	//	}
	//	return c.JSON(http.StatusOK, cambridgeInfo)
	//})

	e.Logger.Fatal(e.StartTLS(":443", config.GetCertPath(), config.GetCertKeyPath()))
	//e.Logger.Fatal(e.Start(":443"))
}
