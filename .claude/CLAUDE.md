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

## Memoria persistente (Engram)

Tenés acceso a memoria persistente entre sesiones vía las herramientas MCP de Engram (`mem_save`, `mem_search`, `mem_context`, etc.). Proyecto: **`mercado-cercano`** (memoria compartida con el resto del ecosistema).

**Cuándo guardar** — sin esperar que te lo pidan:
- Al resolver un bug no trivial: síntoma, causa raíz, fix aplicado.
- Al tomar una decisión de diseño: qué se decidió y por qué.
- Al descubrir un patrón o convención del proyecto que no está documentada.
- Al completar una feature o refactor significativo: qué cambió y dónde.

**Cuándo buscar** — antes de empezar cualquier tarea:
- `mem_context` al inicio de sesión o tras una compaction para recuperar el estado anterior.
- `mem_search` cuando el usuario menciona algo que puede tener historial ("el bug de autenticación", "la migración de la semana pasada").

**Al cerrar sesión**: llamar `mem_session_summary` para dejar un resumen recuperable.

## Convención de referencias — alias + título (SIEMPRE)

Al mencionar cualquier ítem por su alias —épicas (`E0N`), hitos (`H0N`),
propuestas (`PROP-00N` / `p00N`), specs (`S0N`), reglas (`RULE-0N`),
principios (`P-0N`)— incluí SIEMPRE su título/objetivo entre paréntesis la
primera vez que aparece en una respuesta. El alias suelto no se entiende.

- ❌ "Avancemos con E04 y después H1."
- ✅ "Avancemos con E04 (Team AI Harness Integration) y después H1 (Dashboard de Métricas)."
- ✅ "PROP-001 (ingesta multi-origen Hermes) está aprobada."

Si no conocés el título del alias, leé la fuente correspondiente
(`management/roadmap/roadmap.yaml`, `epicas/`, `propuestas/`, `PROJECT.md`)
ANTES de referenciarlo. Nunca uses el alias suelto.
