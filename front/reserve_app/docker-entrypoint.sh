#!/bin/bash

# Función para reemplazar variables de entorno en archivos JavaScript
replace_env_vars() {
    local file="$1"
    local temp_file="${file}.tmp"
    
    # Reemplazar REACT_APP_API_BASE_URL
    sed "s|REACT_APP_API_BASE_URL_PLACEHOLDER|${REACT_APP_API_BASE_URL:-http://localhost:3050}|g" "$file" > "$temp_file"
    mv "$temp_file" "$file"
}

# Buscar y reemplazar variables de entorno en archivos JavaScript
find /usr/share/nginx/html -name "*.js" -type f -exec bash -c 'replace_env_vars "$0"' {} \;

# Función para reemplazar variables de entorno
export -f replace_env_vars

echo "Variables de entorno configuradas:"
echo "REACT_APP_API_BASE_URL: ${REACT_APP_API_BASE_URL:-http://localhost:3050}"

# Ejecutar el comando original
exec "$@" 