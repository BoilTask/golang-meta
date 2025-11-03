package metaredis

type Config struct {
	Addr     string `yaml:"addr"`     // 地址
	Password string `yaml:"password"` // 密码
}
