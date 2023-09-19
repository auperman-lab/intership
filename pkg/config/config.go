package config

import "github.com/spf13/viper"

type Config struct {
	Port        string `mapstructure:"PORT"`
	AuthSvcUrl  string `mapstructure:"AUTH_SVC_URL"`
	UserSvcUrl  string `mapstructure:"USER_SVC_URL"`
	AdminSvcUrl string `mapstructure:"ADMIN_SVC_URL"`
	GRPCPort    string `mapstructure:"GRPCPORT"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
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
