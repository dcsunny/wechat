package mini_message

type Profile struct {
	OpenID    string    `json:"openId"`
	UnionID   string    `json:"unionId"`
	NickName  string    `json:"nickName"`
	Gender    int       `json:"gender"`
	City      string    `json:"city"`
	Province  string    `json:"province"`
	Country   string    `json:"country"`
	AvatarURL string    `json:"avatarUrl"`
	Language  string    `json:"language"`
	Watermark Watermark `json:"watermark"`
}
type Run struct {
	Data []RunRecord `json:"stepInfoList"`
}

type RunRecord struct {
	Timestamp int64 `json:"timestamp"`
	Steps     int   `json:"step"`
}

type Share struct {
	OpenGId   string    `json:"openGId"`
	Watermark Watermark `json:"watermark"`
}

type Watermark struct {
	Timestamp int64  `json:"timestamp"`
	AppID     string `json:"appid"`
}

// PhoneInfo 用户手机号
type PhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}
