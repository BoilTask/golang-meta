package metafeishu

type AppConfig struct {
	AppId             string            `yaml:"app-id"`             // 应用ID
	AppSecret         string            `yaml:"app-secret"`         // 应用密钥
	VerificationToken string            `yaml:"verification-token"` // 验证Token
	EventEncryptKey   string            `yaml:"encrypt-key"`        // 事件加密Key
	OpenIds           map[string]string `yaml:"open-ids"`           // 自己视角内的OpenId
}

func (c *AppConfig) GetOpenId(key string) string {
	return c.OpenIds[key]
}
