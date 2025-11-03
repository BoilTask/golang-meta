package card

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type MessageCardBehavior interface {
	Type() string
}

type MessageCardBehaviorOpenUrl struct {
	DefaultUrl_ string `json:"default_url"`
	AndroidUrl_ string `json:"android_url"`
	IosUrl_     string `json:"ios_url"`
	PcUrl_      string `json:"pc_url"`
}

func NewMessageCardBehaviorOpenUrl() *MessageCardBehaviorOpenUrl {
	return &MessageCardBehaviorOpenUrl{}
}

func (m *MessageCardBehaviorOpenUrl) DefaultUrl(defaultUrl string) *MessageCardBehaviorOpenUrl {
	m.DefaultUrl_ = defaultUrl
	return m
}

func (m *MessageCardBehaviorOpenUrl) AndroidUrl(androidUrl string) *MessageCardBehaviorOpenUrl {
	m.AndroidUrl_ = androidUrl
	return m
}

func (m *MessageCardBehaviorOpenUrl) IosUrl(iosUrl string) *MessageCardBehaviorOpenUrl {
	m.IosUrl_ = iosUrl
	return m
}

func (m *MessageCardBehaviorOpenUrl) PcUrl(pcUrl string) *MessageCardBehaviorOpenUrl {
	m.PcUrl_ = pcUrl
	return m
}

func (m *MessageCardBehaviorOpenUrl) Type() string {
	return "open_url"
}

func (m *MessageCardBehaviorOpenUrl) MarshalJSON() ([]byte, error) {
	return MessageCardBehaviorJson(m)
}

func (m *MessageCardBehaviorOpenUrl) Build() *MessageCardBehaviorOpenUrl {
	return m
}

type MessageCardBehaviorCallback struct {
	Value_ map[string]interface{} `json:"value,omitempty"`
}

func NewMessageCardBehaviorCallback() *MessageCardBehaviorCallback {
	return &MessageCardBehaviorCallback{}
}

func (m *MessageCardBehaviorCallback) Value(value map[string]interface{}) *MessageCardBehaviorCallback {
	m.Value_ = value
	return m
}

func (m *MessageCardBehaviorCallback) Type() string {
	return "callback"
}

func (m *MessageCardBehaviorCallback) MarshalJSON() ([]byte, error) {
	return MessageCardBehaviorJson(m)
}

func (m *MessageCardBehaviorCallback) Build() *MessageCardBehaviorCallback {
	return m
}

type MessageCardButton struct {
	Type_      *larkcard.MessageCardButtonType    `json:"type,omitempty"`
	Text_      larkcard.MessageCardText           `json:"text,omitempty"`
	Confirm_   *larkcard.MessageCardActionConfirm `json:"confirm,omitempty"`
	Behaviors_ []MessageCardBehavior              `json:"behaviors,omitempty"`
}

func NewMessageCardButton() *MessageCardButton {
	return &MessageCardButton{}
}

func (m *MessageCardButton) Type(type_ larkcard.MessageCardButtonType) *MessageCardButton {
	m.Type_ = &type_
	return m
}

func (m *MessageCardButton) Text(text larkcard.MessageCardText) *MessageCardButton {
	m.Text_ = text
	return m
}

func (m *MessageCardButton) Confirm(confirm *larkcard.MessageCardActionConfirm) *MessageCardButton {
	m.Confirm_ = confirm
	return m
}

func (m *MessageCardButton) Behaviors(behaviors []MessageCardBehavior) *MessageCardButton {
	m.Behaviors_ = behaviors
	return m
}

func (m *MessageCardButton) Build() *MessageCardButton {
	return m
}

func (m *MessageCardButton) Tag() string {
	return "button"
}

func (m *MessageCardButton) MarshalJSON() ([]byte, error) {
	return MessageCardElementJson(m)
}
