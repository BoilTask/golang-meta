package card

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type MessageCardFormInput struct {
	Name_         string                         `json:"name,omitempty"`
	PlaceHolder_  *larkcard.MessageCardPlainText `json:"placeholder,omitempty"`
	Label_        *larkcard.MessageCardPlainText `json:"label,omitempty"`
	MaxLength_    int                            `json:"max_length,omitempty"`
	DefaultValue_ string                         `json:"default_value,omitempty"`
	Require_      bool                           `json:"require,omitempty"`
	Disabled_     bool                           `json:"disabled,omitempty"`
}

func NewMessageCardFormInput() *MessageCardFormInput {
	return &MessageCardFormInput{}
}

func (c *MessageCardFormInput) Name(name string) *MessageCardFormInput {
	c.Name_ = name
	return c
}

func (c *MessageCardFormInput) PlaceHolder(placeHolder *larkcard.MessageCardPlainText) *MessageCardFormInput {
	c.PlaceHolder_ = placeHolder
	return c
}

func (c *MessageCardFormInput) Label(label *larkcard.MessageCardPlainText) *MessageCardFormInput {
	c.Label_ = label
	return c
}

func (c *MessageCardFormInput) MaxLength(maxLength int) *MessageCardFormInput {
	c.MaxLength_ = maxLength
	return c
}

func (c *MessageCardFormInput) DefaultValue(defaultValue string) *MessageCardFormInput {
	c.DefaultValue_ = defaultValue
	return c
}

func (c *MessageCardFormInput) Require(require bool) *MessageCardFormInput {
	c.Require_ = require
	return c
}

func (c *MessageCardFormInput) Disabled(disabled bool) *MessageCardFormInput {
	c.Disabled_ = disabled
	return c
}

func (c *MessageCardFormInput) Build() *MessageCardFormInput {
	return c
}

func (c *MessageCardFormInput) Tag() string {
	return "input"
}

func (c *MessageCardFormInput) MarshalJSON() ([]byte, error) {
	return MessageCardElementJson(c)
}
