package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"weather-mcp-server/config"
	"weather-mcp-server/models"
)

// GetCurrentWeather obtiene el clima actual para una ubicación
func GetCurrentWeather(cfg *config.Config, params map[string]interface{}) (interface{}, error) {
	location, ok := params["location"].(string)
	if !ok || location == "" {
		return nil, fmt.Errorf("parámetro 'location' es requerido")
	}

	// Parámetro opcional para calidad del aire
	aqi := "no"
	if aqiParam, exists := params["aqi"]; exists {
		if aqiStr, ok := aqiParam.(string); ok {
			aqi = aqiStr
		}
	}

	// Construir URL
	baseURL := fmt.Sprintf("%s/current.json", cfg.BaseURL)
	params_url := url.Values{}
	params_url.Add("key", cfg.WeatherAPIKey)
	params_url.Add("q", location)
	params_url.Add("aqi", aqi)

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
	var weatherResp models.CurrentWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	// Formatear respuesta legible
	result := fmt.Sprintf(`🌤️ CLIMA ACTUAL
📍 Ubicación: %s, %s, %s
🌡️  Temperatura: %.1f°C (%.1f°F)
🌦️  Condición: %s
💨 Viento: %.1f km/h %s
💧 Humedad: %d%%
☁️  Nubosidad: %d%%
👁️  Visibilidad: %.1f km
🌡️  Sensación térmica: %.1f°C
📊 Presión: %.1f mb
🌧️  Precipitación: %.1f mm
☀️  Índice UV: %.1f

⏰ Última actualización: %s`,
		weatherResp.Location.Name,
		weatherResp.Location.Region,
		weatherResp.Location.Country,
		weatherResp.Current.TempC,
		weatherResp.Current.TempF,
		weatherResp.Current.Condition.Text,
		weatherResp.Current.WindKph,
		weatherResp.Current.WindDir,
		weatherResp.Current.Humidity,
		weatherResp.Current.Cloud,
		weatherResp.Current.VisKm,
		weatherResp.Current.FeelslikeC,
		weatherResp.Current.PressureMb,
		weatherResp.Current.PrecipMm,
		weatherResp.Current.UV,
		weatherResp.Current.LastUpdated,
	)

	// Agregar calidad del aire si está disponible
	if weatherResp.Current.AirQuality != nil {
		result += fmt.Sprintf(`

🌬️  CALIDAD DEL AIRE
• CO: %.1f µg/m³
• NO2: %.1f µg/m³  
• O3: %.1f µg/m³
• SO2: %.1f µg/m³
• PM2.5: %.1f µg/m³
• PM10: %.1f µg/m³
• Índice EPA: %d
• Índice DEFRA: %d`,
			weatherResp.Current.AirQuality.CO,
			weatherResp.Current.AirQuality.NO2,
			weatherResp.Current.AirQuality.O3,
			weatherResp.Current.AirQuality.SO2,
			weatherResp.Current.AirQuality.PM25,
			weatherResp.Current.AirQuality.PM10,
			weatherResp.Current.AirQuality.USEPAIndex,
			weatherResp.Current.AirQuality.GBDefraIndex,
		)
	}

	return result, nil
}
