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
	"telegram-bot-golang/config"
	"telegram-bot-golang/db"
	"telegram-bot-golang/service/dictionary/cambridge"
	telegram "telegram-bot-golang/telegram"
	telegramConfig "telegram-bot-golang/telegram/config"
)

func sayPolo(body telegram.WebhookReqBody) error {

	var response telegram.SendMessageReqBody
	var err error

	fromTelegram, err := json.Marshal(body)
	if err != nil {
		return err
	}
	fmt.Println("from telegram json:" + string(fromTelegram))

	switch body.GetChatText() {
	case "/start":
		response = telegram.GetTelegramRequest(
			body.GetChatId(),
			telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
		)
	case "/ru_en":
		db.Set(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()), "ru_en")
		response = telegram.GetTelegramRequest(
			body.GetChatId(),
			telegram.GetBaseMsg(body.GetUsername(), body.GetUserId())+telegram.GetChangeTranslateMsg("RU -> EN"),
		)

	case "/en_ru":
		db.Set(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()), "en_ru")
		response = telegram.GetTelegramRequest(
			body.GetChatId(),
			telegram.GetBaseMsg(body.GetUsername(), body.GetUserId())+telegram.GetChangeTranslateMsg("EN -> RU"),
		)
	default:
		fmt.Println(fmt.Sprintf("chat text: %s", body.GetChatText()))
		state, _ := db.Get(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()))
		response = telegram.Reply(body, state)
	}

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

		if err := sayPolo(*body); err != nil {
			fmt.Println("error in sending reply:", err)
			return err
		}

		return c.JSON(http.StatusOK, "")
	})

	e.GET("/dictionary/:query", func(c echo.Context) error {
		query := c.Param("query")
		return c.JSON(http.StatusOK, cambridge.Get(query))
	})

	e.Logger.Fatal(e.StartTLS(":443", config.GetCertPath(), config.GetCertKeyPath()))
	//e.Logger.Fatal(e.Start(":443"))
}

func getHtmlPage(webPage string) (string, error) {

	resp, err := http.Get(webPage)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		return "", err
	}

	return string(body), nil
}
