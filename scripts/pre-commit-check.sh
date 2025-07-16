#!/bin/bash

# Script de pre-commit para validar archivos grandes
# Uso: ./scripts/pre-commit-check.sh

echo "🔍 Verificando archivos grandes antes del commit..."

# Verificar archivos mayores a 10MB (solo archivos que git está trackeando)
LARGE_FILES=$(git ls-files | xargs -I {} find {} -type f -size +10M 2>/dev/null | head -10)

if [ ! -z "$LARGE_FILES" ]; then
    echo "❌ ERROR: Se encontraron archivos grandes (>10MB) en el repositorio:"
    echo "$LARGE_FILES"
    echo ""
    echo "💡 Recomendaciones:"
    echo "  - Verifica que estos archivos estén en .gitignore"
    echo "  - Si son necesarios, considera usar Git LFS"
    echo "  - Si son temporales, elimínalos antes del commit"
    exit 1
fi

# Verificar archivos de cache de BuildKit (solo los trackeados)
BUILDKIT_FILES=$(git ls-files | grep -E "\.buildkit/" | head -10)

if [ ! -z "$BUILDKIT_FILES" ]; then
    echo "❌ ERROR: Se encontraron archivos de cache de BuildKit en el repositorio:"
    echo "$BUILDKIT_FILES"
    echo ""
    echo "💡 Solución:"
    echo "  - Agrega las carpetas .buildkit/ al .gitignore"
    echo "  - Ejecuta: git rm --cached -r infra/.buildkit/ infra/scripts/.buildkit/"
    exit 1
fi

# Verificar binarios (solo los trackeados)
BINARY_FILES=$(git ls-files | xargs -I {} file {} 2>/dev/null | grep -E "(executable|binary)" | cut -d: -f1 | head -10)

if [ ! -z "$BINARY_FILES" ]; then
    echo "⚠️  ADVERTENCIA: Se encontraron archivos binarios en el repositorio:"
    echo "$BINARY_FILES"
    echo ""
    echo "💡 Verifica que estos archivos sean necesarios en el repositorio"
fi

echo "✅ Verificación completada. No se encontraron archivos problemáticos."
exit 0
