# 🚀 Guía de Despliegue - Central Reserve

## 📋 Resumen

**Central Reserve** está desplegado en **AWS ECR Público** y listo para usar en cualquier entorno.

- **📦 Imagen**: `public.ecr.aws/d3a6d4r1/cam/reserve`
- **📏 Tamaño**: 55.4MB (optimizada con Alpine Linux)
- **🔒 Seguridad**: Usuario no-root, imagen minimalista
- **🌐 Galería**: https://gallery.ecr.aws/d3a6d4r1/cam/reserve

---

## 🚀 Uso Rápido

### 1. **Ejecutar directamente desde ECR**
```bash
# Crear archivo .env con tus variables
touch .env

# Ejecutar la aplicación
docker run --env-file .env -p 3050:3050 public.ecr.aws/d3a6d4r1/cam/reserve:latest
```

### 2. **Usando docker-compose (Producción)**
```bash
# Usar el stack completo de producción
docker-compose -f docker/docker-compose.prod.yml up -d

# Verificar que esté funcionando
curl http://localhost:3050/health
```

### 3. **Descargar y ejecutar localmente**
```bash
# Descargar imagen
docker pull public.ecr.aws/d3a6d4r1/cam/reserve:latest

# Ejecutar
docker run --env-file .env -p 3050:3050 public.ecr.aws/d3a6d4r1/cam/reserve:latest
```

---

## 📋 Versiones Disponibles

| Tag | Descripción | Comando |
|-----|-------------|---------|
| `latest` | Última versión estable | `docker pull public.ecr.aws/d3a6d4r1/cam/reserve:latest` |
| `v1.0.0` | Primera versión de producción | `docker pull public.ecr.aws/d3a6d4r1/cam/reserve:v1.0.0` |
| `v1.0.1` | Versión mejorada | `docker pull public.ecr.aws/d3a6d4r1/cam/reserve:v1.0.1` |

---

## 🔧 Despliegue y Desarrollo

### **Desplegar nueva versión**
```bash
# Versión automática (latest + timestamp)
./scripts/deploy.sh

# Versión específica
./scripts/deploy.sh v1.0.2

# Versión de desarrollo
./scripts/deploy.sh dev
```

### **Desarrollo local**
```bash
# Construir imagen local
docker build -f docker/Dockerfile -t central-reserve .

# Ejecutar en desarrollo
docker run --env-file .env -p 3050:3050 central-reserve
```

### **CI/CD Automático**
El proyecto incluye GitHub Actions que automáticamente:
- ✅ Ejecuta tests
- ✅ Construye la imagen
- ✅ Deploya a ECR en cada push a `main`
- ✅ Crear tags automáticos para releases

---

## 🌐 Configuración de Entornos

### **Desarrollo**
```bash
# .env para desarrollo
APP_ENV=development
HTTP_PORT=3050
LOG_LEVEL=debug
DB_HOST=localhost
# ... más variables
```

### **Producción**
```bash
# .env para producción
APP_ENV=production
HTTP_PORT=3050
LOG_LEVEL=info
DB_HOST=prod-database-host
JWT_SECRET=production-super-secret-key
# ... más variables
```

### **Staging**
```bash
# Usar tag específico para staging
docker run --env-file .env.staging -p 3050:3050 public.ecr.aws/d3a6d4r1/cam/reserve:v1.0.1
```

---

## 📊 Monitoreo y Salud

### **Healthcheck**
```bash
# Verificar salud de la aplicación
curl http://localhost:3050/health

# Respuesta esperada: 200 OK
```

### **Logs**
```bash
# Ver logs en tiempo real
docker logs -f central_reserve_prod

# Logs con docker-compose
docker-compose -f docker/docker-compose.prod.yml logs -f central_reserve
```

### **Métricas**
```bash
# Swagger UI disponible en:
http://localhost:3050/docs

# API docs:
http://localhost:3050/api/v1/docs
```

---

## 🛠️ Comandos Útiles

### **Gestión de Imágenes**
```bash
# Limpiar imágenes locales
docker image prune -f

# Ver todas las imágenes del proyecto
docker images | grep central-reserve

# Eliminar imagen específica
docker rmi public.ecr.aws/d3a6d4r1/cam/reserve:old-version
```

### **Troubleshooting**
```bash
# Ejecutar contenedor en modo interactivo
docker run -it --env-file .env public.ecr.aws/d3a6d4r1/cam/reserve:latest sh

# Verificar variables de entorno
docker run --env-file .env public.ecr.aws/d3a6d4r1/cam/reserve:latest env

# Verificar conectividad a base de datos
docker run --env-file .env --rm public.ecr.aws/d3a6d4r1/cam/reserve:latest ping $DB_HOST
```

---

## 🔐 Seguridad

### **Variables de Entorno**
- ❌ **NUNCA** hardcodear credenciales en la imagen
- ✅ Usar archivos `.env` diferentes por entorno
- ✅ Rotar credenciales regularmente
- ✅ Usar gestores de secretos en producción

### **Configuración Segura**
```bash
# Generar JWT secret fuerte
openssl rand -base64 32

# Ejecutar con usuario no-root (ya configurado)
docker run --user 1000:1000 --env-file .env -p 3050:3050 public.ecr.aws/d3a6d4r1/cam/reserve:latest
```

---

## 📞 Soporte

### **Información del Sistema**
- **Go Version**: 1.23
- **Base Image**: Alpine Linux 3.19
- **Architecture**: Multi-stage build optimizado
- **Size**: 55.4MB
- **User**: appuser (non-root)

### **Puertos**
- **HTTP**: 3050
- **Healthcheck**: 3050/health
- **Docs**: 3050/docs

### **Contacto**
- **Repositorio**: https://github.com/your-repo/central-reserve
- **ECR Gallery**: https://gallery.ecr.aws/d3a6d4r1/cam/reserve
- **Issues**: GitHub Issues

---

## 🎯 Próximos Pasos

1. **Configurar monitoreo** (Prometheus, Grafana)
2. **Implementar alertas** (PagerDuty, Slack)
3. **Configurar backup automático** de la base de datos
4. **Implementar scaling horizontal** (Docker Swarm, Kubernetes)
5. **Configurar CDN** para assets estáticos 