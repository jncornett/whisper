package whisper

type Message struct {
	Author  string `json:"author"`
	Token   string `json:"token"`
	Content string `json:"content"`
}
