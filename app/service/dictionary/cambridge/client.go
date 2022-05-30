package cambridge

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"strings"
	"unicode"
)

func Get(query string) Info {
	info := Info{}

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

	return info
}
