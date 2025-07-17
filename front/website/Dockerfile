# Etapa 1: Build del proyecto Astro
FROM --platform=linux/arm64 node:20-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias y instalarlas
COPY package*.json ./
RUN npm ci

# Copiar el resto del código fuente
COPY . .

# Workaround para el bug de Rollup en ARM64
RUN npm install @rollup/rollup-linux-arm64-musl

# Construir la aplicación. Ya no necesitamos pasar variables de entorno.
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
