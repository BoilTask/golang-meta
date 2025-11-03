package metamongo

type Config struct {
	Uri      string `yaml:"uri"`      // 地址
	Username string `yaml:"username"` // 用户名
	Password string `yaml:"password"` // 密码
}
