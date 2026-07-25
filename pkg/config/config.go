package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Neo4j    Neo4jConfig
	Python   PythonConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	Env      string `mapstructure:"env"`
	Debug    bool   `mapstructure:"debug"`
}

type DatabaseConfig struct {
	DSN             string `mapstructure:"dsn"`
	MaxConnections  int    `mapstructure:"max_connections"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Neo4jConfig struct {
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type PythonConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type SecurityConfig struct {
	JWTSecret       string `mapstructure:"jwt_secret"`
	APIKeyRequired  bool   `mapstructure:"api_key_required"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/hlfa")
	viper.AddConfigPath("$HOME/.hlfa")
	viper.AddConfigPath(".")

	// Default values
	setDefaults()

	// Read environment variables
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.env", "development")
	viper.SetDefault("server.debug", true)

	viper.SetDefault("database.max_connections", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 300)

	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("neo4j.uri", "bolt://localhost:7687")

	viper.SetDefault("python.host", "localhost")
	viper.SetDefault("python.port", 5000)

	// Get from environment variables with defaults
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		viper.SetDefault("database.dsn", dsn)
	} else {
		viper.SetDefault("database.dsn", "postgres://user:password@localhost:5432/hlfa")
	}
}
