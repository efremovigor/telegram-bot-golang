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
		page.Transcription["us"] = strings.TrimSpace(strings.ToLower(htmlquery.InnerText(value)))
	}

	return
}
