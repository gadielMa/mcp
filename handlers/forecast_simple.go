package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"weather-mcp-server/config"
)

// SimpleForecastResponse estructura simplificada sin datos astronómicos detallados
type SimpleForecastResponse struct {
	Location struct {
		Name      string  `json:"name"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		Localtime string  `json:"localtime"`
	} `json:"location"`
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Day  struct {
				MaxtempC      float64 `json:"maxtemp_c"`
				MintempC      float64 `json:"mintemp_c"`
				AvgtempC      float64 `json:"avgtemp_c"`
				MaxwindKph    float64 `json:"maxwind_kph"`
				TotalprecipMm float64 `json:"totalprecip_mm"`
				Avghumidity   float64 `json:"avghumidity"`
				UV            float64 `json:"uv"`
				Condition     struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
					Code int    `json:"code"`
				} `json:"condition"`
			} `json:"day"`
			Astro struct {
				Sunrise   string `json:"sunrise"`
				Sunset    string `json:"sunset"`
				Moonrise  string `json:"moonrise"`
				Moonset   string `json:"moonset"`
				MoonPhase string `json:"moon_phase"`
			} `json:"astro"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

// GetForecastSimple obtiene el pronóstico del tiempo de forma simplificada
func GetForecastSimple(cfg *config.Config, params map[string]interface{}) (interface{}, error) {
	location, ok := params["location"].(string)
	if !ok || location == "" {
		return nil, fmt.Errorf("parámetro 'location' es requerido")
	}

	// Parámetros opcionales
	days := 3
	if daysParam, exists := params["days"]; exists {
		if daysFloat, ok := daysParam.(float64); ok {
			days = int(daysFloat)
		}
	}
	if days < 1 || days > 10 {
		days = 3
	}

	// Construir URL
	baseURL := fmt.Sprintf("%s/forecast.json", cfg.BaseURL)
	params_url := url.Values{}
	params_url.Add("key", cfg.WeatherAPIKey)
	params_url.Add("q", location)
	params_url.Add("days", strconv.Itoa(days))
	params_url.Add("aqi", "no")
	params_url.Add("alerts", "no")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params_url.Encode())

	// Hacer request
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error conectando con WeatherAPI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error de WeatherAPI (código %d): verificar ubicación y API key", resp.StatusCode)
	}

	// Decodificar respuesta
	var forecastResp SimpleForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	// Formatear respuesta legible
	result := fmt.Sprintf(`🌦️ PRONÓSTICO DEL TIEMPO (%d días)
📍 Ubicación: %s, %s, %s
🕐 Hora local: %s

`,
		days,
		forecastResp.Location.Name,
		forecastResp.Location.Region,
		forecastResp.Location.Country,
		forecastResp.Location.Localtime,
	)

	// Agregar cada día del pronóstico
	for i, day := range forecastResp.Forecast.Forecastday {
		result += fmt.Sprintf(`📅 DÍA %d - %s
🌡️  Temperatura: %.1f°C - %.1f°C (promedio: %.1f°C)
🌦️  Condición: %s
💧 Humedad promedio: %.0f%%
💨 Viento máximo: %.1f km/h
🌧️  Precipitación total: %.1f mm
☀️  Índice UV: %.1f
🌅 Amanecer: %s | 🌇 Atardecer: %s
🌙 Fase lunar: %s

`,
			i+1,
			day.Date,
			day.Day.MintempC,
			day.Day.MaxtempC,
			day.Day.AvgtempC,
			day.Day.Condition.Text,
			day.Day.Avghumidity,
			day.Day.MaxwindKph,
			day.Day.TotalprecipMm,
			day.Day.UV,
			day.Astro.Sunrise,
			day.Astro.Sunset,
			day.Astro.MoonPhase,
		)
	}

	return result, nil
}
