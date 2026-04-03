# CLAUDE.md - IAM Service

Servicio de Identity & Access Management del ecosistema SaaS multi-tenant.
Autenticación JWT/OAuth2, gestión de tenants, roles, permisos y planes de suscripción.

**Puerto**: 8080 | **Stack**: Go + Gin + PostgreSQL | **Arquitectura**: Hexagonal modular

Hablame siempre en español.

## Comandos esenciales

```bash
go run src/main.go                    # Ejecutar (puerto 8080)
go test -v -cover ./test/...          # Tests con cobertura
make dev-start                        # Docker (todos los backends)
```

## Regla: MCP Go Generator primero

Antes de implementar código Go, consultar MCP:
```bash
analyzeUsecaseWorkflow --service_name="iam" --entity_name="user"
generateWorkflowRoadmap --service_name="iam"
generateComponentByStep --step_type="usecase" --entity_name="role"
```

## Contexto on-demand (cargar según necesidad)

| Archivo | Cuándo cargar |
|---------|---------------|
| `iam-service-management/api-endpoints.md` | Al trabajar con endpoints, rutas o curl |
| `iam-service-management/architecture.md` | Al crear módulos, entidades, casos de uso |
| `iam-service-management/config.md` | Al configurar env vars, Docker, MCP |

## Reglas compartidas (cargar según contexto)

| Regla | Archivo |
|-------|---------|
| Arquitectura hexagonal | `ai-tools/rules/architecture.md` |
| Multi-tenancy | `ai-tools/rules/multi-tenant.md` |
| API Gateway / Kong | `ai-tools/rules/api-gateway.md` |
| Testing standards | `ai-tools/rules/testing-standards.md` |
| Formato respuesta API | `ai-tools/rules/api-response-format.md` |
