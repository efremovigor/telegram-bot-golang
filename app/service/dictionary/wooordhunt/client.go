package wooordhunt

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
	"telegram-bot-golang/helper"
)

func Get(query string) (page Page) {
	page.Transcription = make(map[string]string)
	html1, _ := htmlquery.LoadURL(
		fmt.Sprintf("https://wooordhunt.ru/word/%s",
			strings.NewReplacer(
				" ", "+",
			).Replace(cleanTextField(query)),
		),
	)

	if value, err := htmlquery.Query(html1, "//div[@id='wd_title']/h1"); err == nil && value != nil {
		page.Title = cleanNodeField(value.FirstChild)
	}

	if helper.IsEn(query) {
		if value, err := htmlquery.Query(html1, "//div[@id='us_tr_sound']/span[contains(@class, 'transcription')]"); err == nil && value != nil {
			page.Transcription["us"] = cleanNodeField(value)
		}

		if value, err := htmlquery.Query(html1, "//div[@id='uk_tr_sound']/span[contains(@class, 'transcription')]"); err == nil && value != nil {
			page.Transcription["uk"] = cleanNodeField(value)
		}

		if value, err := htmlquery.Query(html1, "//div[@id='us_tr_sound']//source[contains(@type,'audio/mpeg')]"); err == nil && value != nil {
			page.VoicePath.US = cleanTextField(htmlquery.SelectAttr(value, "src"))
		}

		if value, err := htmlquery.Query(html1, "//div[@id='uk_tr_sound']//source[contains(@type,'audio/mpeg')]"); err == nil && value != nil {
			page.VoicePath.UK = cleanTextField(htmlquery.SelectAttr(value, "src"))
		}

		explanations, _ := htmlquery.QueryAll(html1, "//div[@id='content_in_russian']/*")

		for _, explanation := range explanations {
			if explanation.DataAtom.String() == "div" &&
				(htmlquery.SelectAttr(explanation, "class") == "tr" || htmlquery.SelectAttr(explanation, "class") == "block") {
				continue
			}
			if htmlquery.SelectAttr(explanation, "class") == "t_inline_en" {
				page.GeneralTranslate = strings.Split(cleanNodeField(explanation), ", ")
				continue
			}
			if explanation.DataAtom.String() == "h4" {
				nextIsDivTr, _ := htmlquery.Query(explanation, "following-sibling::div[contains(@class,'tr')]")
				if nextIsDivTr != nil {
					wordTypeNode, err := htmlquery.Query(explanation, "node()[1]")
					if err != nil {
						fmt.Println(err)
					}
					page.Info = append(
						page.Info,
						Info{
							Type:    cleanNodeField(wordTypeNode),
							Meaning: getMeaningBlock(nextIsDivTr),
						},
					)
					continue
				}
			}
			if strings.Contains(htmlquery.SelectAttr(explanation, "class"), "phrases") {
				page.Phrases = getPhrasesFromBlock(explanation)
				continue
			}

			if explanation.DataAtom.String() == "h3" && htmlquery.InnerText(explanation) == "Примеры" {
				nextIsDivBlock, _ := htmlquery.Query(explanation, "following-sibling::div[contains(@class,'block')]")
				if nextIsDivBlock != nil {
					page.Examples = getExamplesFromBlock(nextIsDivBlock)
				}
			}

			if explanation.DataAtom.String() == "h3" && htmlquery.InnerText(explanation) == "Фразовые глаголы" {
				nextIsDivBlock, _ := htmlquery.Query(explanation, "following-sibling::div[contains(@class,'block')]")
				if nextIsDivBlock != nil {
					page.PhraseVerb = getPhrasesWithLinkFromBlock(nextIsDivBlock)
				}
			}

			if explanation.DataAtom.String() == "h3" && htmlquery.InnerText(explanation) == "Возможные однокоренные слова" {
				nextIsDivBlock, _ := htmlquery.Query(explanation, "following-sibling::div[contains(@class,'block')]")
				if nextIsDivBlock != nil {
					page.PossibleCognates = getPhrasesWithLinkFromBlock(nextIsDivBlock)
				}
			}

			if explanation.DataAtom.String() == "div" && htmlquery.SelectAttr(explanation, "class") == "word_form_block" {
				page.Form = append(page.Form, getWordFormFromBlock(explanation))
			}
		}
	} else {
		if value, err := htmlquery.Query(html1, "//div[@id='wd_title']/p"); err == nil && value != nil {
			page.GeneralTranslate = strings.Split(cleanNodeField(value), ", ")
		}
		if value, err := htmlquery.Query(html1, "//div[@id='wd_content']/div[contains(@class,'word_ex')]"); err == nil && value != nil {
			page.Phrases = getPhrasesFromRuBlock(value)
		}

	}

	if !helper.IsEmpty(page.Title) && len(page.GeneralTranslate) > 0 {
		page.RequestText = query
	}
	return
}

