package card

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type MessageCardFormButton struct {
	ActionType_ string                             `json:"action_type,omitempty"`
	Name_       string                             `json:"name,omitempty"`
	Text_       larkcard.MessageCardText           `json:"text,omitempty"`
	Url_        string                             `json:"url,omitempty"`
	MultiUrl_   *larkcard.MessageCardURL           `json:"multi_url,omitempty"`
	ButtonType_ string                             `json:"button_type,omitempty"`
	Value_      map[string]interface{}             `json:"value,omitempty"`
	Confirm_    *larkcard.MessageCardActionConfirm `json:"confirm,omitempty"`
}

func NewMessageCardFormButton() *MessageCardFormButton {
	return &MessageCardFormButton{}
}

func (c *MessageCardFormButton) Name(name string) *MessageCardFormButton {
	c.Name_ = name
	return c
}

func (c *MessageCardFormButton) Text(text larkcard.MessageCardText) *MessageCardFormButton {
	c.Text_ = text
	return c
}

func (c *MessageCardFormButton) Url(url string) *MessageCardFormButton {
	c.Url_ = url
	return c
}

func (c *MessageCardFormButton) MultiUrl(multiUrl *larkcard.MessageCardURL) *MessageCardFormButton {
	c.MultiUrl_ = multiUrl
	return c
}

func (c *MessageCardFormButton) ButtonType(buttonType larkcard.MessageCardButtonType) *MessageCardFormButton {
	c.ButtonType_ = string(buttonType)
	return c
}

func (c *MessageCardFormButton) Value(value map[string]interface{}) *MessageCardFormButton {
	c.Value_ = value
	return c
}

func (c *MessageCardFormButton) Confirm(confirm *larkcard.MessageCardActionConfirm) *MessageCardFormButton {
	c.Confirm_ = confirm
	return c
}

func (c *MessageCardFormButton) ActionType(actionType string) *MessageCardFormButton {
	c.ActionType_ = actionType
	return c
}

func (c *MessageCardFormButton) Build() *MessageCardFormButton {
	return c
}

func (c *MessageCardFormButton) Tag() string {
	return "button"
}

func (c *MessageCardFormButton) MarshalJSON() ([]byte, error) {
	return MessageCardElementJson(c)
}
