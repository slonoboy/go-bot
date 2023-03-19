package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		Bot      Bot
		HTTP     HTTP
		Database Database
	}
	Bot struct {
		Token    string
		Username string
	}

	HTTP struct {
		Domain string
		Port   string
	}
	Database struct {
		Host     string
		User     string
		Password string
		Name     string
		Port     string
	}
)

func Init() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	setFromEnv(&cfg)

	return &cfg, nil
}

func setFromEnv(cfg *Config) {
	cfg.Bot.Token = viper.GetString("BOT_TOKEN")
	cfg.Bot.Username = viper.GetString("BOT_USERNAME")

	cfg.HTTP.Domain = viper.GetString("HTTP_DOMAIN")
	cfg.HTTP.Port = viper.GetString("HTTP_PORT")

	cfg.Database.Host = viper.GetString("POSTGRES_HOST")
	cfg.Database.Port = viper.GetString("POSTGRES_PORT")
	cfg.Database.Name = viper.GetString("POSTGRES_DB")
	cfg.Database.User = viper.GetString("POSTGRES_USER")
	cfg.Database.Password = viper.GetString("POSTGRES_PASSWORD")
}
