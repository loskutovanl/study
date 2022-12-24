package config

import "os"

// DatabaseConfig определяет поля для подключения к postgre sql
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

// New возвращает экземпляр DatabaseConfig
func New() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("host"),
		Port:     getEnv("port"),
		User:     getEnv("user"),
		Password: getEnv("password"),
		Dbname:   getEnv("dbname"),
	}
}

// getEnv возвращает значение из переменных окружения или пустую строку в случае отсутствия значения
func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}
