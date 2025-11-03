package card

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type MessageCardColumn struct {
	Width_         MessageCardWidth              `json:"width,omitempty"`
	Weight_        int                           `json:"weight,omitempty"`
	VerticalAlign_ MessageCardVerticalAlign      `json:"vertical_align,omitempty"`
	Elements_      []larkcard.MessageCardElement `json:"elements,omitempty"`
}

func NewMessageCardColumn() *MessageCardColumn {
	return &MessageCardColumn{}
}

func (c *MessageCardColumn) Width(width MessageCardWidth) *MessageCardColumn {
	c.Width_ = width
	return c
}

func (c *MessageCardColumn) Weight(weight int) *MessageCardColumn {
	c.Weight_ = weight
	return c
}

func (c *MessageCardColumn) VerticalAlign(verticalAlign MessageCardVerticalAlign) *MessageCardColumn {
	c.VerticalAlign_ = verticalAlign
	return c
}

func (c *MessageCardColumn) Elements(elements []larkcard.MessageCardElement) *MessageCardColumn {
	c.Elements_ = elements
	return c
}

func (c *MessageCardColumn) Build() *MessageCardColumn {
	return c
}

func (c *MessageCardColumn) Tag() string {
	return "column"
}

func (c *MessageCardColumn) MarshalJSON() ([]byte, error) {
	return MessageCardElementJson(c)
}
