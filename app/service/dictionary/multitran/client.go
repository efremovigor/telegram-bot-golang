package multitran

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"strings"
	"telegram-bot-golang/helper"
)

func Get(query string) Page {
	info := Page{}

	html, err := htmlquery.LoadURL(fmt.Sprintf("https://www.multitran.com/m.exe?l1=1&l2=2&s=%s&langlist=2", query))

	if err != nil {
		fmt.Println("error getting html data: " + err.Error())
		return info
	}
	rows := htmlquery.Find(html, "//table[1]//tr")
	var option Info
	for _, row := range rows {
		cols := htmlquery.Find(row, "//td")
		switch len(cols) {
		case 1:
			if htmlquery.SelectAttr(cols[0], "class") == "gray" {
				if !helper.IsEmpty(option.Text) && len(option.Explanation) > 0 {
					info.Options = append(info.Options, option)
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
	if !helper.IsEmpty(option.Text) && len(option.Explanation) > 0 {
		info.Options = append(info.Options, option)
	}
	if len(info.Options) > 0 {
		info.RequestText = query
	}
	return info
}
