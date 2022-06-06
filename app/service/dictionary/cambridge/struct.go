package cambridge

import "telegram-bot-golang/helper"

const Url = "https://dictionary.cambridge.org"

type Info struct {
	Text          string        `json:"text"`
	Type          string        `json:"type"`
	Transcription string        `json:"transcription"`
	VoicePath     VoicePath     `json:"voice_path"`
	Explanation   []Explanation `json:"explanation"`
}

type VoicePath struct {
	UK string `json:"uk"`
	US string `json:"us"`
}

type Explanation struct {
	SemanticDescription string   `json:"semantic_description"`
	Level               string   `json:"level"`
	Description         string   `json:"description"`
	Translate           string   `json:"translate"`
	Example             []string `json:"example"`
}

func (i Info) IsValid() bool {
	return !helper.IsEmpty(i.Text)
}

const xpathBLockDescriptionEnRu = "//article[@id='page-content']//div[contains(@class, 'entry-body')]//div[contains(@class, 'entry-body__el')]"
const xpathTitle = xpathBLockDescriptionEnRu + "//div[contains(@class, 'di-title')]/span/span"
const xpathType = xpathBLockDescriptionEnRu + "//div[contains(@class, 'posgram')]/span"
const xpathTranscription = xpathBLockDescriptionEnRu + "//span/span[contains(@class, 'pron')]"
const xpathUK = xpathBLockDescriptionEnRu + "//span[contains(@class, 'uk')]//amp-audio//source[contains(@type,'audio/mpeg')]"
const xpathUS = xpathBLockDescriptionEnRu + "//span[contains(@class, 'us')]//amp-audio//source[contains(@type,'audio/mpeg')]"
const xpathExplanations = xpathBLockDescriptionEnRu + "//div[contains(@class, 'pos-body')]/div[contains(@class, 'dsense')]"
const xpathExplanationsSemanticDescription = "//h3[contains(@class, 'dsense_h')]"
const xpathExplanationsLevel = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'ddef_h')]/span/span"
const xpathExplanationsDescription = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'ddef_h')]/div"
const xpathExplanationsTranslate = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'def-body')]/span"
const xpathExplanationsExamples = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'examp')]"
