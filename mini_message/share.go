package mini_message

type Share struct {
	OpenGId   string    `json:"openGId"`
	Watermark Watermark `json:"watermark"`
}
