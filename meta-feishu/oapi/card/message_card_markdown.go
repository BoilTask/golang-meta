package card

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type MessageCardMarkdown struct {
	TextAlign_ MessageCardTextAlign                `json:"text_align,omitempty"`
	Content_   string                              `json:"content,omitempty"`
	Href_      map[string]*larkcard.MessageCardURL `json:"href,omitempty"`
}

func NewMessageCardMarkdown() *MessageCardMarkdown {
	return &MessageCardMarkdown{}
}

func (markDown *MessageCardMarkdown) TextAlign(textAlign_ MessageCardTextAlign) *MessageCardMarkdown {
	markDown.TextAlign_ = textAlign_
	return markDown
}

func (markDown *MessageCardMarkdown) Content(content string) *MessageCardMarkdown {
	markDown.Content_ = content
	return markDown
}

func (markDown *MessageCardMarkdown) Href(href map[string]*larkcard.MessageCardURL) *MessageCardMarkdown {
	markDown.Href_ = href
	return markDown
}

func (markDown *MessageCardMarkdown) Build() *MessageCardMarkdown {
	return markDown
}

func (m *MessageCardMarkdown) Tag() string {
	return "markdown"
}

func (m *MessageCardMarkdown) MarshalJSON() ([]byte, error) {
	return MessageCardElementJson(m)
}
