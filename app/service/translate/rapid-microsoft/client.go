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
	"telegram-bot-golang/db"
	"telegram-bot-golang/env"
)

const url = "https://microsoft-translator-text.p.rapidapi.com/translate?to=%s&from=%s&api-version=3.0&profanityAction=NoAction&textType=plain"
const cacheKey = "translate_rapid_from_%s_to_%s_text_%s"

func GetTranslate(text string, to string, from string) (microsoftTranslateResponse []MicrosoftTranslate, err error) {
	url := fmt.Sprintf(url, to, from)

	payload := strings.NewReader("[\n    {\n        \"Text\": \"" + text + "\"\n    }\n]")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Host", "microsoft-translator-text.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", env.GetEnvVariable("MICROSOFT_API_TOKEN"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	buf, _ := ioutil.ReadAll(res.Body)
	b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
	if err != nil {
		log.Fatalln(err)
	}

	if err := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(b))).Decode(&microsoftTranslateResponse); err != nil {
		fmt.Println("could not decode microsoft response", err)
	} else {
		db.Set(fmt.Sprintf(cacheKey, from, to, text), string(buf))
	}
	return
}
