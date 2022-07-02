package wooordhunt

type Page struct {
	RequestText   string            `json:"request_text"`
	Options       []Info            `json:"options"`
	VoicePath     VoicePath         `json:"voice_path"`
	Transcription map[string]string `json:"transcription"`
}

type Info struct {
	Text          string `json:"text"`
	Type          string `json:"type"`
	Transcription string `json:"transcription"`
}

type VoicePath struct {
	UK string `json:"uk"`
	US string `json:"us"`
}
