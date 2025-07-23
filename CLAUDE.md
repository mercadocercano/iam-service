# CLAUDE.md - IAM Service

Este archivo proporciona guía a Claude Code para trabajar con el servicio de Identity & Access Management (IAM) del sistema SaaS multi-tenant.

## Descripción del Servicio

El servicio IAM es el núcleo de autenticación y autorización del ecosistema, manejando:
- Autenticación de usuarios (JWT + OAuth)
- Gestión de tenants (multi-tenancy)
- Roles y permisos
- Planes de suscripción
- Control de acceso a recursos

## Arquitectura

### Estructura Hexagonal Modular
```
src/
├── main.go              # Punto de entrada
├── shared/              # Código compartido entre módulos
├── api/                 # Configuración API común
├── auth/                # Módulo de autenticación
├── user/                # Módulo de usuarios
├── tenant/              # Módulo de tenants
├── role/                # Módulo de roles
└── plan/                # Módulo de planes
```

### Patrón por Módulo
Cada módulo sigue la arquitectura hexagonal:
```
módulo/
├── domain/
│   ├── entity/         # Entidades de dominio
│   ├── valueobject/    # Value objects
│   ├── exception/      # Excepciones de dominio
│   └── port/           # Interfaces (puertos)
├── application/
│   ├── usecase/        # Casos de uso
│   ├── dto/            # DTOs request/response
│   └── mapper/         # Mappers DTO ↔ Domain
└── infrastructure/
    ├── http/           # Controladores REST
    ├── persistence/    # Repositorios PostgreSQL
    └── config/         # Configuración del módulo
```

## Comandos de Desarrollo

### Local
```bash
# Ejecutar servicio
go run src/main.go

# Tests con cobertura
go test -v -cover ./test/...

# Test específico de módulo
go test ./test/user/application/usecase/
```

### Docker (Recomendado)
```bash
# Desde el directorio raíz del proyecto
make dev-start    # Inicia todos los servicios backend
make dev-logs     # Ver logs del servicio IAM
make dev-restart  # Reiniciar servicios
```

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

### Autenticación (`/auth`)
```http
POST   /auth/login      # Login (email/password o Google OAuth)
POST   /auth/refresh    # Refrescar access token
GET    /auth/validate   # Validar token activo
POST   /auth/logout     # Cerrar sesión
```

### Usuarios (`/users`)
```http
POST   /users           # Crear usuario
GET    /users/:id       # Obtener usuario
PUT    /users/:id       # Actualizar usuario
DELETE /users/:id       # Eliminar usuario
GET    /users           # Listar usuarios (paginado)
```

### Tenants (`/tenants`)
```http
POST   /tenants                  # Crear tenant
GET    /tenants/:id              # Obtener tenant
GET    /tenants/by-slug/:slug    # Obtener por slug
PUT    /tenants/:id              # Actualizar tenant
DELETE /tenants/:id              # Eliminar tenant
GET    /tenants                  # Listar tenants
POST   /tenants/:id/plan         # Asignar plan
GET    /tenants/:id/features     # Obtener features
PUT    /tenants/:id/features     # Actualizar features
```

### Roles (`/roles`)
```http
POST   /roles           # Crear rol
GET    /roles/:id       # Obtener rol
PUT    /roles/:id       # Actualizar rol
DELETE /roles/:id       # Eliminar rol
GET    /roles           # Listar roles
```

### Planes (`/plans`)
```http
POST   /plans           # Crear plan (admin)
GET    /plans/:id       # Obtener plan
GET    /plans           # Listar planes (público)
```

## Headers Requeridos

### Autenticación
```http
Authorization: Bearer <jwt_token>
```

### Multi-tenancy
```http
X-Tenant-ID: <tenant_uuid>
```

### Content Type
```http
Content-Type: application/json
```

## Patrones de Respuesta

### Listado con Paginación
```json
{
  "items": [...],
  "total_count": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```

