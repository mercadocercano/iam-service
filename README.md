# saas-mt-iam-service

Este proyecto fue extraído del monorepo SaaS Marketplace como parte de la migración a repositorios independientes.

## Descripción

Servicio iam del ecosistema SaaS Marketplace.

## Tecnología

- **Tipo**: go
- **Lenguaje**: Go
- **Framework**: Gin/Fiber
- **Base de datos**: PostgreSQL

## Desarrollo

### Prerrequisitos

- Go 1.21+
- PostgreSQL 15+
- Docker (opcional)

### Instalación

```bash
# Clonar el repositorio
git clone https://github.com/trinityweb/saas-mt-iam-service.git
cd saas-mt-iam-service

# Instalar dependencias
go mod download

# Configurar variables de entorno
cp .env.example .env
# Ver "Autenticación S2S" más abajo para poblar las keys de dev.

# Ejecutar
go run src/main.go
```

### Docker

```bash
# Construir imagen
docker build -t saas-mt-iam-service .

# Ejecutar contenedor
docker run -p 8080:8080 saas-mt-iam-service
```

## Configuración

Copiá `.env.example` a `.env` y configurá las variables necesarias.

### Autenticación S2S (service-to-service)

El IAM ya no usa una única `S2S_API_KEY` global. Cada servicio consumidor tiene su propia key y un scope acotado:

| Servicio        | Env var                  | Scope            | Qué puede hacer hoy                           |
|-----------------|--------------------------|------------------|-----------------------------------------------|
| whatsapp-agent  | `S2S_KEY_WHATSAPP_AGENT` | `tenant:provision` | Crear tenant + usuario owner (wizard)         |
| onboarding      | `S2S_KEY_ONBOARDING`     | `system:admin`    | Todo lo que hacía la god-key legacy (migrar)  |
| sales           | `S2S_KEY_SALES`          | `system:admin`    | Igual legacy                                  |
| pim             | `S2S_KEY_PIM`            | `system:admin`    | Igual legacy                                  |

#### Generar keys de desarrollo

**Nunca comitear keys.** Para dev local generá valores dummy distintos por servicio:

```bash
openssl rand -hex 32
```

Repetí el comando una vez por servicio y pegá los resultados en tu `.env` personal:

```bash
# ejemplo (no uses estos valores en prod)
S2S_KEY_WHATSAPP_AGENT=2a3b...
S2S_KEY_ONBOARDING=4c5d...
S2S_KEY_SALES=6e7f...
S2S_KEY_PIM=8g9h...
```

En CI/CD los valores reales vienen de GitHub Actions Secrets. En k3s se inyectan vía sealed-secrets. La app solo lee `os.Getenv("S2S_KEY_<SERVICE>")`; no hace fetch a GitHub en runtime.

#### Rotación de keys

Cambiar la env var de un solo servicio y reiniciar el pod es suficiente. Las demás keys no se ven afectadas. No hay key en código ni en commits.

## Documentación

- [Documentación de API](./api-docs/)
- [Documentación técnica](./documentation/)

## Migración desde Monorepo

Este proyecto fue extraído del monorepo original manteniendo todo su historial de git.

|**Repositorio original**: https://github.com/trinityweb/saas-marketplace.git

## Contribución

1. Fork el proyecto
2. Crea una rama para tu feature
3. Commit tus cambios
4. Push a la rama
5. Abre un Pull Request

## Licencia

[Licencia del proyecto]
