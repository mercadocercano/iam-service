# Multi-Tenants IAM


## build app
```bash
go build -o iam src/main.go
```

## build image
```bash
docker-compose build --no-cache api
```


## Ejecutar Migraciones Manualmente

Si prefieres ejecutar las migraciones manualmente en lugar de usar el contenedor de migraciones, puedes seguir estos pasos:

1. Asegúrate de que el contenedor de PostgreSQL esté corriendo:
```bash
docker compose up postgres -d
```

2. Ejecuta cada migración en orden usando los siguientes comandos:
```bash
# Crear tabla de planes
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000001_create_plans_table.up.sql

# Crear tabla de roles
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000002_create_roles_table.up.sql

# Crear tabla de tenants
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000003_create_tenants_table.up.sql

# Crear tabla de usuarios
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000004_create_users_table.up.sql

# Agregar tipo saas_type
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000005_add_all_saas_type.up.sql

# Agregar soporte de autenticación
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000006_add_auth_support.up.sql
```

Para revertir las migraciones, ejecuta los scripts down en orden inverso:
```bash
# Eliminar soporte de autenticación
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000006_add_auth_support.down.sql

# Eliminar tipo saas_type
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000005_add_all_saas_type.down.sql

# Eliminar tabla de usuarios
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000004_create_users_table.down.sql

# Eliminar tabla de tenants
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000003_create_tenants_table.down.sql

# Eliminar tabla de roles
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000002_create_roles_table.down.sql

# Eliminar tabla de planes
docker exec -i iam-db psql -U postgres -d multi_tenants_iam < src/infrastructure/persistence/migrations/000001_create_plans_table.down.sql
