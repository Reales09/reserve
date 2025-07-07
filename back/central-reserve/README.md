# 🚀 Proyecto Base Backend en Go

Este repositorio sirve como una plantilla robusta y escalable para iniciar proyectos de backend en Go. Incluye una arquitectura limpia, configuración para dos tipos de servidores (HTTP y gRPC), conexión a base de datos, y flujos de trabajo de desarrollo automatizados.

---

## ✨ Características Principales

- **🌐 Servidor HTTP**: Implementado con [Gin](https://gin-gonic.com/), uno de los frameworks más rápidos y populares de Go.
- **🔌 Servidor gRPC**: Listo para comunicación de alto rendimiento entre microservicios.
- **🗄️ Base de Datos**: Configurado para [PostgreSQL](https://www.postgresql.org/), con un repositorio listo para usar.
- **📧 Sistema de Email**: Notificaciones automáticas por email para confirmaciones y cancelaciones de reservas.
- **📄 Documentación de API**:
    - **OpenAPI (Swagger)** para el servidor HTTP, totalmente interactiva.
    - **HTML Estático** para los servicios gRPC, con estilos personalizados.
- **⚙️ Tareas Automatizadas**: Un `Makefile` para simplificar tareas comunes como la generación de documentación.
- **📝 Logging Estructurado**: Logs claros y consistentes para facilitar la depuración.
- **🔑 Gestión de Entorno**: Carga de configuración desde archivos `.env`.
- **🐳 Soporte para Docker**: Preparado para ser contenedorizado.

---

## 📋 Prerrequisitos

Antes de empezar, asegúrate de tener instalado lo siguiente:

- **Go**: Versión 1.18 o superior.
- **Make**: Para ejecutar los comandos del `Makefile`.
- **Docker**: (Opcional) Si deseas levantar la base de datos PostgreSQL con Docker.

---

## 🚀 Guía de Inicio Rápido

Sigue estos pasos para poner en marcha el proyecto en tu máquina local:

1.  **Clonar el repositorio:**
    ```bash
    git clone [URL_DEL_REPOSITORIO]
    cd central_reserve
    ```

2.  **Configurar las variables de entorno:**
    ```bash
    # Crear archivo .env basado en las variables requeridas
    touch .env
    ```
    
    **⚠️ IMPORTANTE - Seguridad de Variables de Entorno:**
    
    Las siguientes variables son **OBLIGATORIAS** y contienen información sensible:
    ```bash
    # Configuración de la aplicación
    APP_ENV=development
    HTTP_PORT=3050
    GRPC_PORT=9090
    LOG_LEVEL=debug
    
    # 🔐 CRÍTICO: Usa un JWT secret fuerte en producción
    JWT_SECRET=your-super-secret-jwt-key-here-change-this-in-production
    
    # 🗄️ Configuración de base de datos PostgreSQL
    DB_HOST=localhost
    DB_USER=your_db_user
    DB_PASS=your_db_password
    DB_PORT=5432
    DB_NAME=central_reserve
    DB_LOG_LEVEL=info
    PGSSLMODE=disable
    
    # 📚 Configuración de Swagger
    URL_BASE_SWAGGER=http://localhost:3050
    
    # 📧 Configuración de Email (Opcional)
    SMTP_HOST=smtp.gmail.com
    SMTP_PORT=587
    SMTP_USER=tu-email@gmail.com
    SMTP_PASS=tu-contraseña-de-aplicación
    FROM_EMAIL=reservas@trattorialabella.com
    ```
    
    **🛡️ Mejores Prácticas de Seguridad:**
    - ❌ **NUNCA** subas el archivo `.env` al repositorio
    - ❌ **NUNCA** hardcodees credenciales en el código
    - ✅ Usa diferentes valores para dev/staging/prod
    - ✅ Genera JWT secrets seguros: `openssl rand -base64 32`
    - ✅ Usa gestores de secretos en producción (AWS Secrets Manager, HashiCorp Vault, etc.)
    - ✅ Para Gmail, usa contraseñas de aplicación en lugar de tu contraseña normal

3.  **Instalar dependencias:**
    ```bash
    go mod tidy
    ```

4.  **Levantar la base de datos (Opcional):**
    Si usas Docker, puedes iniciar una instancia de PostgreSQL con:
    ```bash
    # (Asegúrate de tener un docker-compose.yml en la carpeta /docker)
    docker-compose -f docker/docker-compose.yml up -d
    ```

5.  **Ejecutar la aplicación:**
    ```bash
    go run ./cmd/main.go
    ```
    ¡El servidor debería estar corriendo! Los logs de inicio te mostrarán las URLs disponibles.

---

## 🐳 Despliegue con Docker

### 🚀 Inicio Rápido con Docker

```bash
# Opción 1: Script automatizado (Recomendado)
./scripts/build-docker.sh dev

# Opción 2: Makefile
make docker-dev

# Opción 3: Docker Compose directo
cd docker && docker-compose -f docker-compose.dev.yml up -d
```

### 📋 Servicios Incluidos
- **API Backend**: http://localhost:3050
- **Swagger Docs**: http://localhost:3050/docs
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **NATS**: localhost:4222
- **NATS Dashboard**: http://localhost:8111
- **Adminer (DB)**: http://localhost:8080

### 🔒 Características de Seguridad
- ✅ Usuario no-root para ejecución
- ✅ Imagen minimalista (Alpine)
- ✅ Variables sensibles NO hardcodeadas
- ✅ Certificados SSL incluidos
- ✅ Healthcheck configurado
- ✅ Red aislada para servicios
- ✅ Volúmenes persistentes para datos

### 📖 Documentación Completa
Para más detalles sobre Docker, consulta: [README-DOCKER.md](README-DOCKER.md)

---

## ☁️ Despliegue a AWS ECR

La imagen está disponible públicamente en AWS ECR:

```bash
# 🌐 Imagen pública disponible
docker pull public.ecr.aws/d3a6d4r1/cam/reserve:latest

# 🚀 Ejecutar desde ECR
docker run --env-file .env -p 3050:3050 public.ecr.aws/d3a6d4r1/cam/reserve:latest
```

### 📦 **Despliegue Automatizado**

Para desplegar nuevas versiones a ECR:

```bash
# Desplegar versión latest
./scripts/deploy.sh

# Desplegar versión específica
./scripts/deploy.sh v1.0.1

# Desplegar versión de desarrollo
./scripts/deploy.sh dev
```

### 🔧 **Configuración Inicial de ECR**

Si necesitas configurar ECR desde cero:

```bash
# 1. Configurar permisos IAM para ECR público
# Agregar política: AmazonElasticContainerRegistryPublicFullAccess
# O crear política personalizada con:
#   - ecr-public:*
#   - sts:GetServiceBearerToken

# 2. Hacer login
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

# 3. Usar el script de despliegue
./scripts/deploy.sh
```

### 📋 **Versiones Disponibles**

- `latest`: Última versión estable
- `v1.0.0`: Primera versión de producción
- `YYYYMMDD_HHMMSS`: Versiones con timestamp automático

**🌐 Galería ECR:** https://gallery.ecr.aws/d3a6d4r1/cam/reserve

---

## 🛠️ Comandos Disponibles

Hemos configurado un `Makefile` para simplificar algunas tareas:

-   **`make docs`**: Regenera toda la documentación de la API gRPC (lee los `.proto`, aplica estilos y personalizaciones).
-   **`make clean`**: Elimina los binarios de compilación y la documentación generada.

---

## 📚 Documentación de API

Una vez que el servidor esté corriendo, puedes acceder a la documentación en las siguientes rutas:

-   **HTTP (OpenAPI)**:
    -   Visita `http://localhost:[PUERTO_HTTP]/docs`

-   **gRPC (Estática)**:
    -   Visita `http://localhost:[PUERTO_HTTP]/grpc-docs`

*(Reemplaza `[PUERTO_HTTP]` por el puerto que configuraste en tu archivo `.env`)*

---

## 📧 Sistema de Email

El proyecto incluye un sistema completo de notificaciones por email que envía automáticamente:

- ✅ **Confirmaciones de reserva** cuando se crea una nueva reserva
- ✅ **Cancelaciones de reserva** cuando se cancela una reserva existente

### Características del Sistema de Email:
- **Envío asíncrono**: No bloquea la respuesta de la API
- **Templates HTML profesionales**: Diseño responsivo con branding del restaurante
- **Soporte múltiples proveedores**: Gmail, Outlook, SendGrid, etc.
- **Logging detallado**: Seguimiento completo de envíos y errores
- **Configuración flexible**: Variables de entorno para diferentes entornos

### Documentación Completa:
Para más detalles sobre la configuración y uso del sistema de email, consulta:
- 📖 [README-EMAIL.md](README-EMAIL.md) - Guía completa del sistema de email
- 📋 [env-template-email.txt](env-template-email.txt) - Ejemplos de configuración
- 🧪 [examples/email-test.go](examples/email-test.go) - Ejemplo de uso
