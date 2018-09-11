package mini_message

type Share struct {
	OpenGId   string    `json:"open_gid"`
	Watermark Watermark `json:"watermark"`
}
