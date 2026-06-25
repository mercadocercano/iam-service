# ADR-002: S2S Authentication with Scoped API Keys

**Estado**: Propuesto  
**Fecha**: 2026-06-24  
**Deciders**: @dev-architect, @dev-security  

## Contexto

El servicio `iam-service` usaba una única `S2S_API_KEY` global compartida por todos los servicios internos (onboarding, sales, pim). Esa key otorgaba bypass total de autorización: cualquier servicio con la key podía llamar cualquier endpoint de gestión del IAM, saltándose rol y namespace. Esto viola least-privilege y crea riesgo de corrupción cross-tenant si un servicio comprometido o mal configurado accede a datos de otro tenant.

Caso motivador: `whatsapp-agent` necesita crear tenants y usuarios owner durante el wizard de onboarding, pero NO debe poder acceder al resto de los endpoints admin del IAM.

## Decisión

Reemplazar la god-key por un **registro de credenciales S2S por servicio**, cada una con un **scope acotado**. El secreto (key) se separa de la política (scopes):

- **Secreto**: una env var por servicio (`S2S_KEY_<SERVICE>`), inyectada desde GitHub Actions Secrets / sealed-secrets. La app solo hace `os.Getenv`.
- **Política**: mapa `servicio → scopes` en código (`s2s.ServicePolicy`), no secreto.
- **Enforcement**: cada grupo de rutas declara el scope requerido. El middleware `Authorize` resuelve `X-API-Key` contra el registro en memoria y permite solo si el scope requerido pertenece a la credencial.

## Alternativas consideradas

### Opción A — Scope por servicio (elegida)

- **Pros**: least-privilege real, rotación independiente por servicio, fail-closed, auditable por `s2s_service`.
- **Cons**: requiere configurar N env vars en vez de 1; consumidores legacy necesitan migración transicional.

### Opción B — Mantener god-key + path-based ACL

- **Pros**: un solo secret, cambio mínimo.
- **Cons**: no separa secreto de política; path-based ACL es más frágil y menos expresivo que scopes; no resuelve el caso de whatsapp-agent.

### Opción C — JWT firmados por servicio (mTLS o SPIFFE)

- **Pros**: estándar de largo plazo, identidad criptográfica fuerte.
- **Cons**: requiere infra adicional (CA interna, sidecars, rotación de certs). Se deja como evolución futura; el scope-based es un paso intermedio compatible.

## Consecuencias

- ✅ Eliminamos el bypass total de S2S.
- ✅ whatsapp-agent puede operar con `tenant:provision` sin acceso a users/roles/tenants/plans admin.
- ✅ onboarding/sales/pim siguen funcionando con `system:admin` (compatibilidad legacy).
- ⚠️ Los consumidores legacy tienen scopes amplios; se requiere épica de hardening para reducirlos.
- ⚠️ DevOps debe inyectar 4 variables de entorno en producción antes del deploy.
- 🚤 No se implementa rate-limiting específico S2S ni rotación automática en este cambio.

## Funcionamiento

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Request                                        │
│                    X-API-Key: <secret>                                      │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Authorize Middleware                                │
│  1. Busca X-API-Key en S2SRegistry (constant-time compare)                  │
│  2. Si matchea, verifica scope requerido por la ruta                        │
│     • scope OK  → set context s2s=true, s2s_service=<name> → NEXT            │
│     • scope NOK → 403 forbidden: insufficient scope                        │
│  3. Si no matchea, continúa a flujo JWT+rol (humano)                         │
│     • JWT válido + namespace OK + rol permitido → NEXT                     │
│     • cualquier otra cosa → 401/403                                          │
└─────────────────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Grupos de rutas                                   │
│  adminGroup        requires system:admin   → tenants, plans                 │
│  tenantScopedGroup requires system:admin   → users, roles (legacy)          │
│                    or tenant:admin                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                         S2SRegistry (boot)                                  │
│  ServicePolicy (código)              Env vars (secrets)                      │
│  whatsapp-agent → tenant:provision  ×  S2S_KEY_WHATSAPP_AGENT = ***        │
│  onboarding     → system:admin      ×  S2S_KEY_ONBOARDING     = ***        │
│  sales          → system:admin      ×  S2S_KEY_SALES          = ***        │
│  pim            → system:admin      ×  S2S_KEY_PIM              = ***        │
│                                      ↓                                      │
│                              map[key] → {service, scopes}                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Anti-corruption cross-tenant / cross-project

- El middleware NO introduce acceso cross-tenant: un servicio S2S con `tenant:provision` solo puede crear un nuevo tenant; no puede listar ni modificar tenants existentes.
- Los endpoints tenant-scoped (`users`, `roles`) requieren `tenant:admin` o `system:admin`. Un servicio con `tenant:provision` no puede acceder a ellos.
- El namespace check del flujo JWT humano permanece intacto.
- La capa de repositorio mantiene su filtro `tenant_id` (RULE-04). El gate solo decide si se permite la llamada; el aislamiento de datos sigue siendo responsabilidad de la capa de persistencia.

## Scopes definidos

| Scope            | Descripción                                              |
|------------------|----------------------------------------------------------|
| `tenant:provision` | Crear tenant + usuario owner inicial.                    |
| `tenant:admin`     | Operaciones tenant-scoped (users, roles).              |
| `system:admin`     | Acceso equivalente al god-key legacy (full admin IAM).   |

## Asignación por servicio

| Servicio        | Scope            | Motivo                                                                 |
|-----------------|------------------|------------------------------------------------------------------------|
| whatsapp-agent  | `tenant:provision` | Nuevo caso: wizard de onboarding, solo crea tenant + owner.           |
| onboarding      | `system:admin`    | Legacy: compatibilidad con god-key anterior.                            |
| sales           | `system:admin`    | Legacy.                                                                |
| pim             | `system:admin`    | Legacy.                                                                |

## Revisión prevista

Revisar en 90 días o al inicio de la épica de hardening:
- Reducir scopes de onboarding/sales/pim a lo mínimo real.
- Evaluar migración a JWT firmados por servicio o mTLS.
- Revisar métricas de `s2s_auth_denied_total` si se implementaron.