func getMeaningBlock(node *html.Node) (meanings []Meaning) {
	currentMeaning := Meaning{}
	each, _ := htmlquery.QueryAll(node, "node()")
	var buffer string
	for _, explanation := range each {
		tag := explanation.DataAtom.String()
		if tag == "br" && helper.Len(buffer) > 0 {
			currentMeaning.Text = buffer
			meanings = append(meanings, currentMeaning)
			currentMeaning = Meaning{}
			buffer = ""
			continue
		}
		if tag == "span" {
			currentMeaning.Text = cleanNodeField(explanation)
			continue
		}
		if tag == "div" && strings.Contains(htmlquery.SelectAttr(explanation, "class"), "ex") {
			currentMeaning.Phrases = getPhrasesFromBlock(explanation)

			if !helper.IsEmpty(currentMeaning.Text) {
				meanings = append(meanings, currentMeaning)
			}
			currentMeaning = Meaning{}
			continue
		}
		if isHasMoreBLock(explanation) {
			meanings = append(meanings, getMeaningBlock(explanation)...)
			continue
		}
		text := htmlquery.InnerText(explanation)
		if string([]rune(text)[0:2]) == "- " {
			text = string([]rune(text)[2:])
		}
		if !helper.IsEmpty(text) {
			buffer += text
		}
	}
	return
}

func getPhrasesFromBlock(node *html.Node) (phrases []Phrase) {
	each, err := htmlquery.QueryAll(node, "node()")
	if err != nil {
		fmt.Println(err)
	}
	var buffer string
	for _, explanation := range each {
		tag := explanation.DataAtom.String()
		text := cleanNodeField(explanation)

		if helper.IsEmpty(text) {
			continue
		}

		if explanation.DataAtom.String() == "i" {
			phrases = append(
				phrases,
				Phrase{
					Text:      cleanTextField(buffer),
					Translate: text,
				},
			)
			buffer = ""
			continue
		}
		nextIsI, _ := htmlquery.Query(explanation, "following-sibling::i")

		if nextIsI != nil && string([]rune(text)[helper.Len(text)-1:]) == "—" {
			text = string([]rune(text)[:helper.Len(text)-1])
		}
		if tag == "a" {
			text = " " + text + " "
		}
		buffer += text
	}
	return
}

func getExamplesFromBlock(node *html.Node) (phrases []Phrase) {
	each, err := htmlquery.QueryAll(node, "node()")
	if err != nil {
		fmt.Println(err)
	}
	for _, explanation := range each {
		if htmlquery.SelectAttr(explanation, "id") == "hidden_ex" {
			phrases = append(phrases, getExamplesFromBlock(explanation)...)
		}
		if htmlquery.SelectAttr(explanation, "class") == "ex_o" {
			nextBlock, _ := htmlquery.Query(explanation, "following-sibling::p[contains(@class,'ex_t')]")
			if nextBlock != nil {
				phrases = append(phrases, Phrase{Text: cleanNodeField(explanation), Translate: cleanNodeField(nextBlock)})
			}
		}
	}
	return
}

