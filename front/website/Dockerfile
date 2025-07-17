# Etapa 1: Build del proyecto Astro
FROM --platform=linux/arm64 node:20-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias y instalarlas
COPY package*.json ./
RUN npm ci

# Copiar el resto del código fuente
COPY . .

# Limpiar y reinstalar dependencias para evitar bug de npm en ARM64
RUN rm -rf node_modules package-lock.json && npm ci

# Construir la aplicación
RUN npm run build

# Etapa 2: Servir los archivos estáticos con Nginx
FROM --platform=linux/arm64 nginx:alpine

# Copiar el build de la etapa anterior
COPY --from=builder /app/dist /usr/share/nginx/html

# Copiar la nueva configuración de Nginx (que incluye el reverse proxy)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Exponer el puerto 80
EXPOSE 80

# Iniciar Nginx directamente
CMD ["nginx", "-g", "daemon off;"]
