package config

import (
	"github.com/spf13/viper"
	"infac/internal/models"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	SUNAT       SUNATConfig       `mapstructure:"sunat"`
	Certificate CertificateConfig `mapstructure:"certificate"`
	Issuer      models.Company    `mapstructure:"issuer"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type SUNATConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	OSE      OSEConfig `mapstructure:"ose"`
}

type OSEConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Provider string `mapstructure:"provider"`
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type CertificateConfig struct {
	PFXPath  string `mapstructure:"pfx_path"`
	Password string `mapstructure:"password"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	
	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("sunat.url", "https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService")
	viper.SetDefault("sunat.ose.enabled", false)
	
	// Environment variables
	viper.SetEnvPrefix("INFAC")
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	
	return &config, nil
}