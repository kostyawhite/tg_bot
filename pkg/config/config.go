package config

import (
	"github.com/spf13/viper"
)

type Responses struct {
	ReplyStartTemplate     string `mapstructure:"reply_start_template"`
	ReplyAlreadyAuthorized string `mapstructure:"reply_already_authorized"`
	LinkSavedSuccessfully  string `mapstructure:"link_saved_successfully"`
}

type Errors struct {
	InvalidLink    string `mapstructure:"invalid_link"`
	FailedToSave   string `mapstructure:"failed_to_save"`
	Unauthorized   string `mapstructure:"unauthorized"`
	UnknownCommand string `mapstructure:"unknown_command"`
	UnknownError   string `mapstructure:"unknown_error"`
}
type Messages struct {
	Responses
	Errors
}

type Config struct {
	Token         string
	ConsumerKey   string
	AuthServerUrl string
	DbPath        string `mapstructure:"db_path"`
	TgBotUrl      string `mapstructure:"tg_bot_url"`

	Messages Messages
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	err := unmarshal(&config)
	if err != nil {
		return nil, err
	}

	if err := config.parseEnvs(); err != nil {
		return nil, err
	}

	return &config, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return err
	}
	return nil
}

func (c *Config) parseEnvs() error {
	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	if err := viper.BindEnv("consumer_key"); err != nil {
		return err
	}
	if err := viper.BindEnv("auth_server_url"); err != nil {
		return err
	}

	c.Token = viper.GetString("token")
	c.ConsumerKey = viper.GetString("consumer_key")
	c.AuthServerUrl = viper.GetString("auth_server_url")

	return nil
}
