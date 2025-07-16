# Scripts de Validación y Limpieza

Este directorio contiene scripts para mantener el repositorio limpio y evitar archivos grandes.

## 📋 Scripts Disponibles

### 🔍 `pre-commit-check.sh`
Valida que no se agreguen archivos grandes o problemáticos al repositorio.

**Uso:**
```bash
./scripts/pre-commit-check.sh
```

**Qué verifica:**
- Archivos mayores a 10MB
- Archivos de cache de BuildKit
- Archivos binarios ejecutables

### 🧹 `clean-cache.sh`
Limpia caches y archivos temporales del proyecto.

**Uso:**
```bash
./scripts/clean-cache.sh
```

**Qué limpia:**
- Cache de Docker BuildKit
- Carpetas `.buildkit/` locales
- Binarios de Go (`main`, `bin/`)
- Archivos temporales (`.tmp`, `.log`, `.cache`)
- `node_modules` (opcional)

## 🔧 Configuración Automática

### Hook de Pre-commit
Se ha configurado un hook de Git que ejecuta automáticamente `pre-commit-check.sh` antes de cada commit.

**Ubicación:** `.git/hooks/pre-commit`

**Funcionamiento:**
- Se ejecuta automáticamente al hacer `git commit`
- Si encuentra archivos problemáticos, cancela el commit
- Muestra mensajes informativos sobre cómo solucionar problemas

## 🚨 Problemas Comunes

### Error: "Se encontraron archivos grandes"
**Solución:**
1. Verifica que el archivo esté en `.gitignore`
2. Si es necesario, considera usar Git LFS
3. Si es temporal, elimínalo antes del commit

### Error: "Se encontraron archivos de cache de BuildKit"
**Solución:**
```bash
# Limpiar cache
./scripts/clean-cache.sh

# Remover del control de versiones
git rm --cached -r infra/.buildkit/ infra/scripts/.buildkit/
```

### Error: "Se encontraron archivos binarios"
**Solución:**
1. Verifica que el binario sea necesario
2. Si es generado automáticamente, agrégalo al `.gitignore`
3. Si es necesario, considera usar Git LFS

## 📝 Buenas Prácticas

1. **Ejecuta limpieza regularmente:**
   ```bash
   ./scripts/clean-cache.sh
   ```

2. **Verifica antes de commits importantes:**
   ```bash
   ./scripts/pre-commit-check.sh
   ```

3. **Revisa el `.gitignore`** si agregas nuevos tipos de archivos

4. **Usa Git LFS** para archivos grandes necesarios (>50MB)

## 🔄 Mantenimiento

### Actualizar el hook de pre-commit
Si modificas `pre-commit-check.sh`, el hook se actualiza automáticamente.

### Agregar nuevas validaciones
Edita `pre-commit-check.sh` para agregar nuevas reglas de validación.

### Compartir con el equipo
Los hooks de Git no se comparten automáticamente. Para compartir:
1. Copia `.git/hooks/pre-commit` a `scripts/pre-commit-hook`
2. Agrega instrucciones en este README para que el equipo lo instale manualmente 