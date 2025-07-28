# 🌤️ Weather MCP Server

Un servidor MCP (Model Context Protocol) en Go que proporciona herramientas del clima usando la API de WeatherAPI.com.

## 🚀 Características

- **4 Herramientas del clima** completamente funcionales
- **API gratuita** de WeatherAPI.com con hasta 1 millón de llamadas/mes
- **Respuestas en español** formateadas y legibles
- **Manejo robusto de errores** y validación
- **Configuración simple** via variables de entorno
- **Documentación completa** para presentaciones

## 🛠️ Herramientas Disponibles

### 1. `get_current_weather`
Obtiene el clima actual para cualquier ubicación.

**Parámetros:**
- `location` (requerido): Ciudad, código postal, coordenadas o IP
- `aqi` (opcional): Incluir calidad del aire (`yes`/`no`)

**Ejemplo de uso:**
```json
{
  "location": "Buenos Aires",
  "aqi": "yes"
}
```

### 2. `get_forecast`
Obtiene pronóstico del tiempo hasta 10 días.

**Parámetros:**
- `location` (requerido): Ciudad, código postal, coordenadas o IP
- `days` (opcional): Número de días (1-10, default: 3)
- `aqi` (opcional): Incluir calidad del aire (`yes`/`no`)
- `alerts` (opcional): Incluir alertas meteorológicas (`yes`/`no`)

**Ejemplo de uso:**
```json
{
  "location": "Madrid",
  "days": 5,
  "aqi": "yes"
}
```

### 3. `search_locations`
Busca ubicaciones por nombre para obtener información precisa.

**Parámetros:**
- `query` (requerido): Nombre de la ciudad o ubicación

**Ejemplo de uso:**
```json
{
  "query": "London"
}
```

### 4. `get_astronomy`
Obtiene datos astronómicos (amanecer, atardecer, fases lunares).

**Parámetros:**
- `location` (requerido): Ciudad, código postal, coordenadas o IP
- `date` (opcional): Fecha en formato YYYY-MM-DD (default: hoy)

**Ejemplo de uso:**
```json
{
  "location": "Santiago",
  "date": "2024-12-25"
}
```

## 📦 Instalación y Configuración

