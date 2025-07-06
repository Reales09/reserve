# Multi-stage build para React
FROM node:18-alpine AS builder

# Argumentos de build para variables de entorno
ARG REACT_APP_API_BASE_URL=REACT_APP_API_BASE_URL_PLACEHOLDER
ARG REACT_APP_NAME="Reserve App"
ARG REACT_APP_VERSION="1.0.0"

# Establecer variables de entorno para el build
ENV REACT_APP_API_BASE_URL=${REACT_APP_API_BASE_URL}
ENV REACT_APP_NAME=${REACT_APP_NAME}
ENV REACT_APP_VERSION=${REACT_APP_VERSION}

# Establecer directorio de trabajo
WORKDIR /app

# Copiar package.json y package-lock.json
COPY package*.json ./

# Instalar dependencias
RUN npm ci --only=production

# Copiar el código fuente
COPY . .

# Construir la aplicación para producción
RUN npm run build

# Etapa de producción con nginx
FROM nginx:alpine

# Copiar archivos construidos desde la etapa anterior
COPY --from=builder /app/build /usr/share/nginx/html

# Copiar configuración personalizada de nginx
COPY nginx.conf /etc/nginx/nginx.conf

# Crear script para reemplazar variables de entorno en tiempo de ejecución
RUN apk add --no-cache bash
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Exponer puerto 80
EXPOSE 80

# Usar script personalizado como entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"] 