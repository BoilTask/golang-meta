package card

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type MessageCardFormSelectStatic struct {
	*larkcard.MessageCardEmbedSelectMenuStatic
	Name_ string `json:"name,omitempty"`
}

func NewMessageCardFormSelectStatic() *MessageCardFormSelectStatic {
	return &MessageCardFormSelectStatic{}
}

func (c *MessageCardFormSelectStatic) MessageCardFormSelectStatic(static *larkcard.MessageCardEmbedSelectMenuStatic) *MessageCardFormSelectStatic {
	c.MessageCardEmbedSelectMenuStatic = static
	return c
}

func (c *MessageCardFormSelectStatic) Name(name string) *MessageCardFormSelectStatic {
	c.Name_ = name
	return c
}

func (c *MessageCardFormSelectStatic) Build() *MessageCardFormSelectStatic {
	return c
}

func (c *MessageCardFormSelectStatic) Tag() string {
	return "select_static"
}

func (c *MessageCardFormSelectStatic) MarshalJSON() ([]byte, error) {
	return MessageCardElementJson(c)
}
