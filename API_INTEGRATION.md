# Guía de Integración con IAM APU

## API Gateway

La API está expuesta a través de Kong Gateway en `http://localhost:8000`. Todas las peticiones deben realizarse a esta URL base.

### Características del Gateway

1. **CORS**: Configurado para permitir:
   - Todos los orígenes (`*`)
   - Métodos: GET, POST, PUT, DELETE, OPTIONS
   - Headers: Accept, Authorization, Content-Type, X-Refresh-Token
   - Credentials: true
   - Max Age: 3600

2. **Rate Limiting**:
   - 60 peticiones por minuto
   - 1000 peticiones por hora

3. **Rutas**:
   - Base path: `/iam/api/v1`
   - No se hace strip_path, usar las rutas completas

## Autenticación

La API utiliza autenticación JWT. Para acceder a los endpoints protegidos:

1. Obtén el token mediante el endpoint `/auth/login`
2. Incluye el token en el header de las peticiones:
   ```
   Authorization: Bearer <tu_token>
   ```

### Flujo de Autenticación

1. **Login Local**:
   ```typescript
   const loginData = {
     email: "usuario@ejemplo.com",
     password: "contraseña",
     provider: "LOCAL",
     tenant_id: "uuid-del-tenant"
   };
   
   const response = await fetch('http://localhost:8000/iam/api/v1/auth/login', {
     method: 'POST',
     headers: { 'Content-Type': 'application/json' },
     body: JSON.stringify(loginData)
   });
   ```

2. **Login con Google**:
   ```typescript
   const loginData = {
     email: "usuario@ejemplo.com",
     google_token: "token-de-google",
     provider: "GOOGLE",
     tenant_id: "uuid-del-tenant"
   };
   ```

3. **Refrescar Token**:
   ```typescript
   const response = await fetch('/iam/api/v1/auth/refresh', {
     method: 'POST',
     headers: { 'Content-Type': 'application/json' },
     body: JSON.stringify({ refresh_token: "tu-refresh-token" })
   });
   ```

## Endpoints Públicos

Los siguientes endpoints no requieren autenticación:

- `GET /plans` - Listar todos los planes
- `POST /tenants` - Crear nuevo tenant
- `GET /tenants/by-email-key` - Buscar tenant por email key

## Manejo de Errores

La API retorna errores en el siguiente formato:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Descripción del error"
  }
}
```

## Ejemplo de Cliente API

```typescript
class ApiClient {
  private baseUrl = 'http://localhost:8000/iam/api/v1';
  private token: string | null = null;

  setToken(token: string) {
    this.token = token;
  }

  private async fetch(endpoint: string, options: RequestInit = {}) {
    const headers = {
      'Content-Type': 'application/json',
      ...(this.token ? { Authorization: `Bearer ${this.token}` } : {}),
      ...options.headers
    };

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    return response.json();
  }

  // Auth
  async login(data: LoginRequest) {
    return this.fetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  // Plans
  async getPlans() {
    return this.fetch('/plans');
  }

  // Tenants
  async createTenant(data: TenantCreate) {
    return this.fetch('/tenants', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  async getTenantByEmailKey(emailKey: string) {
    return this.fetch(`/tenants/by-email-key?email_key=${emailKey}`);
  }

  // Users
  async createUser(data: UserCreate) {
    return this.fetch('/users', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  async getUsersByTenant(tenantId: string) {
    return this.fetch(`/users/by-tenant/${tenantId}`);
  }
}
```

## Consideraciones Importantes

1. **Manejo de Tokens**:
   - Almacena el access_token y refresh_token de forma segura
   - Implementa un interceptor para refrescar el token automáticamente
   - Borra los tokens al hacer logout

2. **CORS**:
   - La API está configurada para aceptar peticiones del frontend
   - No deberías tener problemas de CORS si el frontend está en el dominio correcto

3. **Validaciones**:
   - Valida los datos antes de enviarlos a la API
   - Maneja los errores de forma adecuada en el frontend

4. **Tipos TypeScript**:
   - Usa los tipos definidos en el swagger.yaml para tipar tus requests/responses
   - Puedes usar herramientas como openapi-typescript-codegen para generar los tipos automáticamente

Para más detalles técnicos, consulta la especificación OpenAPI completa en `swagger.yaml`.
