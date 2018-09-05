package mini_message

type Run struct {
	Data []RunRecord `json:"stepInfoList"`
}

type RunRecord struct {
	Timestamp int64 `json:"timestamp"`
	Steps     uint  `json:"step"`
}