### Parámetros de Query
```
?page=1&page_size=10&sort_by=created_at&sort_dir=desc
```

## Testing

### Estructura de Tests
```
test/
├── user/               # Tests del módulo user
├── tenant/             # Tests del módulo tenant
├── auth/               # Tests del módulo auth
├── role/               # Tests del módulo role
├── plan/               # Tests del módulo plan
└── shared/             # Utilidades compartidas de test
```

### Patrón Object Mother
```go
// test/user/domain/user_mother.go
func UserMother() *entity.User {
    return UserMotherWithTenant(uuid.New())
}
```

### Ejecutar Tests
```bash
# Todos los tests
go test ./test/...

# Tests de un módulo
go test ./test/user/...

# Test específico con verbose
go test -v ./test/auth/application/usecase/
```

## Base de Datos

### Migraciones
```bash
# Las migraciones se ejecutan automáticamente al iniciar
# Archivos en: migrations/
```

### Esquema Principal
- `users` - Usuarios del sistema
- `tenants` - Organizaciones/empresas
- `roles` - Roles de usuario
- `plans` - Planes de suscripción
- `user_roles` - Relación usuario-rol
- `tenant_plans` - Relación tenant-plan
- `plan_features` - Features por plan

## Configuración

### Variables de Entorno
```bash
# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=iam_db
DB_SSLMODE=disable

# Servidor
PORT=8080

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=7d

# Métricas
PROMETHEUS_ENABLED=true
```

## Consideraciones de Desarrollo

### Multi-tenancy
- Todos los recursos están aislados por `tenant_id`
- El middleware valida automáticamente el header `X-Tenant-ID`
- Los repositorios filtran por tenant automáticamente

### Autenticación JWT
- Access tokens: 15 minutos (renovables)
- Refresh tokens: 7 días
- Los tokens incluyen: user_id, tenant_id, roles

### Seguridad
- Contraseñas hasheadas con bcrypt (cost: 10)
- Validación de permisos por rol
- Soporte para OAuth2 (Google)

### Casos de Uso
Cada operación se implementa como un caso de uso:
```go
// Ejemplo: CreateUserUseCase
type CreateUserUseCase struct {
    userRepo     port.UserRepositoryPort
    roleRepo     port.RoleRepositoryPort
    eventBus     port.EventBusPort
}
```

### Manejo de Errores
- Excepciones de dominio específicas
- Códigos HTTP apropiados
- Mensajes de error consistentes

## Flujo de Desarrollo Recomendado

1. **Entender el módulo**: Revisar domain/entity y domain/port
2. **Implementar caso de uso**: En application/usecase
3. **Crear DTOs**: En application/dto si es necesario
4. **Implementar controller**: En infrastructure/http
5. **Escribir tests**: Usando Object Mother pattern
6. **Actualizar rutas**: En el router del módulo

## Integración con otros Servicios

El servicio IAM es consultado por:
- **API Gateway (Kong)**: Validación de tokens
- **Otros microservicios**: Verificación de permisos
- **Frontends**: Autenticación y autorización

## Monitoreo

### Health Check
```http
GET /health
```

### Métricas Prometheus
```http
GET /metrics
```

### Logs
- Formato estructurado con Gin
- Niveles: INFO, WARN, ERROR
- Incluyen: timestamp, method, path, status, duration

## Tips de Desarrollo

1. **Siempre validar tenant_id** en operaciones multi-tenant
2. **Usar transacciones** para operaciones múltiples
3. **Cachear roles y permisos** cuando sea posible
4. **Documentar casos de uso complejos** en el código
5. **Mantener tests actualizados** con Object Mother

## Comandos Útiles

```bash
# Verificar que el servicio responde
curl http://localhost:8080/health

# Obtener planes disponibles (público)
curl http://localhost:8080/api/v1/plans

# Login de prueba
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'

# Listar usuarios (requiere token)
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <token>" \
  -H "X-Tenant-ID: <tenant-uuid>"
```