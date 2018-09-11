package mini_message

type Run struct {
	Data []RunRecord `json:"data"`
}

type RunRecord struct {
	Timestamp int64 `json:"timestamp"`
	Steps     int   `json:"step"`
}
