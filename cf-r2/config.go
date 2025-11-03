package cfr2

type Config struct {
	Url    string `yaml:"url"`    // R2 数据服务地址
	Token  string `yaml:"token"`  // R2 数据服务 token
	Key    string `yaml:"key"`    // R2 数据服务密钥
	Secret string `yaml:"secret"` // R2 数据服务密钥
}
