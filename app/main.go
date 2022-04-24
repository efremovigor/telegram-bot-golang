package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"telegram-bot-golang/env"

	"github.com/labstack/echo/v4"
)

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func Handler(res http.ResponseWriter, req *http.Request) {

}

//The below code deals with the process of sending a response message
// to the user

// Create a struct to conform to the JSON body
// of the send message request
// https://core.telegram.org/bots/api#sendmessage
type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// sayPolo takes a chatID and sends "polo" to them
func sayPolo(chatID int64) error {
	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   "Polo!!",
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot"+env.GetEnvVariable("TELEGRAM_API_TOKEN")+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
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
	fmt.Println(env.GetEnvVariable("TELEGRAM_API_TOKEN"))
	e.POST("/"+env.GetEnvVariable("TELEGRAM_API_TOKEN")+"/", func(c echo.Context) error {
		// First, decode the JSON response body
		body := &webhookReqBody{}
		if err := json.NewDecoder(c.Request().Body).Decode(body); err != nil {
			fmt.Println("could not decode request body", err)
			return err
		}

		// Check if the message contains the word "marco"
		// if not, return without doing anything
		if !strings.Contains(strings.ToLower(body.Message.Text), "marco") {
			if err := sayPolo(body.Message.Chat.ID); err != nil {
				fmt.Println("error in sending reply:", err)
				return err
			}
		}

		// If the text contains marco, call the `sayPolo` function, which
		// is defined below
		if err := sayPolo(body.Message.Chat.ID); err != nil {
			fmt.Println("error in sending reply:", err)
			return err
		}

		// log a confirmation message if the message is sent successfully
		fmt.Println("reply sent")
		return c.JSON(http.StatusOK, "")
	})
	e.Logger.Fatal(e.StartTLS(":443", "./build/domain.crt", "./build/domain.key"))
}
