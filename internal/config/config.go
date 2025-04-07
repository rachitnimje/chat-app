package config

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Port int
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

func DefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "cricket360",
			DBName:   "yap_up",
		},
		Server: ServerConfig{
			Port: 8080,
		},
	}
}
