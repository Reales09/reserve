#!/bin/bash

# Script para limpiar caches y archivos temporales
# Uso: ./scripts/clean-cache.sh

echo "🧹 Limpiando caches y archivos temporales..."

# Limpiar cache de BuildKit
echo "📦 Limpiando cache de BuildKit..."
docker buildx prune -f 2>/dev/null || echo "⚠️  Docker BuildKit no disponible"

# Limpiar carpetas de cache locales
echo "🗂️  Limpiando carpetas de cache locales..."
rm -rf infra/.buildkit/ 2>/dev/null
rm -rf infra/scripts/.buildkit/ 2>/dev/null
rm -rf .buildkit/ 2>/dev/null

# Limpiar binarios de Go
echo "🔧 Limpiando binarios de Go..."
find . -name "main" -type f -delete 2>/dev/null
find . -path "*/bin/*" -type f -delete 2>/dev/null

# Limpiar archivos temporales
echo "📄 Limpiando archivos temporales..."
find . -name "*.tmp" -delete 2>/dev/null
find . -name "*.log" -delete 2>/dev/null
find . -name "*.cache" -delete 2>/dev/null

# Limpiar node_modules (opcional)
read -p "¿Deseas eliminar node_modules? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "📦 Eliminando node_modules..."
    find . -name "node_modules" -type d -exec rm -rf {} + 2>/dev/null
fi

echo "✅ Limpieza completada!"
echo ""
echo "💡 Para verificar que todo esté limpio, ejecuta:"
echo "   ./scripts/pre-commit-check.sh" 