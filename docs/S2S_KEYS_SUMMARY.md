# Resumen: S2S keys por servicio con scope

## Cambios realizados

- `src/auth/infrastructure/s2s/registry.go`: nuevo registro S2S que carga `S2S_KEY_<SERVICE>` del entorno y cruza con `ServicePolicy` (política en código, no secreta).
- `src/auth/infrastructure/middleware/authorize.go`: `Authorize` ahora requiere un scope S2S. Rama S2S hace lookup constant-time y permite solo si la credencial tiene el scope de la ruta. Key desconocida cae a JWT/401.
- `src/main.go`: reemplaza la única `S2S_API_KEY` por `s2s.LoadFromEnv()` + `ScopeMiddlewareFactory`. Todos los grupos de gestión requieren `system:admin` (el scope más amplio) por ahora; se puede cambiar por ruta cuando se agreguen scopes más finos.
- `src/auth/infrastructure/s2s/registry_test.go` y `authorize_test.go`: tests para lookup, constant-time, scope insuficiente, key desconocida, JWT intacto.
- `.env.example`: agrega `S2S_KEY_WHATSAPP_AGENT`, `S2S_KEY_ONBOARDING`, `S2S_KEY_SALES`, `S2S_KEY_PIM` (sin valores).
- `README.md`: documentación de S2S, tabla de scopes, comando `openssl rand -hex 32` para dev, política de rotación.
- `.github/workflows/deploy.yml`: comentarios sobre dónde deben vivir las S2S keys y TODO para automatizar el apply del Secret/sealed-secret.

## Scopes definidos

| Scope            | Descripción                                              |
|------------------|----------------------------------------------------------|
| `tenant:provision` | Crear tenant + usuario owner inicial.                    |
| `system:admin`     | Acceso equivalente al god-key legacy (full admin IAM).   |

## Asignación por servicio

| Servicio        | Scope            | Motivo                                                                 |
|-----------------|------------------|------------------------------------------------------------------------|
| whatsapp-agent  | `tenant:provision` | Nuevo caso: wizard de onboarding, solo necesita crear tenant + owner. |
| onboarding      | `system:admin`    | Legacy: mantener compatibilidad con lo que hacía la god-key anterior.   |
| sales           | `system:admin`    | Idem legacy.                                                           |
| pim             | `system:admin`    | Idem legacy.                                                           |

## Migración elegida

Opción (a): los consumidores legacy siguen con `system:admin` (equivalente al god-key anterior). No se rompe ningún cliente existente. Se recomienda crear una épica de hardening para reducir los scopes de onboarding/sales/pim y luego retirar el `system:admin` legacy.

## Notas para @architect / @security

1. **Namespace check humano intacto**: no se modificó el flujo JWT ni el namespace check.
2. **Comparación constant-time**: se mantiene `subtle.ConstantTimeCompare` en `Registry.Lookup`.
3. **No se loguean keys**: el middleware setea `s2s_service` en el contexto para logs, nunca la key.
4. **Fail-closed**: env ausente, key desconocida o scope insuficiente → denegar.
5. **Rotación**: cambiar la env var del servicio y reiniciar el pod es suficiente; no afecta a los demás.
6. **CI/CD**: el workflow no inyecta secrets todavía. Hace falta que DevOps aplique el Secret/sealed-secret del deployment con las 4 nuevas variables antes del próximo deploy.
7. **Coverage**: `go test ./... -race -timeout 120s` y `go vet ./...` pasan.

## Pendiente de revisión

- Validar si `tenantScopedGroup` (users, roles) debería seguir con `system:admin` o cambiar a un scope más fino como `tenant:admin` para S2S, dejando `system:admin` solo para `adminGroup`.
- Decidir si se quiere opción (b) de migración (mover cada consumidor a su key propia en el mismo cambio) en lugar de (a).
