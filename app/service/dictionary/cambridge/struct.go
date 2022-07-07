package cambridge

import "telegram-bot-golang/helper"

const Url = "https://dictionary.cambridge.org"
const SearchUrl = Url + "/autocomplete/amp?dataset=english&q=%s&__amp_source_origin=" + Url

type Page struct {
	RequestText string `json:"request_text"`
	Options     []Info `json:"options"`
}

type Info struct {
	Text          string            `json:"text"`
	Type          string            `json:"type"`
	Forms         []Forms           `json:"forms"`
	Transcription map[string]string `json:"transcription"`
	Explanation   []Explanation     `json:"explanation"`
	VoicePath     VoicePath         `json:"voice_path"`
	Image         string            `json:"image"`
}

type Forms struct {
	Desc  string `json:"desc"`
	Value string `json:"value"`
}

type VoicePath struct {
	UK string `json:"uk"`
	US string `json:"us"`
}

type Explanation struct {
	Text                string   `json:"text"`
	SemanticDescription string   `json:"semantic_description"`
	Level               string   `json:"level"`
	Description         string   `json:"description"`
	Translate           string   `json:"translate"`
	Example             []string `json:"example"`
}

func (i Page) IsValid() bool {
	return !helper.IsEmpty(i.RequestText)
}

type SearchResponse struct {
	RequestWord string          `json:"request_word"`
	Founded     []SearchElement `json:"founded"`
}

type SearchElement struct {
	Word string `json:"word"`
	Path string `json:"url"`
}

func (s SearchResponse) IsValid() bool {
	return !helper.IsEmpty(s.RequestWord)
}

const xpathBlockDescriptionEnRu = "//article[@id='page-content']//div[contains(@class, 'entry-body__el')]"
const xpathAltBlockDescriptionEnRu = "//article[@id='page-content']//div[contains(@class, 'di-body')]"
const xpathTitle = "//div[contains(@class, 'di-title')]"
const xpathType = "//div[contains(@class, 'posgram')]"
const xpathComplexType = "//span[contains(@class, 'di-info')]/div"
const xpathForms = "//span[contains(@class, 'irreg-infls')]/span[contains(@class, 'inf-group')]"
const xpathImage = "//div[contains(@class, 'dimg')]//amp-img[contains(@class, 'dimg_i')]"
const xpathTranscriptionUK = "//span[contains(@class, 'uk')]/span[contains(@class, 'pron')]"
const xpathTranscriptionUS = "//span[contains(@class, 'us')]/span[contains(@class, 'pron')]"
const xpathUK = "//span[contains(@class, 'uk')]//amp-audio//source[contains(@type,'audio/mpeg')]"
const xpathUS = "//span[contains(@class, 'us')]//amp-audio//source[contains(@type,'audio/mpeg')]"
const xpathExplanations = "//div[contains(concat(\" \", normalize-space(@class), \" \"), \" dsense \")]"
const xpathExplanationsSemanticDescription = "//h3[contains(@class, 'dsense_h')]"
const xpathExplanationsWord = "//span[contains(@class, 'phrase-title')]"
const xpathExplanationsLevel = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'ddef_h')]/span/span"
const xpathExplanationsDescription = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'ddef_h')]/div"
const xpathExplanationsTranslate = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'def-body')]/span"
const xpathExplanationsExamples = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'dexamp')]"
const xpathExplanationsMoreExamples = "//div[contains(@class, 'daccord')]//ul/li[contains(@class, 'dexamp')]"
