package entity

type Config struct {
	Redis RedisConfig `yaml:"redis"`
}

type RedisConfig struct {
	LogBase    RedisDBConfig `yaml:"logBase"`
	ChangeBase RedisDBConfig `yaml:"changeBase"`
}

type RedisDBConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
