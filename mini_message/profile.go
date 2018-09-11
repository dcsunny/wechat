package mini_message

type Profile struct {
	OpenID    string    `json:"openid"`
	UnionID   string    `json:"unionid"`
	NickName  string    `json:"nickname"`
	Gender    int       `json:"gender"`
	City      string    `json:"city"`
	Province  string    `json:"province"`
	Country   string    `json:"country"`
	AvatarURL string    `json:"avatar_url"`
	Language  string    `json:"language"`
	Watermark Watermark `json:"watermark"`
}
