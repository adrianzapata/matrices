# Matrices Coding Challenge — Interseguro

**Desplegado en:** https://go-api-7q4y.onrender.com (usuario `admin`, ver
credenciales entregadas por correo). API de estadísticas (uso interno de
`go-api`): https://node-api-srj3.onrender.com

> Nota: ambos servicios están en el plan gratuito de Render, que apaga la
> instancia tras ~15 min de inactividad. La primera petición después de eso
> puede tardar ~50s en responder mientras arranca de nuevo.

Dos APIs REST que se comunican por HTTP:

1. **`go-api`** (Go + Fiber): recibe una matriz, la rota 90° y calcula su
   factorización QR (reflexiones de Householder), y reenvía los resultados a
   la API de Node.js.
2. **`node-api`** (Node.js + Express): recibe las matrices calculadas por la
   API en Go y devuelve estadísticas (máximo, mínimo, promedio, suma total y
   si cada matriz es diagonal).

Incluye además un **frontend** estático opcional (HTML/JS puro) que consume
ambas APIs a través de la API en Go, y autenticación **JWT** opcional en
ambas APIs.

## Por qué "rotación" y "QR" al mismo tiempo

El enunciado es inconsistente entre secciones:

- *"Arquitectura de la solución"* dice que la API en Go **rota** la matriz y
  envía el resultado a Node.js.
- *"Funcionalidad requerida"* dice que la API en Go debe **devolver la
  factorización QR** de la matriz.

Ambas son requisitos explícitos, así que en vez de descartar uno, la API en
Go hace las dos cosas sobre la matriz de entrada: la rota (satisface la
arquitectura descrita) y calcula su QR (satisface la funcionalidad
requerida). Reenvía a Node.js las tres matrices resultantes — `rotated`,
`q`, `r` — lo cual además encaja mejor con el texto de "Operación adicional"
("verificar si **alguna matriz** es diagonal", en plural), algo que un único
resultado de rotación no explicaría tan bien. Esta decisión se documenta
aquí para poder sustentarla en la entrevista, tal como pide el enunciado.

## Arquitectura

```
                    POST /api/auth/login
   Frontend ───────────────────────────────►  go-api
   (HTML/JS)  ◄─────────────────────────────  (Fiber, :8080)
                    POST /api/matrix/process
                                │
                                │ POST /api/stats
                                │ (JWT de servicio)
                                ▼
                            node-api
                           (Express, :4000)
```

- `go-api` valida la matriz de entrada, calcula `rotated` (rotación 90°
  horaria) y `q`/`r` (factorización QR vía reflexiones de Householder —
  numéricamente estable y válida para matrices rectangulares, no solo
  cuadradas), y llama a `node-api` para obtener las estadísticas.
- `node-api` recibe un conjunto de matrices nombradas y calcula máximo,
  mínimo, promedio, suma total, y si cada una es diagonal.
- La comunicación entre ambas usa HTTP + JSON, autenticada con un JWT de
  servicio que `go-api` genera al arrancar (firmado con el mismo
  `JWT_SECRET` que comparten ambos servicios).

## Requisitos previos

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker](https://www.docker.com/) + Docker Compose (para contenerización)

## Cómo correr todo con Docker (recomendado)

```bash
cp .env.example .env
# edita .env si quieres cambiar credenciales/secreto
docker compose up --build
```

- Frontend: http://localhost:8081
- API Go: http://localhost:8080
- API Node: http://localhost:4000 (normalmente solo la llama `go-api`
  internamente, pero queda expuesta para pruebas directas)

## Cómo correr cada servicio por separado (desarrollo local)

### API Go

```bash
cd go-api
go mod download
go run .
```

Variables de entorno (con sus defaults):

| Variable        | Default                    | Descripción                                   |
|-----------------|-----------------------------|------------------------------------------------|
| `PORT`          | `8080`                     | Puerto HTTP                                    |
| `NODE_API_URL`  | `http://localhost:4000`    | Base URL de la API de Node                     |
| `JWT_SECRET`    | `dev-secret-change-me`     | Secreto HS256 compartido con `node-api`        |
| `AUTH_USERNAME` | `admin`                    | Usuario demo para `/api/auth/login`            |
| `AUTH_PASSWORD` | `admin`                    | Password demo para `/api/auth/login`           |
| `CORS_ORIGINS`  | `*`                        | Orígenes permitidos por CORS                   |

### API Node.js

```bash
cd node-api
npm install
npm start
```

Variables de entorno:

