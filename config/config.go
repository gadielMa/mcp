package config

import (
	"fmt"
	"os"
)

// Config contiene la configuración del servidor
type Config struct {
	WeatherAPIKey string
	BaseURL       string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() (*Config, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY es requerida. Obtén una gratis en https://www.weatherapi.com/")
	}

	return &Config{
		WeatherAPIKey: apiKey,
		BaseURL:       "https://api.weatherapi.com/v1",
	}, nil
}

// MaskAPIKey enmascara la API key para logging seguro
func (c *Config) MaskAPIKey() string {
	if len(c.WeatherAPIKey) < 8 {
		return "****"
	}
	return c.WeatherAPIKey[:4] + "..." + c.WeatherAPIKey[len(c.WeatherAPIKey)-4:]
}
