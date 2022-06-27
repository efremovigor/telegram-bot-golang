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
	"strings"
	"telegram-bot-golang/command"
	"telegram-bot-golang/config"
	"telegram-bot-golang/env"
	"telegram-bot-golang/helper"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
	"telegram-bot-golang/statistic"
	"telegram-bot-golang/telegram"
	telegramConfig "telegram-bot-golang/telegram/config"
)

type Context struct {
	echo.Context
}

func bindListener(listener telegram.Listener) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("listener", listener)
			return next(c)
		}
	}
}

func Handle(listener telegram.Listener) {
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

		buf, _ := ioutil.ReadAll(c.Request().Body)
		b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("raw json from telegram:" + string(b))
		fmt.Println("----")

		webhookMessage := &telegram.WebhookMessage{}
		callbackQuery := &telegram.CallbackQuery{}

		if err := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buf))).Decode(webhookMessage); err == nil && webhookMessage.IsValid() {
			if err := cc.reply(webhookMessage); err != nil {
				fmt.Println("error in sending reply:", err)
				return err
			}
		} else if err := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buf))).Decode(callbackQuery); err == nil && callbackQuery.IsValid() {
			if err := cc.reply(callbackQuery); err != nil {
				fmt.Println("error in sending reply:", err)
				return err
			}
		} else {
			fmt.Println("could not decode request body", err)
		}

		return cc.JSON(http.StatusOK, "")
	})

	if !env.IsProd() {
		e.GET("/cambridge/:query", func(c echo.Context) error {
			query := c.Param("query")
			info := cambridge.Get(query)
			if info.IsValid() {
				statistic.Consider(query, 1)
			}
			return c.JSON(http.StatusOK, info)
		})
		//e.GET("/collins/:query", func(c echo.Context) error {
		//	query := c.Param("query")
		//	info := collins.Get(query)
		//	if info.IsValid() {
		//		statistic.Consider(query, 1)
		//	}
		//	return c.JSON(http.StatusOK, info)
		//})
		e.GET("/multitran/:query", func(c echo.Context) error {
			query := c.Param("query")
			info := multitran.Get(query)
			return c.JSON(http.StatusOK, info)
		})

		e.GET("/is-en/:query", func(c echo.Context) error {
			query := c.Param("query")
			return c.JSON(http.StatusOK, helper.IsEn(query))
		})
		e.Logger.Fatal(e.Start(":443"))
	} else {
		e.Logger.Fatal(e.StartTLS(":443", config.GetCertPath(), config.GetCertKeyPath()))
	}
}

func (c Context) reply(query telegram.IncomingTelegramQueryInterface) error {

	fmt.Println("from telegram message:" + string(query.GetChatText()))
	listener := c.Get("listener").(telegram.Listener)

	switch query.GetChatText() {
	case command.StartCommand:
		listener.Message <- telegram.NewRequestChannelTelegram("text", command.SayHello(query))
	case command.HelpCommand:
		listener.Message <- telegram.NewRequestChannelTelegram("text", command.Help(query))
	case command.RuEnCommand:
		listener.Message <- telegram.NewRequestChannelTelegram("text", command.ChangeTranslateTransition(command.RuEnCommand, query))
	case command.EnRuCommand:
		listener.Message <- telegram.NewRequestChannelTelegram("text", command.ChangeTranslateTransition(command.EnRuCommand, query))
	case command.AutoTranslateCommand:
		listener.Message <- telegram.NewRequestChannelTelegram("text", command.ChangeTranslateTransition(command.AutoTranslateCommand, query))
	case command.GetAllTopCommand:
		listener.Message <- telegram.NewRequestChannelTelegram("text", command.GetTop10(query))
	case command.GetMyTopCommand:
		listener.Message <- telegram.NewRequestChannelTelegram("text", command.GetTop10ForUser(query))
	default:
		words := strings.Split(query.GetChatText(), " ")
		switch words[0] {
		case telegram.NextRequestMessage:
			if message, err := command.GetNextMessage(query.GetUserId(), strings.Join(words[1:], " ")); err == nil {
				listener.Message <- message
			}
			return nil
		case telegram.ShowRequestVoice:
			if words[1] == telegram.CountryUs || words[1] == telegram.CountryUk {
				command.GetVoice(query, words[1], strings.Join(words[2:], " "))
			}
			return nil
		case telegram.SearchRequest:
			query.SetChatText(strings.Join(words[1:], " "))
			command.General(query)
			if message, err := command.GetNextMessage(query.GetUserId(), query.GetChatText()); err == nil {
				listener.Message <- message
			}
			return nil
		}
		command.General(query)
		if message, err := command.GetNextMessage(query.GetUserId(), query.GetChatText()); err == nil {
			listener.Message <- message
		}
	}

	return nil
}
