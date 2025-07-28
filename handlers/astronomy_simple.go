package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"weather-mcp-server/config"
)

// SimpleAstronomyResponse estructura simplificada sin moon_illumination
type SimpleAstronomyResponse struct {
	Location struct {
		Name      string  `json:"name"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		TzID      string  `json:"tz_id"`
		Localtime string  `json:"localtime"`
	} `json:"location"`
	Astronomy struct {
		Astro struct {
			Sunrise   string `json:"sunrise"`
			Sunset    string `json:"sunset"`
			Moonrise  string `json:"moonrise"`
			Moonset   string `json:"moonset"`
			MoonPhase string `json:"moon_phase"`
		} `json:"astro"`
	} `json:"astronomy"`
}

// GetAstronomySimple obtiene datos astronómicos de forma simplificada
func GetAstronomySimple(cfg *config.Config, params map[string]interface{}) (interface{}, error) {
	location, ok := params["location"].(string)
	if !ok || location == "" {
		return nil, fmt.Errorf("parámetro 'location' es requerido")
	}

	// Fecha opcional (por defecto hoy)
	date := time.Now().Format("2006-01-02")
	if dateParam, exists := params["date"]; exists {
		if dateStr, ok := dateParam.(string); ok && dateStr != "" {
			// Validar formato de fecha
			if _, err := time.Parse("2006-01-02", dateStr); err != nil {
				return nil, fmt.Errorf("formato de fecha inválido. Use YYYY-MM-DD")
			}
			date = dateStr
		}
	}

	// Construir URL
	baseURL := fmt.Sprintf("%s/astronomy.json", cfg.BaseURL)
	params_url := url.Values{}
	params_url.Add("key", cfg.WeatherAPIKey)
	params_url.Add("q", location)
	params_url.Add("dt", date)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params_url.Encode())

	// Hacer request
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error conectando con WeatherAPI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error de WeatherAPI (código %d): verificar ubicación, fecha y API key", resp.StatusCode)
	}

	// Decodificar respuesta
	var astroResp SimpleAstronomyResponse
	if err := json.NewDecoder(resp.Body).Decode(&astroResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	// Formatear respuesta legible
	result := fmt.Sprintf(`🌌 DATOS ASTRONÓMICOS
📍 Ubicación: %s, %s, %s
📅 Fecha: %s
🌐 Coordenadas: %.2f, %.2f
🕐 Hora local: %s

🌅 DATOS SOLARES
• Amanecer: %s
• Atardecer: %s

🌙 DATOS LUNARES
• Salida de luna: %s
• Puesta de luna: %s
• Fase lunar: %s

💡 Zona horaria: %s`,
		astroResp.Location.Name,
		astroResp.Location.Region,
		astroResp.Location.Country,
		date,
		astroResp.Location.Lat,
		astroResp.Location.Lon,
		astroResp.Location.Localtime,
		astroResp.Astronomy.Astro.Sunrise,
		astroResp.Astronomy.Astro.Sunset,
		astroResp.Astronomy.Astro.Moonrise,
		astroResp.Astronomy.Astro.Moonset,
		astroResp.Astronomy.Astro.MoonPhase,
		astroResp.Location.TzID,
	)

	return result, nil
}
