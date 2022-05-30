package cambridge

type Info struct {
	Text        string        `json:"text"`
	Type        string        `json:"type"`
	Explanation []Explanation `json:"explanation"`
}

type Explanation struct {
	SemanticDescription string   `json:"semantic_description"`
	Level               string   `json:"level"`
	Description         string   `json:"description"`
	Translate           string   `json:"translate"`
	Example             []string `json:"example"`
}

const xpathBLockDescriptionEnRu = "//article[@id=\"page-content\"]//div[contains(@class, 'entry-body')]//div[contains(@class, 'entry-body__el')]"
const xpathTitle = xpathBLockDescriptionEnRu + "//div[contains(@class, 'di-title')]/span/span"
const xpathType = xpathBLockDescriptionEnRu + "//div[contains(@class, 'posgram')]/span"
const xpathExplanations = xpathBLockDescriptionEnRu + "//div[contains(@class, 'pos-body')]/div[contains(@class, 'dsense')]"
const xpathExplanationsSemanticDescription = "//h3[contains(@class, 'dsense_h')]"
const xpathExplanationsLevel = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'ddef_h')]/span/span"
const xpathExplanationsDescription = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'ddef_h')]/div"
const xpathExplanationsTranslate = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'def-body')]/span"
const xpathExplanationsExamples = "//div[contains(@class, 'sense-body')]//div[contains(@class, 'examp')]"
