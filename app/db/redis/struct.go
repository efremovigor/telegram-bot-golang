package redis

type VoiceFile struct {
	Lang string `json:"lang"`
	Url  string `json:"url"`
	Word string `json:"word"`
}

type PicFile struct {
	Url  string `json:"url"`
	Word string `json:"word"`
}