| Variable        | Default                 | Descripción                            |
|-----------------|--------------------------|------------------------------------------|
| `PORT`          | `4000`                  | Puerto HTTP                              |
| `JWT_SECRET`    | `dev-secret-change-me`  | Secreto HS256 (debe coincidir con Go)    |
| `CORS_ORIGINS`  | `*`                     | Orígenes permitidos por CORS             |

### Frontend

Es HTML/JS estático; ábrelo directamente (`frontend/index.html`) o sírvelo
con cualquier servidor estático. Al cargarlo, indica la URL de `go-api` (por
defecto `http://localhost:8080`).

## Pruebas

```bash
# Go: unit tests de rotación, validación y QR (incluye verificación de que
# Q*R reconstruye la matriz original y que Q es ortogonal)
cd go-api && go test ./... -v

# Node.js: unit tests de estadísticas/validación + tests de integración
# del endpoint /api/stats (auth JWT incluida)
cd node-api && npm test
```

## Contratos de las APIs

### `POST /api/auth/login` (go-api, público)

```json
// request
{ "username": "admin", "password": "admin" }

// response 200
{ "token": "<jwt>", "expiresInSeconds": 3600 }
```

### `POST /api/matrix/process` (go-api, requiere `Authorization: Bearer <jwt>`)

```json
// request
{ "matrix": [[12, -51, 4], [6, 167, -68], [-4, 24, -41]] }

// response 200
{
  "original": [[12, -51, 4], ...],
  "rotated":  [[-4, 6, 12], ...],
  "q":        [[...]],
  "r":        [[...]],
  "stats": {
    "max": 167,
    "min": -68,
    "average": 12.9,
    "sum": 116.4,
    "diagonal": { "rotated": false, "q": false, "r": false }
  }
}
```

Errores: `400` (JSON malformado), `422` (matriz vacía o no rectangular),
`401` (token ausente/ inválido), `502` (falla al contactar `node-api`).

### `POST /api/stats` (node-api, requiere `Authorization: Bearer <jwt>`)

Pensada para ser llamada por `go-api`, pero puede probarse directamente:

```json
// request
{ "matrices": { "rotated": [[...]], "q": [[...]], "r": [[...]] } }

// response 200
{
  "max": 167,
  "min": -68,
  "average": 12.9,
  "sum": 116.4,
  "diagonal": { "rotated": false, "q": false, "r": false }
}
```

Errores: `400` (JSON malformado), `422` (matriz vacía/ no rectangular/ con
valores no numéricos), `401` (token ausente/inválido).

## Decisiones de diseño relevantes

- **QR vía Householder, no Gram-Schmidt**: funciona para matrices
  rectangulares (no solo cuadradas) sin el requisito de `filas >= columnas`
  que tiene Gram-Schmidt clásico, y es numéricamente más estable.
- **JWT compartido entre servicios**: en vez de exponer `node-api` sin
  protección, `go-api` firma un token de servicio al arrancar (con el mismo
  secreto) para autenticar sus llamadas a `node-api`. En producción, cada
  servicio debería tener su propio secreto/rol vía un mecanismo como OAuth2
  client credentials; aquí se simplificó dado el alcance del reto.
- **Estadísticas sobre múltiples matrices**: `node-api` calcula max/min/
  promedio/suma sobre el conjunto combinado de valores de todas las matrices
  recibidas, y la verificación de "diagonal" por separado para cada una.
- **Sin base de datos ni estado**: ambas APIs son *stateless*; no se
  requirió persistencia según el enunciado.

## Despliegue en la nube

Este repositorio queda listo para desplegarse en cualquier plataforma que
soporte contenedores Docker (Render, Railway, Fly.io, AWS ECS/App Runner,
Azure Container Apps, etc.). Pasos generales:

1. Sube el repo a GitHub.
2. En la plataforma elegida, crea un servicio por cada carpeta
   (`go-api`, `node-api`, `frontend`) usando su `Dockerfile`, o usa
   `docker-compose.yml` si la plataforma lo soporta directamente.
3. Configura las variables de entorno de la tabla de arriba como secretos
   (especialmente `JWT_SECRET` y `AUTH_PASSWORD`, nunca los valores de
   ejemplo).
4. Apunta `NODE_API_URL` en `go-api` a la URL pública/interna de `node-api`
   en esa plataforma.
5. Si el frontend se sirve desde otro dominio, ajusta `CORS_ORIGINS` en
   ambas APIs a esa URL en vez de `*`.
