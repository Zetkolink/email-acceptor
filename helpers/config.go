package helpers

import (
	"github.com/spf13/viper"
	"log"
)

func InitConfig() Config {
	cfg := Config{}
	cfg.loadConfig()
	return cfg
}

type Config struct {
	Db        DbConfig     `yaml:"db"`
	Smtp      SmtpConfig   `yaml:"smtp"`
	Rb        RabbitConfig `yaml:"rb"`
	LogLevel  string       `yaml:"logLevel"`
	LogFormat string       `yaml:"logFormat"`
	Addr      string       `yaml:"addr"`
}

type DbConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

type SmtpConfig struct {
	Identity string `yaml:"identity"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
}

type RabbitConfig struct {
	Amqp     string `yaml:"amqp"`
	Queue    string `yaml:"queue"`
	Exchange string `yaml:"exchange"`
}

func (c *Config) loadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}
}
