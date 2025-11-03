package metaemail

type Config struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}
