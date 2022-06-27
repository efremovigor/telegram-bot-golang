package cambridge

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"net/http"
	"strings"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/helper"
	"unicode"
)

func Get(query string) Page {
	page := Page{}
	cachedInfo, errGetCache := redis.Get(fmt.Sprintf(redis.InfoCambridgePageKey, query))
	if errGetCache != nil {
		html, err := htmlquery.LoadURL(Url + "/dictionary/english-russian/" + query)
		if err != nil {
			fmt.Println("error getting html data: " + err.Error())
			return page
		}
		nodes, _ := htmlquery.QueryAll(html, xpathBlockDescriptionEnRu)

		if len(nodes) == 0 {
			nodes, _ = htmlquery.QueryAll(html, xpathAltBlockDescriptionEnRu)
			if len(nodes) == 0 {
				html, err = htmlquery.LoadURL(Url + "/dictionary/english/" + query)
				if err != nil {
					fmt.Println("error getting html data: " + err.Error())
					return page
				}
				nodes, _ = htmlquery.QueryAll(html, xpathBlockDescriptionEnRu)
				if len(nodes) == 0 {
					nodes, _ = htmlquery.QueryAll(html, xpathAltBlockDescriptionEnRu)
				}
			}
		}
		for _, node := range nodes {
			info := Info{}
			if node, err := htmlquery.Query(node, xpathTitle); err == nil && node != nil {
				info.Text = strings.ToLower(strings.TrimSpace(htmlquery.InnerText(node)))
			}
			pathWayType := ""
			if nodeWordType, err := htmlquery.Query(node, xpathType); err == nil && nodeWordType != nil {
				pathWayType = xpathType
			} else if nodeWordType, err := htmlquery.Query(node, xpathComplexType); err == nil && nodeWordType != nil {
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
			if node, err := htmlquery.Query(node, xpathUK); helper.Len(page.VoicePath.UK) == 0 && err == nil && node != nil {
				page.VoicePath.UK = strings.TrimSpace(htmlquery.SelectAttr(node, "src"))
			}
			if node, err := htmlquery.Query(node, xpathUS); helper.Len(page.VoicePath.US) == 0 && err == nil && node != nil {
				page.VoicePath.US = strings.TrimSpace(htmlquery.SelectAttr(node, "src"))
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
				if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsWord); node != nil && err == nil {
					explanation.Text = strings.TrimSpace(htmlquery.InnerText(node))
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

				if xpathExamples, err := htmlquery.QueryAll(xpathExplanation, xpathExplanationsExamples); xpathExamples != nil && err == nil && len(xpathExamples) > 0 {
					for _, xpathExample := range xpathExamples {
						explanation.Example = append(explanation.Example, strings.TrimSpace(htmlquery.InnerText(xpathExample)))
					}
				}
				info.Explanation = append(info.Explanation, explanation)
			}
			page.Options = append(page.Options, info)
		}
		if len(page.Options) > 0 {
			page.RequestText = strings.TrimSpace(query)
		}

		if infoInJson, err := json.Marshal(page); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(infoInJson))
			redis.Set(fmt.Sprintf(redis.InfoCambridgePageKey, page.RequestText), infoInJson, 0)
		}
	} else {
		fmt.Println("get cambridge info from cache")
		if err := json.Unmarshal([]byte(cachedInfo), &page); err != nil {
			fmt.Println(err)
		}
	}

	return page
}

func Search(query string) (response SearchResponse) {
	res, err := http.Get(Url + "/dictionary/english-russian/" + query)

	if err != nil {
		fmt.Println("error getting search of result: " + err.Error())
		return
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil || !response.IsValid() {
		return
	}

	response.RequestWord = query
	return
}
