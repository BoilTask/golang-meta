package bitable

type FiledLinkRecord struct {
	RecordIds []string `json:"record_ids"`
	TableId   string   `json:"table_id"`
	Text      string   `json:"text"`
	TextArr   []string `json:"text_arr"`
	Type      string   `json:"type"`
}
