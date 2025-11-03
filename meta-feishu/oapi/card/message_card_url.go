package card

type MessageCardURL struct {
	URL_        string  `json:"url,omitempty"`
	AndroidURL_ *string `json:"android_url,omitempty"`
	IosUrl_     *string `json:"ios_url,omitempty"`
	PcUrl_      *string `json:"pc_url,omitempty"`
}

func NewMessageCardURL() *MessageCardURL {
	return &MessageCardURL{}
}

func (cardUrl *MessageCardURL) Url(url string) *MessageCardURL {
	cardUrl.URL_ = url
	return cardUrl
}

func (cardUrl *MessageCardURL) AndroidUrl(androidUrl string) *MessageCardURL {
	cardUrl.AndroidURL_ = &androidUrl
	return cardUrl
}

func (cardUrl *MessageCardURL) IoSUrl(iosUrl string) *MessageCardURL {
	cardUrl.IosUrl_ = &iosUrl
	return cardUrl
}

func (cardUrl *MessageCardURL) PcUrl(pcURL string) *MessageCardURL {
	cardUrl.PcUrl_ = &pcURL
	return cardUrl
}

func (cardUrl *MessageCardURL) Build() *MessageCardURL {
	if cardUrl.PcUrl_ == nil {
		cardUrl.PcUrl_ = &cardUrl.URL_
	}
	if cardUrl.AndroidURL_ == nil {
		cardUrl.AndroidURL_ = &cardUrl.URL_
	}
	if cardUrl.IosUrl_ == nil {
		cardUrl.IosUrl_ = &cardUrl.URL_
	}
	return cardUrl
}
