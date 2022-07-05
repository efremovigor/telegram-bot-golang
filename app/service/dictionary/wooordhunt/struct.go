package wooordhunt

type Page struct {
	RequestText      string            `json:"request_text"`
	Title            string            `json:"title"`
	Info             []Info            `json:"info"`
	Phrases          []Phrase          `json:"phrases"`
	Examples         []Phrase          `json:"examples"`
	PhraseVerb       []PhraseLink      `json:"phrase_verb"`
	PossibleCognates []PhraseLink      `json:"possible_cognates"`
	Form             []Forms           `json:"form"`
	VoicePath        VoicePath         `json:"voice_path"`
	Transcription    map[string]string `json:"transcription"`
	GeneralTranslate []string          `json:"generalTranslate"`
}

type Phrase struct {
	Text      string `json:"text"`
	Translate string `json:"translate"`
}

type PhraseLink struct {
	Text      string `json:"text"`
	Link      string `json:"link"`
	Translate string `json:"translate"`
}

type Forms struct {
	Type string     `json:"type"`
	Form []WordForm `json:"form"`
}

type WordForm struct {
	Info  string `json:"text"`
	Link  string `json:"link"`
	Value string `json:"translate"`
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
