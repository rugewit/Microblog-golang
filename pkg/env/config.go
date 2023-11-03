package env

import "github.com/spf13/viper"

type Config struct {
	DbName      string `mapstructure:"DB_NAME"`
	DbUrl       string `mapstructure:"DB_URL"`
	UserColName string `mapstructure:"USER_COL_NAME"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./pkg/common/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}