func getPhrasesFromRuBlock(node *html.Node) (phrases []Phrase) {
	each, err := htmlquery.QueryAll(node, "node()")
	if err != nil {
		fmt.Println(err)
	}
	for _, explanation := range each {
		if explanation.DataAtom.String() == "span" {

			beforeBlock, _ := htmlquery.Query(explanation, "preceding-sibling::text()")
			if beforeBlock != nil {
				text := cleanNodeField(beforeBlock)
				if helper.IsEmpty(text) {
					continue
				}
				if string([]rune(text)[helper.Len(text)-1:]) == "—" {
					text = cleanTextField(string([]rune(text)[:helper.Len(text)-1]))
				}
				phrases = append(phrases, Phrase{Text: text, Translate: cleanNodeField(explanation)})
			}
		}

		if explanation.DataAtom.String() == "div" && htmlquery.SelectAttr(explanation, "id") == "ex_sup" {
			phrases = append(phrases, getPhrasesFromRuBlock(explanation)...)
		}
	}
	return
}

func getPhrasesWithLinkFromBlock(node *html.Node) (phrases []PhraseLink) {
	each, err := htmlquery.QueryAll(node, "node()")
	if err != nil {
		fmt.Println(err)
	}
	for _, explanation := range each {
		if explanation.DataAtom.String() == "a" {
			nextBlock, _ := htmlquery.Query(explanation, "following-sibling::text()")
			if nextBlock != nil {
				translate := cleanNodeField(nextBlock)

				if string([]rune(translate)[:1]) == "—" {
					translate = string([]rune(translate)[1:helper.Len(translate)])
				}

				phrases = append(
					phrases,
					PhraseLink{
						Text:      cleanNodeField(explanation),
						Link:      cleanTextField(htmlquery.SelectAttr(explanation, "href")),
						Translate: cleanTextField(translate),
					},
				)

			}
		}
	}
	return
}

func getWordFormFromBlock(node *html.Node) (form Forms) {
	each, err := htmlquery.QueryAll(node, "node()")
	if err != nil {
		fmt.Println(err)
	}

	var bufferInfo string
	var bufferValue string
	var bufferLink string
	for _, explanation := range each {
		if explanation.DataAtom.String() == "i" {
			form.Type = cleanNodeField(explanation)
			continue
		}
		if explanation.DataAtom.String() == "a" {
			bufferLink = cleanTextField(htmlquery.SelectAttr(explanation, "href"))
			bufferValue += cleanNodeField(explanation)
			continue
		}
		if explanation.DataAtom.String() == "span" {
			bufferInfo += cleanNodeField(explanation)
			continue
		}
		if explanation.DataAtom.String() == "br" {
			if !helper.IsEmpty(bufferValue) {
				if !helper.IsEmpty(bufferInfo) {
					if string([]rune(bufferInfo)[helper.Len(bufferInfo)-1:]) == ":" {
						bufferInfo = string([]rune(bufferInfo)[:helper.Len(bufferInfo)-1])
					}
				}
				form.Form = append(form.Form, WordForm{Info: bufferInfo, Link: bufferLink, Value: bufferValue})
			}
			bufferInfo = ""
			bufferValue = ""
			bufferLink = ""
			continue
		}

		bufferValue += cleanNodeField(explanation)

	}
	return
}

func isHasMoreBLock(node *html.Node) bool {
	if node.DataAtom.String() != "div" {
		return false
	}
	id := htmlquery.SelectAttr(node, "id")
	if id == idMoreVerb || id == idMoreNoun || id == idMoreAdjective || id == idMoreAdverb || id == idMorePreposition {
		return true
	}

	return false
}

func cleanTextField(text string) string {
	return strings.TrimSpace(strings.ToLower(text))
}
func cleanNodeField(node *html.Node) string {
	return cleanTextField(htmlquery.InnerText(node))
}
