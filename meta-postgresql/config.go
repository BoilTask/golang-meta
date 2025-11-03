package metapostgresql

type Config struct {
	Host     string `yaml:"host"`     // 地址
	Port     int    `yaml:"port"`     // 端口
	Username string `yaml:"username"` // 用户名
	Password string `yaml:"password"` // 密码
	Database string `yaml:"database"` // 数据库
}
