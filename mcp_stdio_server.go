package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"weather-mcp-server/config"
	"weather-mcp-server/handlers"
)

// MCPStdioRequest estructura de request MCP estándar
type MCPStdioRequest struct {
	ID     interface{}            `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// MCPStdioResponse estructura de response MCP estándar
type MCPStdioResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	ID      interface{}    `json:"id"`
	Result  interface{}    `json:"result,omitempty"`
	Error   *MCPStdioError `json:"error,omitempty"`
}

// MCPStdioError estructura de error MCP
type MCPStdioError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToolDef define las herramientas disponibles
type ToolDef struct {
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

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req MCPStdioRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(req.ID, 400, "Request JSON inválido")
			continue
		}

		handleRequest(cfg, req)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error leyendo stdin: %v", err)
	}
}

func handleRequest(cfg *config.Config, req MCPStdioRequest) {
	switch req.Method {
	case "initialize":
		handleInitialize(req)
	case "tools/list":
		handleToolsList(req)
	case "tools/call":
		handleToolCall(cfg, req)
	default:
		sendError(req.ID, 404, "Método no encontrado: "+req.Method)
	}
}

func handleInitialize(req MCPStdioRequest) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "weather-mcp-server",
			"version": "1.0.0",
		},
	}
	sendResponse(req.ID, result)
}

func handleToolsList(req MCPStdioRequest) {
	tools := []ToolDef{
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
	}

	sendResponse(req.ID, map[string]interface{}{
		"tools": tools,
	})
}

func handleToolCall(cfg *config.Config, req MCPStdioRequest) {
	params, ok := req.Params["arguments"].(map[string]interface{})
	if !ok {
		sendError(req.ID, 400, "Argumentos inválidos")
		return
	}

	toolName, ok := req.Params["name"].(string)
	if !ok {
		sendError(req.ID, 400, "Nombre de herramienta requerido")
		return
	}

	var result interface{}
	var err error

	switch toolName {
	case "get_current_weather":
		result, err = handlers.GetCurrentWeather(cfg, params)
	case "get_forecast":
		result, err = handlers.GetForecastSimple(cfg, params)
	case "search_locations":
		result, err = handlers.SearchLocations(cfg, params)
	case "get_astronomy":
		result, err = handlers.GetAstronomySimple(cfg, params)
	default:
		sendError(req.ID, 404, "Herramienta no encontrada")
		return
	}

	if err != nil {
		sendError(req.ID, 500, err.Error())
		return
	}

	sendResponse(req.ID, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("%v", result),
			},
		},
	})
}

func sendResponse(id interface{}, result interface{}) {
	response := MCPStdioResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
	}

	jsonData, _ := json.Marshal(response)
	fmt.Println(string(jsonData))
}

func sendError(id interface{}, code int, message string) {
	response := MCPStdioResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Error: &MCPStdioError{
			Code:    code,
			Message: message,
		},
	}

	jsonData, _ := json.Marshal(response)
	fmt.Println(string(jsonData))
}
