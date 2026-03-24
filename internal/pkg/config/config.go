package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App           AppConfig           `mapstructure:"app"`
	Server        ServerConfig        `mapstructure:"server"`
	LLM           LLMConfig           `mapstructure:"llm"`
	DataSources   DataSourcesConfig   `mapstructure:"datasources"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Observability ObservabilityConfig `mapstructure:"observability"`
}

type AppConfig struct {
	Name     string `mapstructure:"name"`
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
}

type ServerConfig struct {
	HTTPPort     int    `mapstructure:"http_port"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

type LLMConfig struct {
	Provider    string  `mapstructure:"provider"`
	APIKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

type DataSourcesConfig struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Loki       LokiConfig       `mapstructure:"loki"`
}

type PrometheusConfig struct {
	URL     string `mapstructure:"url"`
	Timeout string `mapstructure:"timeout"`
}

type LokiConfig struct {
	URL     string `mapstructure:"url"`
	Timeout string `mapstructure:"timeout"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type PostgresConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type ObservabilityConfig struct {
	Tracing TracingConfig `mapstructure:"tracing"`
	Metrics MetricsConfig `mapstructure:"metrics"`
}

type TracingConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

type MetricsConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
