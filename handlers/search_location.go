package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"weather-mcp-server/config"
	"weather-mcp-server/models"
)

// SearchLocations busca ubicaciones por nombre
func SearchLocations(cfg *config.Config, params map[string]interface{}) (interface{}, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("parámetro 'query' es requerido")
	}

	// Construir URL
	baseURL := fmt.Sprintf("%s/search.json", cfg.BaseURL)
	params_url := url.Values{}
	params_url.Add("key", cfg.WeatherAPIKey)
	params_url.Add("q", query)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params_url.Encode())

	// Hacer request
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error conectando con WeatherAPI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error de WeatherAPI (código %d): verificar consulta y API key", resp.StatusCode)
	}

	// Decodificar respuesta
	var searchResp models.SearchLocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	if len(searchResp) == 0 {
		return "❌ No se encontraron ubicaciones para la consulta: " + query, nil
	}

	// Formatear respuesta legible
	result := fmt.Sprintf(`🔍 BÚSQUEDA DE UBICACIONES
📝 Consulta: "%s"
📍 Resultados encontrados: %d

`, query, len(searchResp))

	for i, location := range searchResp {
		result += fmt.Sprintf(`%d. %s
   📍 Región: %s, %s
   🌐 Coordenadas: %.2f, %.2f
   🔗 ID: %d

`,
			i+1,
			location.Name,
			location.Region,
			location.Country,
			location.Lat,
			location.Lon,
			location.ID,
		)
	}

	result += "💡 Tip: Puedes usar cualquiera de estos nombres en las otras herramientas del clima."

	return result, nil
}
