package multitran

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"strings"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/helper"
)

func Get(query string) Page {
	page := Page{}

	cachedInfo, errGetCache := redis.Get(fmt.Sprintf(redis.InfoMultitranPageKey, query))
	if errGetCache != nil {
		fmt.Println("get multitran page from service")
		html, err := htmlquery.LoadURL(
			fmt.Sprintf("https://www.multitran.com/m.exe?l1=1&l2=2&s=%s&langlist=2",
				strings.NewReplacer(
					" ", "+",
				).Replace(strings.TrimSpace(query)),
			),
		)
		if err != nil {
			fmt.Println("error getting html data: " + err.Error())
			return page
		}
		var option Info
		if rows, err := htmlquery.QueryAll(html, "//table//tr"); err == nil && rows != nil {
			for _, row := range rows {
				if cols, err := htmlquery.QueryAll(row, "//td"); err == nil && cols != nil {
					switch len(cols) {
					case 1:
						if htmlquery.SelectAttr(cols[0], "class") == "gray" {
							if !helper.IsEmpty(option.Text) && len(option.Explanation) > 0 {
								page.Options = append(page.Options, option)
							}
							option = Info{}
							if node, err := htmlquery.Query(cols[0], "/a"); err == nil && node != nil {
								option.Text = htmlquery.InnerText(node)
							}
							if node, err := htmlquery.Query(cols[0], "/span"); err == nil && node != nil {
								option.Transcription = htmlquery.InnerText(node)
							}
							if node, err := htmlquery.Query(cols[0], "/em"); err == nil && node != nil {
								option.Type = htmlquery.InnerText(node)
							}
						}
					case 2:
						if htmlquery.SelectAttr(cols[0], "class") == "subj" && htmlquery.SelectAttr(cols[1], "class") == "trans1" {
							title := strings.ToLower(htmlquery.SelectAttr(htmlquery.FindOne(cols[0], "/a"), "title"))
							if isGeneralType(title) {
								var texts []string
								for _, element := range strings.Split(htmlquery.InnerText(cols[1]), ";") {
									texts = append(texts, strings.TrimSpace(element))
								}
								expl := Explanation{}
								expl.Type = title
								expl.Text = texts
								option.Explanation = append(option.Explanation, expl)
							}
						}
					}
				}
			}
		} else if err != nil {
			fmt.Println(err)
		}
		if !helper.IsEmpty(option.Text) && len(option.Explanation) > 0 {
			page.Options = append(page.Options, option)
		}
		if len(page.Options) > 0 {
			page.RequestText = query
		}

		if infoInJson, err := json.Marshal(page); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(infoInJson))
			redis.Set(fmt.Sprintf(redis.InfoMultitranPageKey, page.RequestText), infoInJson, 0)
		}
	} else {
		fmt.Println("get multitran page from cache")
		if err := json.Unmarshal([]byte(cachedInfo), &page); err != nil {
			fmt.Println(err)
		}
	}
	return page
}
