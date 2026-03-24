package config

type Config struct {
	App *AppConfig
	DB  *DBConfig
}

type DBConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"open_conns"`
	MaxIdleConns    int    `mapstructure:"idle_conns"`
	MaxIdleLifetime int    `mapstructure:"Idle_lifetime"`
	MaxRetries      int    `mapstructure:"db_retries"`
	RetryDelay      int    `mapstructure:"db_retry_delay"`
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}
