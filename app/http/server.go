package http

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
	"telegram-bot-golang/env"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
	"telegram-bot-golang/statistic"
	"telegram-bot-golang/telegram"
	telegramConfig "telegram-bot-golang/telegram/config"
)

type Context struct {
	echo.Context
}

func bindListener(listener telegram.TelegramListener) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("listener", listener)
			return next(c)
		}
	}
}

func Handle(listener telegram.TelegramListener) {
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &Context{Context: c}
			return next(cc)
		}
	})

	e.Use(bindListener(listener))

	e.POST(telegramConfig.GetUrlPrefix(), func(c echo.Context) error {
		cc := c.(*Context)
		body := &telegram.WebhookMessage{}

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

		if err := cc.reply(*body); err != nil {
			fmt.Println("error in sending reply:", err)
			return err
		}

		return cc.JSON(http.StatusOK, "")
	})

	if !env.IsProd() {
		e.GET("/cambridge/:query", func(c echo.Context) error {
			query := c.Param("query")
			cambridgeInfo := cambridge.Get(query)
			if cambridgeInfo.IsValid() {
				statistic.Consider(query, 1)
			}
			return c.JSON(http.StatusOK, cambridgeInfo)
		})
		e.GET("/multitran/:query", func(c echo.Context) error {
			query := c.Param("query")
			cambridgeInfo := multitran.Get(query)
			return c.JSON(http.StatusOK, cambridgeInfo)
		})
		e.Logger.Fatal(e.Start(":443"))
	} else {
		e.Logger.Fatal(e.StartTLS(":443", config.GetCertPath(), config.GetCertKeyPath()))
	}
}

func (c Context) reply(body telegram.WebhookMessage) error {

	fromTelegram, err := json.Marshal(body)
	if err != nil {
		return err
	}
	fmt.Println("from telegram json:" + string(fromTelegram))
	listener := c.Get("listener").(telegram.TelegramListener)
	switch body.GetChatText() {
	case command.StartCommand:
		listener.Msg <- telegram.TelegramTree{Type: "text", Msg: command.SayHello(body)}
	case command.HelpCommand:
		listener.Msg <- telegram.TelegramTree{Type: "text", Msg: command.Help(body)}
	case command.RuEnCommand:
		listener.Msg <- telegram.TelegramTree{Type: "text", Msg: command.ChangeTranslateTransition(command.RuEnCommand, body)}
	case command.EnRuCommand:
		listener.Msg <- telegram.TelegramTree{Type: "text", Msg: command.ChangeTranslateTransition(command.EnRuCommand, body)}
	case command.GetAllTopCommand:
		listener.Msg <- telegram.TelegramTree{Type: "text", Msg: command.GetTop10(body)}
	case command.GetMyTopCommand:
		listener.Msg <- telegram.TelegramTree{Type: "text", Msg: command.GetTop10ForUser(body)}
	default:
		for _, message := range command.General(body) {
			listener.Msg <- message
		}
	}

	return nil
}