### 1. Obtener API Key
1. Regístrate gratis en [WeatherAPI.com](https://www.weatherapi.com/)
2. Copia tu API key del dashboard

### 2. Configurar Variables de Entorno

Copia tu WEATHER_API_KEY="tu_api_key_aqui" en el start.sh

### 3. Ejecutar el Servidor
```bash
./start.sh
```

El MCP Inspector estará disponible en `http://localhost:6274`
La API Rest estará disponible en `http://localhost:3001`

## 🏗️ Arquitectura del Proyecto

```
weather-mcp-server/
├── main.go                     # Servidor MCP principal
├── go.mod                      # Dependencias
├── config/
│   └── config.go              # Configuración y API key
├── handlers/
│   ├── current_weather.go     # Clima actual
│   ├── forecast.go            # Pronósticos
│   ├── search_location.go     # Búsqueda de ubicaciones
│   └── astronomy.go           # Datos astronómicos
├── models/
│   └── weather.go             # Modelos de datos
├── examples/
│   └── client_example.go      # Ejemplo de cliente
└── README.md                  # Esta documentación
```

## 🔧 Desarrollo

### Estructura MCP
El servidor implementa el protocolo MCP estándar:

- **Herramientas**: Funciones que pueden ser llamadas por modelos de IA
- **Esquemas**: Definición de parámetros de entrada
- **Respuestas**: Formato estándar de respuesta MCP

### Manejo de Errores
- Validación de parámetros requeridos
- Verificación de API key
- Manejo de errores de red
- Respuestas HTTP apropiadas

### Seguridad
- API key enmascarada en logs
- Validación de entrada
- Manejo seguro de errores

## 📊 Ejemplos de Respuesta

### Clima Actual
```
🌤️ CLIMA ACTUAL
📍 Ubicación: Buenos Aires, Buenos Aires F.D., Argentina
🌡️ Temperatura: 25.0°C (77.0°F)
🌦️ Condición: Parcialmente nublado
💨 Viento: 15.0 km/h NO
💧 Humedad: 65%
☁️ Nubosidad: 40%
👁️ Visibilidad: 10.0 km
🌡️ Sensación térmica: 27.0°C
📊 Presión: 1013.0 mb
🌧️ Precipitación: 0.0 mm
☀️ Índice UV: 6.0

⏰ Última actualización: 2024-01-20 15:30
```

### Pronóstico
```
🌦️ PRONÓSTICO DEL TIEMPO (3 días)
📍 Ubicación: Madrid, Madrid, España
🕐 Hora local: 2024-01-20 21:45

📅 DÍA 1 - 2024-01-20
🌡️ Temperatura: 8.0°C - 15.0°C
🌦️ Condición: Soleado
💧 Humedad promedio: 55%
💨 Viento máximo: 12.0 km/h
🌧️ Precipitación total: 0.0 mm
☀️ Índice UV: 3.0
🌅 Amanecer: 08:15 AM | 🌇 Atardecer: 06:30 PM
🌙 Fase lunar: Cuarto creciente (45% iluminada)
```

## 📈 Escalabilidad

### Performance
- Respuestas promedio: < 500ms
- Soporte para múltiples requests concurrentes
- Cache de conexiones HTTP

### Límites de API
- Plan gratuito: 1 millón llamadas/mes
- Planes pagos disponibles para mayor volumen
- Rate limiting implementado

## 🛡️ Mejores Prácticas

### Para Producción
1. **Variables de entorno**: Nunca hardcodear API keys
2. **Logging**: Implementar logging estructurado
3. **Monitoreo**: Métricas de health y performance
4. **HTTPS**: Usar TLS en producción
5. **Docker**: Containerizar para deployment

### Para Desarrollo
1. **Testing**: Unit tests para handlers
2. **Mocking**: Mock de API externa para tests
3. **Linting**: Usar golangci-lint
4. **Documentación**: Mantener README actualizado

## 🔗 Enlaces Útiles

- [WeatherAPI.com](https://www.weatherapi.com/) - Documentación de la API
- [MCP Specification](https://spec.modelcontextprotocol.io/) - Protocolo MCP oficial
- [Go Documentation](https://golang.org/doc/) - Documentación de Go

## 📄 Licencia

MIT License - Ver archivo LICENSE para detalles.

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature
3. Commit tus cambios
4. Push a la rama
5. Abre un Pull Request

## 🐛 Reportar Problemas

Si encuentras algún problema, por favor abre un issue con:
- Descripción del problema
- Pasos para reproducir
- Logs relevantes
- Versión de Go utilizada

---

**¡Listo para tu charla sobre MCP Servers! 🎤** 

## 🧪 Ejemplos de uso con curl

### 1. Health Check
```bash
curl -X GET http://localhost:3003/health
```

### 2. Listar herramientas disponibles
```bash
curl -X GET http://localhost:3003/tools
```

### 3. Obtener clima actual
```bash
curl -X POST http://localhost:3003/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "method": "tools/call",
    "params": {
      "name": "get_current_weather",
      "arguments": {
        "location": "Madrid"
      }
    }
  }'
```

### 4. Obtener pronóstico del tiempo (3 días)
```bash
curl -X POST http://localhost:3003/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "2",
    "method": "tools/call",
    "params": {
      "name": "get_forecast",
      "arguments": {
        "location": "Barcelona",
        "days": 3
      }
    }
  }'
```

### 5. Buscar ubicaciones
```bash
curl -X POST http://localhost:3003/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "3",
    "method": "tools/call",
    "params": {
      "name": "search_locations",
      "arguments": {
        "query": "Sevilla"
      }
    }
  }'
```

### 6. Obtener datos astronómicos
```bash
curl -X POST http://localhost:3003/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "4",
    "method": "tools/call",
    "params": {
      "name": "get_astronomy",
      "arguments": {
        "location": "Valencia",
        "date": "2024-12-19"
      }
    }
  }'
```

### 7. Listar herramientas (método MCP)
```bash
curl -X POST http://localhost:3003/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "5",
    "method": "tools/list",
    "params": {}
  }'
```

### Respuestas esperadas

#### Health Check
```json
{
  "status": "healthy",
  "service": "Weather MCP Server",
  "version": "1.0.0"
}
```

#### Clima actual
```json
{
  "id": "1",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "🌤️ Clima actual en Madrid: 15°C, Parcialmente nublado. Sensación térmica: 13°C..."
      }
    ]
  }
}
``` 