package config

import "github.com/spf13/viper"

type Config struct {
	BlockingTimeIp             int64  `mapstructure:"blocking_time_ip"`
	BlockingTimeApiKey         int64  `mapstructure:"blocking_time_api_key"`
	RequestPerSecondPerIp      int64  `mapstructure:"request_per_second_per_ip"`
	RrequestPerSecondPerApiKey int64  `mapstructure:"request_per_second_per_api_key"`
	AppPort                    string `mapstructure:"app_port"`
	RedisPort                  string `mapstructure:"redis_port"`
	RedisHost                  string `mapstructure:"redis_host"`
}

func LoadConfig(paths []string) (*Config, error) {
	var cfg *Config
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
