package rapid_microsoft

type MicrosoftTranslate struct {
	Translations []Translation `json:"translations"`
}
type Translation struct {
	Text string `json:"text"`
	To   string `json:"to"`
}
