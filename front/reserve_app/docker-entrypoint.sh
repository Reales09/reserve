#!/bin/bash

echo "Variables de entorno disponibles:"
echo "REACT_APP_API_BASE_URL: ${REACT_APP_API_BASE_URL:-http://localhost:3050}"

# Reemplazar en todos los archivos JS
for file in $(find /usr/share/nginx/html -name "*.js" -type f); do
  sed -i "s|REACT_APP_API_BASE_URL_PLACEHOLDER|${REACT_APP_API_BASE_URL:-http://localhost:3050}|g" "$file"
done

# Reemplazar en todos los archivos HTML
for file in $(find /usr/share/nginx/html -name "*.html" -type f); do
  sed -i "s|REACT_APP_API_BASE_URL_PLACEHOLDER|${REACT_APP_API_BASE_URL:-http://localhost:3050}|g" "$file"
done

# Reemplazar en todos los archivos JSON
for file in $(find /usr/share/nginx/html -name "*.json" -type f); do
  sed -i "s|REACT_APP_API_BASE_URL_PLACEHOLDER|${REACT_APP_API_BASE_URL:-http://localhost:3050}|g" "$file"
done

echo "Variables de entorno configuradas:"
echo "REACT_APP_API_BASE_URL: ${REACT_APP_API_BASE_URL:-http://localhost:3050}"

# Ejecutar el comando original
exec "$@" 