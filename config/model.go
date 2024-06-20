package config

type Config struct {
	Binance   BinanceConfig   `mapstructure:"binance"`
	Broadcast BroadcastConfig `mapstructure:"broadcast"`
}

type BinanceConfig struct {
	URL       string `mapstructure:"url"`
	ApiKey    string `mapstructure:"api_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type BroadcastConfig struct {
	URL                 string `mapstructure:"url"`
	CheckStatusInterval int    `mapstructure:"check_status_interval"`
}
