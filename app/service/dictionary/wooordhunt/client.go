package wooordhunt

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"strings"
)

func Get(query string) (page Page) {
	page.Transcription = make(map[string]string)
	html, _ := htmlquery.LoadURL(
		fmt.Sprintf("https://wooordhunt.ru/word/%s",
			strings.NewReplacer(
				" ", "+",
			).Replace(strings.TrimSpace(query)),
		),
	)
	/**
	add is valid block
	*/

	if value, err := htmlquery.Query(html, "//div[@id='wd_title']/h1"); err == nil && value != nil {
		page.RequestText = strings.TrimSpace(strings.ToLower(htmlquery.InnerText(value.FirstChild)))
	}

	if value, err := htmlquery.Query(html, "//div[@id='us_tr_sound']/span[contains(@class, 'transcription')]"); err == nil && value != nil {
		page.Transcription["us"] = strings.TrimSpace(strings.ToLower(htmlquery.InnerText(value)))
	}

	if value, err := htmlquery.Query(html, "//div[@id='uk_tr_sound']/span[contains(@class, 'transcription')]"); err == nil && value != nil {
		page.Transcription["uk"] = strings.TrimSpace(strings.ToLower(htmlquery.InnerText(value)))
	}

	if value, err := htmlquery.Query(html, "//div[@id='us_tr_sound']//source[contains(@type,'audio/mpeg')]"); err == nil && value != nil {
		page.VoicePath.US = strings.TrimSpace(htmlquery.SelectAttr(value, "src"))
	}

	if value, err := htmlquery.Query(html, "//div[@id='uk_tr_sound']//source[contains(@type,'audio/mpeg')]"); err == nil && value != nil {
		page.VoicePath.UK = strings.TrimSpace(htmlquery.SelectAttr(value, "src"))
	}

	explanations, _ := htmlquery.QueryAll(html, "//div[@id='content_in_russian']/*")

	for _, explanation := range explanations {
		if htmlquery.SelectAttr(explanation, "class") == "t_inline_en" {
			page.GeneralTranslate = strings.Split(strings.TrimSpace(strings.ToLower(htmlquery.InnerText(explanation))), ", ")
			continue
		}
		if htmlquery.SelectAttr(explanation, "class") == "tr" {
			explanations, _ := htmlquery.QueryAll(explanation, "/*")
			for _, explanation := range explanations {
				fmt.Println(strings.TrimSpace(strings.ToLower(htmlquery.InnerText(explanation))))
			}
			break
		}
	}

	return
}
