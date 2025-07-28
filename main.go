package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"weather-mcp-server/config"
	"weather-mcp-server/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// MCPServer representa nuestro servidor MCP
type MCPServer struct {
	config *config.Config
}

// MCPRequest estructura de request MCP estándar
type MCPRequest struct {
	ID     string                 `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

// MCPResponse estructura de response MCP estándar
type MCPResponse struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *MCPError   `json:"error,omitempty"`
}

// MCPError estructura de error MCP
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToolDefinition define las herramientas disponibles
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	server := &MCPServer{config: cfg}

	// Configurar rutas
	r := mux.NewRouter()

	// Rutas MCP estándar
	r.HandleFunc("/", server.handleMCPRequest).Methods("POST")
	r.HandleFunc("/tools", server.listTools).Methods("GET")
	r.HandleFunc("/health", server.healthCheck).Methods("GET")

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}

	fmt.Printf("🌤️  Servidor MCP Weather iniciado en puerto %s\n", port)
	fmt.Printf("🔧 Herramientas disponibles:\n")
	fmt.Printf("   - get_current_weather: Clima actual\n")
	fmt.Printf("   - get_forecast: Pronóstico del tiempo\n")
	fmt.Printf("   - search_locations: Buscar ubicaciones\n")
	fmt.Printf("   - get_astronomy: Datos astronómicos\n")
	fmt.Printf("📚 API Key: %s\n", cfg.MaskAPIKey())

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// handleMCPRequest maneja todas las solicitudes MCP
func (s *MCPServer) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, req.ID, 400, "Request JSON inválido")
		return
	}

	switch req.Method {
	case "tools/list":
		s.sendResponse(w, req.ID, s.getToolDefinitions())
	case "tools/call":
		s.handleToolCall(w, req)
	default:
		s.sendError(w, req.ID, 404, "Método no encontrado")
	}
}

// handleToolCall ejecuta las herramientas
func (s *MCPServer) handleToolCall(w http.ResponseWriter, req MCPRequest) {
	params, ok := req.Params["arguments"].(map[string]interface{})
	if !ok {
		s.sendError(w, req.ID, 400, "Argumentos inválidos")
		return
	}

	toolName, ok := req.Params["name"].(string)
	if !ok {
		s.sendError(w, req.ID, 400, "Nombre de herramienta requerido")
		return
	}

	var result interface{}
	var err error

	switch toolName {
	case "get_current_weather":
		result, err = handlers.GetCurrentWeather(s.config, params)
	case "get_forecast":
		result, err = handlers.GetForecastSimple(s.config, params)
	case "search_locations":
		result, err = handlers.SearchLocations(s.config, params)
	case "get_astronomy":
		result, err = handlers.GetAstronomySimple(s.config, params)
	default:
		s.sendError(w, req.ID, 404, "Herramienta no encontrada")
		return
	}

	if err != nil {
		s.sendError(w, req.ID, 500, err.Error())
		return
	}

	s.sendResponse(w, req.ID, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("%v", result),
			},
		},
	})
}

// listTools devuelve la lista de herramientas (endpoint GET)
func (s *MCPServer) listTools(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.getToolDefinitions())
}

// healthCheck endpoint de salud
func (s *MCPServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "Weather MCP Server",
		"version": "1.0.0",
	})
}

// getToolDefinitions devuelve las definiciones de todas las herramientas
func (s *MCPServer) getToolDefinitions() map[string]interface{} {
	return map[string]interface{}{
		"tools": []ToolDefinition{
			{
				Name:        "get_current_weather",
				Description: "Obtiene el clima actual para una ubicación específica",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "Nombre de la ciudad, código postal, coordenadas (lat,lon) o dirección IP",
						},
						"aqi": map[string]interface{}{
							"type":        "string",
							"description": "Incluir datos de calidad del aire (yes/no)",
							"default":     "no",
						},
					},
					"required": []string{"location"},
				},
			},
			{
				Name:        "get_forecast",
				Description: "Obtiene el pronóstico del tiempo para una ubicación",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "Nombre de la ciudad, código postal, coordenadas (lat,lon) o dirección IP",
						},
						"days": map[string]interface{}{
							"type":        "number",
							"description": "Número de días de pronóstico (1-10)",
							"default":     3,
						},
						"aqi": map[string]interface{}{
							"type":        "string",
							"description": "Incluir datos de calidad del aire (yes/no)",
							"default":     "no",
						},
						"alerts": map[string]interface{}{
							"type":        "string",
							"description": "Incluir alertas meteorológicas (yes/no)",
							"default":     "no",
						},
					},
					"required": []string{"location"},
				},
			},
			{
				Name:        "search_locations",
				Description: "Busca ubicaciones por nombre para obtener información detallada",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Nombre de la ciudad o ubicación a buscar",
						},
					},
					"required": []string{"query"},
				},
			},
			{
				Name:        "get_astronomy",
				Description: "Obtiene datos astronómicos para una ubicación y fecha",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "Nombre de la ciudad, código postal, coordenadas (lat,lon) o dirección IP",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Fecha en formato YYYY-MM-DD (opcional, por defecto hoy)",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}
}

// sendResponse envía una respuesta MCP exitosa
func (s *MCPServer) sendResponse(w http.ResponseWriter, id string, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MCPResponse{
		ID:     id,
		Result: result,
	})
}

// sendError envía una respuesta MCP de error
func (s *MCPServer) sendError(w http.ResponseWriter, id string, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MCPResponse{
		ID: id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	})
}
