package mini_message

type Run struct {
	Data []RunRecord `json:"stepInfoList"`
}

type RunRecord struct {
	Timestamp uint `json:"timestamp"`
	Steps     uint `json:"step"`
}
