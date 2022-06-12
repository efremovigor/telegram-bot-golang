package multitran

import "strings"

type MultitranInfo struct {
	RequestText string `json:"request_text"`
	Options     []Info `json:"options"`
}

type Info struct {
	Text          string        `json:"text"`
	Type          string        `json:"type"`
	Transcription string        `json:"transcription"`
	Explanation   []Explanation `json:"explanation"`
}

type Explanation struct {
	Type string   `json:"type"`
	Text []string `json:"text"`
}

func generalTypes() []string {
	return []string{
		"general", "agriculture", "dialectal", "programming", "religion", "slang", "vulgar", "american",
		"biology", "diplomacy", "drilling", "electronics", "figurative", "informal", "folklore", "poetic",
		"figurative", "dated", "information technology", "mathematics", "mining", "technology", "telecommunications",
		"military", "literature", "cooking", "automated equipment", "astronautics",
	}
}

func isGeneralType(text string) bool {
	for _, generalType := range generalTypes() {
		if strings.Contains(text, generalType) {
			return true
		}
	}
	return false
}
