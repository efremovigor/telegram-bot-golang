package rapid_microsoft

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/env"
	"telegram-bot-golang/helper"
)

const url = "https://microsoft-translator-text.p.rapidapi.com/translate?to=%s&from=%s&api-version=3.0&profanityAction=NoAction&textType=plain"

func GetTranslate(text string, to string, from string) string {
	var microsoftTranslateResponse []MicrosoftTranslate
	var err error
	translate, errGetCache := redis.Get(fmt.Sprintf(redis.TranslateRapidMicrosoftKey, from, to, text))
	if errGetCache != nil {
		fmt.Println("get translate from service")

		url := fmt.Sprintf(url, to, from)
		payload := strings.NewReader("[\n    {\n        \"Text\": \"" + text + "\"\n    }\n]")
		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("content-type", "application/json")
		req.Header.Add("X-RapidAPI-Host", "microsoft-translator-text.p.rapidapi.com")
		req.Header.Add("X-RapidAPI-Key", env.GetEnvVariable("MICROSOFT_API_TOKEN"))

		res, _ := http.DefaultClient.Do(req)

		defer CloseConnection(res.Body)

		parseJson, err := helper.ParseJson(res.Body, &microsoftTranslateResponse)
		if err != nil {
			fmt.Println("could not decode microsoft response", err)
		} else {
			redis.Set(fmt.Sprintf(redis.TranslateRapidMicrosoftKey, from, to, text), parseJson, 0)
		}

	} else {
		fmt.Println("get translate from redis")

		if err = json.Unmarshal([]byte(translate), &microsoftTranslateResponse); err != nil {
			fmt.Println(err)
		}
	}

	stringTranslation := ""
	if err == nil {
		for i, response := range microsoftTranslateResponse {
			for _, translation := range response.Translations {
				if i != 0 {
					stringTranslation += ", "
				}
				stringTranslation += translation.Text
			}
		}
	}

	return strings.ToLower(stringTranslation)
}

func CloseConnection(connect io.ReadCloser) {
	err := connect.Close()
	if err != nil {
		fmt.Println("rapid:error close connection:" + err.Error())
	}
}
