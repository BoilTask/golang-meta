package card

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type MessageCardForm struct {
	Name_     string                        `json:"name,omitempty"`
	Elements_ []larkcard.MessageCardElement `json:"elements,omitempty"`
}

func NewMessageCardForm() *MessageCardForm {
	return &MessageCardForm{}
}

func (c *MessageCardForm) Name(name string) *MessageCardForm {
	c.Name_ = name
	return c
}

func (c *MessageCardForm) Elements(elements []larkcard.MessageCardElement) *MessageCardForm {
	c.Elements_ = elements
	return c
}

func (c *MessageCardForm) Build() *MessageCardForm {
	return c
}

func (c *MessageCardForm) Tag() string {
	return "form"
}

func (c *MessageCardForm) MarshalJSON() ([]byte, error) {
	return MessageCardElementJson(c)
}
