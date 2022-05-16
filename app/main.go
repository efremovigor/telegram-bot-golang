package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"telegram-bot-golang/config"
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

	switch body.Message.Text {
	case "/start":
		response = telegram.GetTelegramRequest(
			body.Message.Chat.ID,
			telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
		)
	case "/ru_en":
		_, exist := telegram.Chats[body.Message.Chat.ID]
		if !exist {
			telegram.Chats[body.Message.Chat.ID] = make(map[int]string)
		}
		telegram.Chats[body.Message.Chat.ID][body.Message.From.ID] = "ru_en"
		response = telegram.GetTelegramRequest(
			body.Message.Chat.ID,
			telegram.GetBaseMsg(body.Message.From.FirstName, body.Message.From.ID)+telegram.GetChangeTranslateMsg("RU -> EN"),
		)

	case "/en_ru":
		_, exist := telegram.Chats[body.Message.Chat.ID]
		if !exist {
			telegram.Chats[body.Message.Chat.ID] = make(map[int]string)
		}
		telegram.Chats[body.Message.Chat.ID][body.Message.From.ID] = "en_ru"
		response = telegram.GetTelegramRequest(
			body.Message.Chat.ID,
			telegram.GetBaseMsg(body.Message.From.FirstName, body.Message.From.ID)+telegram.GetChangeTranslateMsg("EN -> RU"),
		)
	default:
		fmt.Println(fmt.Sprintf("chat text: %s", body.Message.Text))
		response = telegram.SayHello(body)

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

		if err := json.NewDecoder(c.Request().Body).Decode(body); err != nil {
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
		html, err := htmlquery.LoadURL("https://dictionary.cambridge.org/dictionary/english-russian/" + query + "?q=" + query)
		if err != nil {
			panic(err)
		}

		list, err := htmlquery.QueryAll(html, "//div[contains(@class, 'entry-body')]//div[contains(@class, 'entry-body__el')]//span[@lang=\"ru\"]")
		var translate string

		for _, n := range list {
			translate += htmlquery.InnerText(n)
		}
		return c.JSON(http.StatusOK, translate)
	})
	e.Logger.Fatal(e.StartTLS(":443", config.GetCertPath(), config.GetCertKeyPath()))
	//e.Logger.Fatal(e.Start(":88"))
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
