package cambridge

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"strings"
	"telegram-bot-golang/db"
	"unicode"
)

const cacheKey = "info_cambridge_page_%s"

func Get(query string) Info {
	info := Info{}

	cachedInfo, errGetCache := db.Get(fmt.Sprintf(cacheKey, strings.ToLower(query)))
	if errGetCache != nil {
		fmt.Println("get cambridge info from service")
		html, err := htmlquery.LoadURL("https://dictionary.cambridge.org/dictionary/english-russian/" + query + "?q=" + query)
		if err != nil {
			fmt.Println(err)
			return info
		}

		if node, err := htmlquery.Query(html, xpathTitle); err == nil && node != nil {
			info.Text = strings.TrimSpace(htmlquery.InnerText(node))
		}
		if node, err := htmlquery.Query(html, xpathType); err == nil && node != nil {
			info.Type = strings.TrimSpace(htmlquery.InnerText(node))
		}
		if node, err := htmlquery.Query(html, xpathTranscription); err == nil && node != nil {
			info.Transcription = strings.TrimSpace(htmlquery.InnerText(node))
		}
		xpathExplanations, err := htmlquery.QueryAll(html, xpathExplanations)

		for _, xpathExplanation := range xpathExplanations {
			explanation := Explanation{}
			if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsSemanticDescription); node != nil && err == nil {
				explanation.SemanticDescription = strings.TrimSpace(
					strings.Map(func(letter rune) rune {
						if unicode.IsGraphic(letter) && unicode.IsPrint(letter) {
							return letter
						}
						return -1
					}, htmlquery.InnerText(node)))
			}
			if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsLevel); node != nil && err == nil {
				explanation.Level = strings.TrimSpace(htmlquery.InnerText(node))
			}
			if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsDescription); node != nil && err == nil {
				explanation.Description = strings.TrimSpace(htmlquery.InnerText(node))

			}
			if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsTranslate); node != nil && err == nil {
				explanation.Translate = strings.TrimSpace(htmlquery.InnerText(node))
			}

			if xpathExamples, err := htmlquery.QueryAll(html, xpathExplanationsExamples); xpathExamples != nil && err == nil && len(xpathExamples) > 0 {
				for _, xpathExample := range xpathExamples {
					explanation.Example = append(explanation.Example, strings.TrimSpace(htmlquery.InnerText(xpathExample)))
				}
			}
			info.Explanation = append(info.Explanation, explanation)
		}
		if json, err := json.Marshal(info); err != nil {
			fmt.Println(err)
		} else {
			db.Set(fmt.Sprintf(cacheKey, strings.ToLower(query)), json)
		}

	} else {
		fmt.Println("get cambridge info from cache")
		if err := json.Unmarshal([]byte(cachedInfo), &info); err != nil {
			fmt.Println(err)
		}
	}

	return info
}
