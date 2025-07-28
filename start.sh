#!/bin/bash

echo "🌤️  INICIANDO WEATHER MCP SERVER - DEMO COMPLETA"
echo "================================================"
echo ""

# Configurar entorno
export WEATHER_API_KEY="f0f8a22058004b99b2d181316252807"

# Verificar directorio
echo "📁 Directorio: $(pwd)"
echo "🔑 API Key: ${WEATHER_API_KEY:0:8}..."
echo ""

# Verificar prerrequisitos
echo "🔍 Verificando prerrequisitos..."
if ! command -v go &> /dev/null; then
    echo "❌ Go no está instalado"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "❌ npm no está instalado" 
    exit 1
fi

echo "✅ Go: $(go version | cut -d' ' -f3)"
echo "✅ npm: $(npm --version)"
echo ""

# Probar servidor stdio
echo "🧪 Probando servidor MCP stdio..."
TEST_RESULT=$(echo '{"id":1,"method":"initialize","params":{}}' | go run mcp_stdio_server.go 2>&1 | head -n 1)
if [[ $TEST_RESULT == *"result"* ]] && [[ $TEST_RESULT == *"serverInfo"* ]]; then
    echo "✅ Servidor MCP stdio funciona correctamente"
else
    echo "❌ Error en servidor stdio - Respuesta: $TEST_RESULT"
    echo "⚠️  Continuando con demo HTTP..."
    echo ""
fi
echo ""

# Menú de opciones
echo "🚀 OPCIONES PARA TU CHARLA:"
echo "=========================="
echo ""
echo "1. 📱 SERVIDOR HTTP (Recomendado para demos rápidas)"
echo "2. 🔍 MCP INSPECTOR (Más visual e interactivo)"
echo "3. 🧪 AMBOS (Para demo completa)"
echo ""

read -p "Elige una opción (1/2/3): " choice

case $choice in
    1)
        echo ""
        echo "🚀 Iniciando Servidor HTTP..."
        echo "URL: http://localhost:3003"
        echo ""
        go run main.go
        ;;
    2)
        echo ""
        echo "🔍 Iniciando MCP Inspector..."
        echo ""
        echo "📋 CONFIGURACIÓN PARA EL INSPECTOR:"
        echo "Command: go run mcp_stdio_server.go"
        echo "Working Directory: $(pwd)"
        echo "Environment: WEATHER_API_KEY=f0f8a22058004b99b2d181316252807"
        echo ""
        sleep 3
        npx @modelcontextprotocol/inspector
        ;;
    3)
        echo ""
        echo "🚀 Iniciando AMBOS servidores..."
        echo ""
        
        # Servidor HTTP en background
        echo "Iniciando servidor HTTP en puerto 3003..."
        go run main.go &
        HTTP_PID=$!
        sleep 3
        
        # Verificar que el HTTP funciona
        if curl -s http://localhost:3003/health > /dev/null; then
            echo "✅ Servidor HTTP listo en http://localhost:3003"
        else
            echo "❌ Error en servidor HTTP"
            kill $HTTP_PID 2>/dev/null
            exit 1
        fi
        
        echo ""
        echo "📋 CONFIGURACIÓN PARA MCP INSPECTOR:"
        echo "Command: go run mcp_stdio_server.go"
        echo "Working Directory: $(pwd)"
        echo "Environment: WEATHER_API_KEY=f0f8a22058004b99b2d181316252807"
        echo ""
        echo "Iniciando MCP Inspector en 3 segundos..."
        sleep 3
        
        # Iniciar MCP Inspector
        npx @modelcontextprotocol/inspector
        
        # Cleanup al salir
        echo "Cerrando servidor HTTP..."
        kill $HTTP_PID 2>/dev/null
        ;;
    *)
        echo "❌ Opción inválida"
        exit 1
        ;;
esac 