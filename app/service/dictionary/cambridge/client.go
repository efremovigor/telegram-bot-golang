package cambridge

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"strings"
	"telegram-bot-golang/db/redis"
	"unicode"
)

func Get(query string) CambridgeInfo {
	cambridgeInfo := CambridgeInfo{}
	cachedInfo, errGetCache := redis.Get(fmt.Sprintf(redis.InfoCambridgePageKey, strings.ToLower(query)))
	if errGetCache != nil {
		html, err := htmlquery.LoadURL("https://dictionary.cambridge.org/dictionary/english-russian/" + query)
		if err != nil {
			fmt.Println(err)
			return cambridgeInfo
		}
		nodes, _ := htmlquery.QueryAll(html, xpathBLockDescriptionEnRu)
		if len(nodes) > 0 {
			cambridgeInfo.RequestText = strings.TrimSpace(query)
		}
		for _, node := range nodes {
			info := Info{}
			if node, err := htmlquery.Query(node, xpathTitle); err == nil && node != nil {
				info.Text = strings.ToLower(strings.TrimSpace(htmlquery.InnerText(node)))
			}
			pathWayType := ""
			if node, err := htmlquery.Query(node, xpathType); err == nil && node != nil {
				pathWayType = xpathType
			} else if node, err := htmlquery.Query(node, xpathComplexType); err == nil && node != nil {
				pathWayType = xpathComplexType
			}
			if len(pathWayType) > 0 {
				if nodes, err := htmlquery.QueryAll(node, pathWayType+"/*"); err == nil && nodes != nil {
					for _, node := range nodes {
						if node.DataAtom.String() == "span" {
							info.Type += strings.TrimSpace(htmlquery.InnerText(node))
							continue
						}
						break
					}
				}
				info.Type = strings.TrimSpace(info.Type)
			}

			if node, err := htmlquery.Query(node, xpathTranscription); err == nil && node != nil {
				info.Transcription = strings.TrimSpace(htmlquery.InnerText(node))
			}
			if node, err := htmlquery.Query(node, xpathUK); len(cambridgeInfo.VoicePath.UK) == 0 && err == nil && node != nil {
				cambridgeInfo.VoicePath.UK = strings.TrimSpace(htmlquery.SelectAttr(node, "src"))
			}
			if node, err := htmlquery.Query(node, xpathUS); len(cambridgeInfo.VoicePath.US) == 0 && err == nil && node != nil {
				cambridgeInfo.VoicePath.US = strings.TrimSpace(htmlquery.SelectAttr(node, "src"))
			}
			xpathExplanations, _ := htmlquery.QueryAll(node, xpathExplanations)

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
			cambridgeInfo.Options = append(cambridgeInfo.Options, info)
		}

		if infoInJson, err := json.Marshal(cambridgeInfo); err != nil {
			fmt.Println(err)
		} else {
			redis.Set(fmt.Sprintf(redis.InfoCambridgePageKey, cambridgeInfo.RequestText), infoInJson)
		}

	} else {
		fmt.Println("get cambridge info from cache")
		if err := json.Unmarshal([]byte(cachedInfo), &cachedInfo); err != nil {
			fmt.Println(err)
		}
	}

	return cambridgeInfo
}
