package rapid_microsoft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/env"
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

		buf, _ := ioutil.ReadAll(res.Body)
		b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
		if err != nil {
			log.Fatalln(err)
		}

		if err = json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(b))).Decode(&microsoftTranslateResponse); err != nil {
			fmt.Println("could not decode microsoft response", err)
		} else {
			redis.Set(fmt.Sprintf(redis.TranslateRapidMicrosoftKey, from, to, text), buf)
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

	return stringTranslation
}

func CloseConnection(connect io.ReadCloser) {
	err := connect.Close()
	if err != nil {
		fmt.Println("rapid:error close connection:" + err.Error())
	}
}
