package message

type MessagePostTag struct {
	Tag      string  `json:"tag"`
	ImageKey *string `json:"image_key"`
	Text     *string `json:"text"`
}

type MessagePost struct {
	Title   *string             `json:"title"`
	Content [][]*MessagePostTag `json:"content"`
}
