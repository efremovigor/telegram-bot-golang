package http

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"telegram-bot-golang/command"
	"telegram-bot-golang/config"
	"telegram-bot-golang/env"
	"telegram-bot-golang/helper"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
	"telegram-bot-golang/service/dictionary/wooordhunt"
	"telegram-bot-golang/statistic"
	"telegram-bot-golang/telegram"
	telegramConfig "telegram-bot-golang/telegram/config"
)

type Context struct {
	echo.Context
}

func (c Context) tryReply(reader io.Reader, message telegram.IncomingTelegramQueryInterface) error {
	if _, err := helper.ParseJson(reader, &message); err == nil && message.IsValid() {
		if err := c.reply(message); err != nil {
			fmt.Println("error in sending reply:", err)
		}
		return nil
	} else {
		return errors.New("invalid struct")
	}
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

		parseJson, _ := ioutil.ReadAll(c.Request().Body)

		if err := cc.tryReply(bytes.NewBuffer(parseJson), &telegram.WebhookMessage{}); err != nil {
			if err = cc.tryReply(bytes.NewBuffer(parseJson), &telegram.CallbackQuery{}); err != nil {
				fmt.Println(err)
			}
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
		e.GET("/multitran/:query", func(c echo.Context) error {
			query := c.Param("query")
			info := multitran.Get(query)
			return c.JSON(http.StatusOK, info)
		})

		e.GET("/wooordhunt/:query", func(c echo.Context) error {
			query := c.Param("query")
			info := wooordhunt.Get(query)
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
				command.SendVoice(query, words[1], strings.Join(words[2:], " "))
			}
			return nil
		case telegram.ShowRequestPic:
			command.SendImage(query, strings.Join(words[1:], " "))
			return nil
		case telegram.SearchRequest:
			if words[1] == "cambridge" {
				query.SetChatText(strings.Join(words[2:], " "))
				command.GetSubCambridge(query)
				if message, err := command.GetNextMessage(query.GetUserId(), query.GetChatText()); err == nil {
					listener.Message <- message
				}
				return nil
			}
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
