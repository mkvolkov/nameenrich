package cfg

import "github.com/spf13/viper"

type Cfg struct {
	Port        string `mapstructure:"PORT"`
	PostgresURL string `mapstructure:"POSTGRES_URL"`
	AtrHost     string `mapstructure:"ATREUGO_HOST"`
	AtrPort     string `mapstructure:"ATREUGO_PORT"`
}

func LoadConfig(cfg *Cfg) error {
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	return nil
}
