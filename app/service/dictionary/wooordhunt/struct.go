package wooordhunt

type Page struct {
	RequestText      string            `json:"request_text"`
	Info             []Info            `json:"info"`
	Phrases          []Phrase          `json:"phrases"`
	Examples         []Phrase          `json:"examples"`
	PhraseVerb       []PhraseVerb      `json:"phrase_verb"`
	VoicePath        VoicePath         `json:"voice_path"`
	Transcription    map[string]string `json:"transcription"`
	GeneralTranslate []string          `json:"generalTranslate"`
}

type Phrase struct {
	Text      string `json:"text"`
	Translate string `json:"translate"`
}

type PhraseVerb struct {
	Text      string `json:"text"`
	Link      string `json:"link"`
	Translate string `json:"translate"`
}

type Info struct {
	Meaning []Meaning `json:"meaning"`
	Type    string    `json:"type"`
}

type Meaning struct {
	Text    string   `json:"text"`
	Phrases []Phrase `json:"phrases"`
}

type VoicePath struct {
	UK string `json:"uk"`
	US string `json:"us"`
}

const idMoreVerb = "pos_verb"
const idMoreNoun = "pos_noun"
const idMoreAdjective = "pos_adjective"
const idMoreAdverb = "pos_adverb"
const idMorePreposition = "pos_preposition"
