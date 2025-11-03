package card

type MessageCardColumnSet struct {
	FlexMode_ MessageCardFlexMode  `json:"flex_mode,omitempty"`
	Columns_  []*MessageCardColumn `json:"columns,omitempty"`
}

func NewMessageCardColumnSet() *MessageCardColumnSet {
	return &MessageCardColumnSet{}
}

func (c *MessageCardColumnSet) FlexMode(flexMode MessageCardFlexMode) *MessageCardColumnSet {
	c.FlexMode_ = flexMode
	return c
}

func (c *MessageCardColumnSet) Columns(columns []*MessageCardColumn) *MessageCardColumnSet {
	c.Columns_ = columns
	return c
}

func (c *MessageCardColumnSet) Build() *MessageCardColumnSet {
	return c
}

func (c *MessageCardColumnSet) Tag() string {
	return "column_set"
}

func (c *MessageCardColumnSet) MarshalJSON() ([]byte, error) {
	return MessageCardElementJson(c)
}
