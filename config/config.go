package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("config", fx.Provide(NewConfig), fx.Invoke(func(c *Config) error {
	return c.Validate()
}))

type Config struct {
	lg         *zap.Logger
	DSN        string
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	IsProd     bool
}

func (c *Config) Validate() error {
	v := reflect.ValueOf(c).Elem()
	var e bool
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			c.lg.Error(fmt.Sprintf("missing required environment variable: %s", v.Type().Field(i).Name))
			e = true
		}
	}
	if e {
		return fmt.Errorf("missing required environment variables")
	}
	return nil
}

func NewConfig(lg *zap.Logger) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		lg:         lg,
		DSN:        "",
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		IsProd:     os.Getenv("IS_PROD") == "true",
	}
	cfg.DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	return cfg, nil
}
