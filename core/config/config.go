package config

import "os"

type Config struct {
	AWS AWSConfig
	DB  DBConfig
}

type AWSConfig struct {
	Region       string
	UserPoolID   string
	ClientID     string
	ClientSecret string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Load() Config {
	return Config{
		AWS: AWSConfig{
			Region:       os.Getenv("AWS_REGION"),
			UserPoolID:   os.Getenv("COGNITO_USER_POOL_ID"),
			ClientID:     os.Getenv("COGNITO_CLIENT_ID"),
			ClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
	}
}
