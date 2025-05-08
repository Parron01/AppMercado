// pkg/config/config.go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    ServerPort         string
    DBHost             string
    DBPort             string
    DBUser             string
    DBPassword         string
    DBName             string
    JWTSecret          string
    JWTExpirationHours int
}

// Load carrega as variáveis de ambiente
func Load() *Config {
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        panic("Erro ao ler o arquivo .env: " + err.Error())
    }

    return &Config{
        ServerPort:         viper.GetString("SERVER_PORT"),
        DBHost:             viper.GetString("DB_HOST"),
        DBPort:             viper.GetString("DB_PORT"),
        DBUser:             viper.GetString("DB_USER"),
        DBPassword:         viper.GetString("DB_PASSWORD"),
        DBName:             viper.GetString("DB_NAME"),
        JWTSecret:          viper.GetString("JWT_SECRET"),
        JWTExpirationHours: viper.GetInt("JWT_EXPIRATION_HOURS"),
    }
}
