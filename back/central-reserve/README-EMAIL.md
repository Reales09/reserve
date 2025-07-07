# Funcionalidad de Email - Sistema de Reservas

## Descripción

Se ha implementado un sistema de notificaciones por email que envía automáticamente confirmaciones y cancelaciones de reservas a los clientes.

## Características

### ✅ Emails Automáticos
- **Confirmación de Reserva**: Se envía automáticamente cuando se crea una nueva reserva
- **Cancelación de Reserva**: Se envía automáticamente cuando se cancela una reserva

### ✅ Diseño Responsivo
- Emails con diseño HTML profesional
- Compatible con la mayoría de clientes de email
- Colores y branding de Trattoria La Bella

### ✅ Envío Asíncrono
- Los emails se envían en background (goroutine)
- No bloquea la respuesta de la API
- Logging detallado de éxito/error

### ✅ Seguridad TLS/STARTTLS
- Soporte completo para STARTTLS (puerto 587) - **RECOMENDADO**
- Soporte para TLS directo (puerto 465)
- Verificación de certificados SSL
- Configuración flexible por variables de entorno

## Configuración

### 1. Variables de Entorno

Agrega las siguientes variables a tu archivo `.env`:

```bash
# Configuración SMTP con STARTTLS (RECOMENDADO)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=tu-email@gmail.com
SMTP_PASS=tu-contraseña-de-aplicación
FROM_EMAIL=reservas@trattorialabella.com
SMTP_USE_STARTTLS=true
SMTP_USE_TLS=false
```

### 2. Proveedores de Email Soportados

#### Gmail (STARTTLS - RECOMENDADO)
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=tu-email@gmail.com
SMTP_PASS=tu-contraseña-de-aplicación
SMTP_USE_STARTTLS=true
SMTP_USE_TLS=false
```

#### Gmail (TLS Directo)
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=465
SMTP_USER=tu-email@gmail.com
SMTP_PASS=tu-contraseña-de-aplicación
SMTP_USE_STARTTLS=false
SMTP_USE_TLS=true
```

**Importante**: Para Gmail, necesitas usar una "Contraseña de aplicación":
1. Ve a https://myaccount.google.com/apppasswords
2. Genera una contraseña de aplicación
3. Usa esa contraseña en `SMTP_PASS`

#### Outlook/Hotmail (STARTTLS)
```bash
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USER=tu-email@outlook.com
SMTP_PASS=tu-contraseña
SMTP_USE_STARTTLS=true
SMTP_USE_TLS=false
```

#### SendGrid (STARTTLS)
```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=tu-api-key-de-sendgrid
SMTP_USE_STARTTLS=true
SMTP_USE_TLS=false
```

## Arquitectura

### Estructura de Archivos
```
internal/
├── domain/
│   └── ports.go                    # Interfaz IEmailService
├── infra/
│   └── secundary/
│       └── email/
│           └── email_service.go    # Implementación del servicio
└── app/
    └── usecasereserve/
        ├── constructor.go          # Inyección de dependencias
        ├── create-reserve.go       # Envío de confirmación
        └── cancel-reservation.go   # Envío de cancelación
```

### Flujo de Email

1. **Creación de Reserva**:
   ```
   API → CreateReserve → Crear Reserva → Enviar Email (background)
   ```

2. **Cancelación de Reserva**:
   ```
   API → CancelReservation → Cancelar Reserva → Enviar Email (background)
   ```

## Templates de Email

### Confirmación de Reserva
- **Asunto**: "Confirmación de Reserva - Trattoria La Bella"
- **Contenido**: Detalles de la reserva, fecha, hora, número de personas
- **Diseño**: Verde/marrón con branding del restaurante

### Cancelación de Reserva
- **Asunto**: "Cancelación de Reserva - Trattoria La Bella"
- **Contenido**: Detalles de la reserva cancelada
- **Diseño**: Rojo con branding del restaurante

## Seguridad TLS/STARTTLS

### 🔒 Métodos de Cifrado Soportados

1. **STARTTLS (RECOMENDADO)**
   - Puerto: 587
   - Inicia conexión sin cifrado y luego la cifra
   - Compatible con la mayoría de servidores SMTP
   - Configuración: `SMTP_USE_STARTTLS=true`

2. **TLS Directo**
   - Puerto: 465
   - Conexión cifrada desde el inicio
   - Más rápido pero menos compatible
   - Configuración: `SMTP_USE_TLS=true`

3. **Sin Cifrado (NO RECOMENDADO)**
   - Solo para desarrollo/testing
   - No usar en producción
   - Configuración: Ambos en `false`

### 🛡️ Características de Seguridad
- ✅ Verificación de certificados SSL
- ✅ Fallback graceful si no hay configuración
- ✅ Logging detallado del método de seguridad usado
- ✅ Configuración flexible por entorno

## Logging

El sistema registra automáticamente:
- ✅ Emails enviados exitosamente (con método de seguridad)
- ❌ Errores de envío
- ⚠️ Configuración SMTP incompleta

### Ejemplos de Logs
```
INFO  Email enviado exitosamente email=cliente@ejemplo.com subject="Confirmación de Reserva - Trattoria La Bella" security=STARTTLS
ERROR Error enviando email error="authentication failed" to=cliente@ejemplo.com method=STARTTLS
WARN  Configuración SMTP incompleta, saltando envío de email
```

## Testing

### Configuración de Pruebas
Para pruebas locales, puedes usar servicios como:
- **Mailtrap**: Para capturar emails en desarrollo
- **Ethereal Email**: Para testing SMTP

### Ejemplo con Mailtrap
```bash
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=tu-usuario-mailtrap
SMTP_PASS=tu-contraseña-mailtrap
FROM_EMAIL=reservas@trattorialabella.com
```

## Seguridad

### ✅ Buenas Prácticas Implementadas
- Contraseñas de aplicación para Gmail
- Envío asíncrono para no bloquear la API
- Logging sin información sensible
- Validación de configuración SMTP

### ⚠️ Consideraciones
- Nunca commits las credenciales SMTP
- Usa variables de entorno en producción
- Considera usar servicios como SendGrid para producción

## Troubleshooting

### Problemas Comunes

1. **"authentication failed"**
   - Verifica las credenciales SMTP
   - Para Gmail, usa contraseña de aplicación

2. **"connection refused"**
   - Verifica `SMTP_HOST` y `SMTP_PORT`
   - Asegúrate de que el puerto esté abierto

3. **"error iniciando STARTTLS"**
   - Verifica que el servidor soporte STARTTLS
   - Intenta cambiar a TLS directo (`SMTP_USE_TLS=true`, puerto 465)
   - Verifica que el firewall permita el puerto

4. **"certificate verify failed"**
   - El servidor tiene un certificado SSL inválido
   - Verifica la configuración del servidor SMTP
   - Contacta al proveedor de email

5. **"Email no se envía"**
   - Revisa los logs para errores específicos
   - Verifica que todas las variables estén configuradas
   - Confirma la configuración de seguridad (STARTTLS/TLS)

### Debug
Para debug, revisa los logs del servidor:
```bash
docker logs central_reserve_prod
```

## Próximas Mejoras

- [ ] Templates personalizables por restaurante
- [ ] Notificaciones push como alternativa
- [ ] Sistema de reintentos para emails fallidos
- [ ] Métricas de entrega de emails
- [ ] Soporte para múltiples idiomas 