package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

	var response telegram.SendMessageReqBody
	var err error

	fromTelegram, err := json.Marshal(body)
	if err != nil {
		return err
	}
	fmt.Println("from telegram json:" + string(fromTelegram))

	switch body.GetChatText() {
	case command.StartCommand:
		response = command.SayHello(body)
	case command.HelpCommand:
		response = command.Help(body)
	case command.RuEnCommand:
		response = command.ChangeTranslateTransition(command.RuEnCommand, body)
	case command.EnRuCommand:
		response = command.ChangeTranslateTransition(command.EnRuCommand, body)
	case command.GetAllTopCommand:
		response = command.GetTop10(body)
	case command.GetMyTopCommand:
		response = command.GetTop10ForUser(body)
	default:
		response = command.General(body)
	}
	if len([]rune(response.Text)) > 0 {
		toTelegram, err := json.Marshal(response)
		if err != nil {
			return err
		}
		fmt.Println("----")
		fmt.Println("to telegram json:" + string(toTelegram))
		fmt.Println("+++")
		fmt.Println("+++")

		res, err := http.Post(telegramConfig.GetTelegramUrl(), "application/json", bytes.NewBuffer(toTelegram))
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(res.Body)
			return errors.New("Unexpected status:" + res.Status + " Message:" + string(body))
		}
		telegram.SendVoice(body.GetChatId())
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
