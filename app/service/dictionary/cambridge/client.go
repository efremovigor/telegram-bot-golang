package cambridge

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
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
		page = DoRequest(query, Url+"/dictionary/english-russian/"+query, Url+"/dictionary/english/"+query)
		if page.IsValid() {
			redis.SetStruct(fmt.Sprintf(redis.InfoCambridgePageKey, page.RequestText), page, 0)
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
	cachedInfo, errGetCache := redis.Get(fmt.Sprintf(redis.InfoCambridgeSearchResult, query))
	if errGetCache != nil {
		fmt.Println("get cambridge search from service")
		res, err := http.Get(fmt.Sprintf(SearchUrl, query))

		if err != nil {
			fmt.Println("error get cambridge search from service: " + err.Error())
			return SearchResponse{}
		}

		if res.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(fmt.Sprintf(
				"error getting search of result - url:%s, Code:%d, Content:%s",
				fmt.Sprintf(SearchUrl, query),
				res.StatusCode,
				body,
			))
			return
		}

		if _, err = helper.ParseJson(res.Body, &response.Founded); err != nil {
			fmt.Println(err.Error())
			return
		}

		for i, found := range response.Founded {
			if found.Word == query {
				response.Founded[i] = response.Founded[len(response.Founded)-1]
				response.Founded = response.Founded[:len(response.Founded)-1]
				break
			}
		}
		if len(response.Founded) > 0 {
			for _, found := range response.Founded {
				redis.Set(fmt.Sprintf(redis.InfoCambridgeSearchValue, found.Word), found.Path, 0)
			}
			response.RequestWord = query
		}

		redis.SetStruct(fmt.Sprintf(redis.InfoCambridgeSearchResult, response.RequestWord), response, 0)

	} else {
		fmt.Println("get cambridge search from cache")
		if err := json.Unmarshal([]byte(cachedInfo), &response); err != nil {
			fmt.Println(err)
		}
	}
	return
}

func getNodes(url string) []*html.Node {
	node, err := htmlquery.LoadURL(url)
	if err != nil {
		fmt.Println("error getting html data: " + err.Error())
		return []*html.Node{}
	}
	nodes, err := htmlquery.QueryAll(node, xpathBlockDescriptionEnRu)

	if err != nil {
		fmt.Println("error getting html data: " + err.Error())
		return nodes
	}

	if len(nodes) == 0 {
		nodes, _ = htmlquery.QueryAll(node, xpathAltBlockDescriptionEnRu)
	}
	return nodes
}

func DoRequest(word string, url string, altUrl string) (page Page) {
	nodes := getNodes(url)
	if !helper.IsEmpty(altUrl) {
		nodes = append(nodes, getNodes(altUrl)...)
	}

	for _, node := range nodes {
		info := Info{Transcription: make(map[string]string)}
		if node, err := htmlquery.Query(node, xpathTitle); err == nil && node != nil {
			info.Text = strings.ToLower(strings.TrimSpace(htmlquery.InnerText(node)))
		}
		pathWayType := ""
		if nodeWordType, err := htmlquery.Query(node, xpathType); err == nil && nodeWordType != nil {
			pathWayType = xpathType
		} else if nodeWordType, err := htmlquery.Query(node, xpathComplexType); err == nil && nodeWordType != nil {
			pathWayType = xpathComplexType
		} else {
			fmt.Println(err)
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

		if node, err := htmlquery.Query(node, xpathTranscriptionUK); err == nil && node != nil {
			info.Transcription["uk"] = strings.TrimSpace(htmlquery.InnerText(node))
		} else {
			fmt.Println(err)
		}

		if node, err := htmlquery.Query(node, xpathTranscriptionUS); err == nil && node != nil {
			info.Transcription["us"] = strings.TrimSpace(htmlquery.InnerText(node))
		} else {
			fmt.Println(err)
		}

		if forms, err := htmlquery.QueryAll(node, xpathForms); err == nil && len(forms) > 0 {
			for _, form := range forms {
				if desc, err := htmlquery.Query(form, "/span[contains(@class,'lab')]"); err == nil && desc != nil {
					if value, err := htmlquery.Query(form, "/b[contains(@class,'inf')]"); err == nil && value != nil {
						info.Forms = append(
							info.Forms,
							Forms{
								Desc:  strings.TrimSpace(htmlquery.InnerText(desc)),
								Value: strings.TrimSpace(htmlquery.InnerText(value)),
							},
						)
					} else {
						fmt.Println(err)
					}
				} else {
					fmt.Println(err)
				}
			}
		}

		if img, err := htmlquery.Query(node, xpathImage); err == nil && img != nil {
			info.Image = strings.TrimSpace(htmlquery.SelectAttr(img, "src"))
		} else {
			fmt.Println(err)
		}

		if node, err := htmlquery.Query(node, xpathUK); helper.Len(info.VoicePath.UK) == 0 && err == nil && node != nil {
			info.VoicePath.UK = strings.TrimSpace(htmlquery.SelectAttr(node, "src"))
		} else {
			fmt.Println(err)
		}
		if node, err := htmlquery.Query(node, xpathUS); helper.Len(info.VoicePath.US) == 0 && err == nil && node != nil {
			info.VoicePath.US = strings.TrimSpace(htmlquery.SelectAttr(node, "src"))
		} else {
			fmt.Println(err)
		}
		xpathExplanations, err := htmlquery.QueryAll(node, xpathExplanations)
		if err != nil {
			fmt.Println(err)
		}
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
			} else {
				fmt.Println(err)
			}
			if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsWord); node != nil && err == nil {
				explanation.Text = strings.TrimSpace(htmlquery.InnerText(node))
			} else {
				fmt.Println(err)
			}
			if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsLevel); node != nil && err == nil {
				explanation.Level = strings.TrimSpace(htmlquery.InnerText(node))
			} else {
				fmt.Println(err)
			}
			if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsDescription); node != nil && err == nil {
				explanation.Description = strings.TrimSpace(htmlquery.InnerText(node))

			} else {
				fmt.Println(err)
			}
			if node, err := htmlquery.Query(xpathExplanation, xpathExplanationsTranslate); node != nil && err == nil {
				explanation.Translate = strings.TrimSpace(htmlquery.InnerText(node))
			} else {
				fmt.Println(err)
			}

			if xpathExamples, err := htmlquery.QueryAll(xpathExplanation, xpathExplanationsExamples); xpathExamples != nil && err == nil && len(xpathExamples) > 0 {
				for _, xpathExample := range xpathExamples {
					explanation.Example = append(explanation.Example, strings.TrimSpace(htmlquery.InnerText(xpathExample)))
				}
			} else {
				fmt.Println(err)
			}
			if xpathExamples, err := htmlquery.QueryAll(xpathExplanation, xpathExplanationsMoreExamples); xpathExamples != nil && err == nil && len(xpathExamples) > 0 {
				for _, xpathExample := range xpathExamples {
					explanation.Example = append(explanation.Example, strings.TrimSpace(htmlquery.InnerText(xpathExample)))
				}
			} else {
				fmt.Println(err)
			}
			info.Explanation = append(info.Explanation, explanation)
		}
		page.Options = append(page.Options, info)
	}
	if len(page.Options) > 0 {
		page.RequestText = strings.TrimSpace(word)
	}

	return
}
