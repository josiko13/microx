# MicroX 🐦 - Plataforma de Microblogging

Una plataforma de microblogging similar a Twitter construida con **Go**, optimizada para escalabilidad y rendimiento.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-Web%20Framework-green.svg)](https://gin-gonic.com/)
[![MySQL](https://img.shields.io/badge/MySQL-Database-orange.svg)](https://www.mysql.com/)
[![Redis](https://img.shields.io/badge/Redis-Cache-red.svg)](https://redis.io/)

## Características

- ✅ **Tweets**: Publicar mensajes cortos (máximo 280 caracteres)
- ✅ **Follow**: Seguir a otros usuarios
- ✅ **Timeline**: Ver tweets de usuarios seguidos
- 🚀 **Escalable**: Diseñado para millones de usuarios
- ⚡ **Optimizado para lecturas**: Cache distribuido y índices optimizados

## Arquitectura

### Tecnologías
- **Backend**: Go 1.21+ con Gin framework
- **Base de datos**: MySQL (datos persistentes) + Redis (cache)
- **Autenticación**: Header X-User-ID (autenticación simple)
- **Logging**: Librería estándar `log` de Go

### Estructura del Proyecto
```
microx/
├── cmd/                      # Puntos de entrada
│   ├── server/              # Servidor principal
│   └── migrate/             # Script de migraciones
├── internal/                 # Código interno de la aplicación
│   ├── api/                 # Handlers HTTP
│   ├── service/             # Lógica de negocio
│   ├── repository/          # Acceso a datos
│   │   ├── mysql/          # Repositorios MySQL
│   │   └── redis/          # Repositorios Redis
│   ├── model/               # Modelos de dominio
│   ├── middleware/          # Middleware
│   └── config/              # Configuraciones
├── migrations/              # Migraciones de BD
├── Dockerfile               # Containerización
├── docker-compose.yml       # Stack completo
└── config.env.example       # Template de configuración
```

## 🚀 Instalación

### Opción 1: Instalación Manual

1. **Clonar el repositorio**
```bash
git clone <repo-url>
cd microx
```

2. **Instalar dependencias**
```bash
go mod download
```

3. **Configurar variables de entorno**
```bash
cp config.env.example config.env
# Editar config.env con tus configuraciones
```

4. **Levantar bases de datos**
```bash
# MySQL
docker run -d --name microx-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql:8.0

# Redis
docker run -d --name microx-redis -p 6379:6379 redis:7-alpine
```

5. **Ejecutar migraciones**
```bash
go run cmd/migrate/main.go
```

6. **Iniciar el servidor**
```bash
go run cmd/server/main.go
```

### Opción 2: Usando Docker Compose (Recomendado)

```bash
# Levantar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down
```



## API Endpoints

### Autenticación
La API usa autenticación simple mediante el header `X-User-ID`. Para rutas protegidas, incluye este header con el ID del usuario.

**Ejemplo:**
```bash
curl -H "X-User-ID: 1" http://localhost:8080/api/tweets
```

### Tweets
- `POST /api/tweets` - Crear un tweet (requiere X-User-ID)
- `GET /api/tweets/:id` - Obtener un tweet específico (requiere X-User-ID)
- `GET /api/users/:id/tweets` - Obtener tweets de un usuario

### Follow
- `POST /api/follow/:user_id` - Seguir a un usuario (requiere X-User-ID)
- `DELETE /api/follow/:user_id` - Dejar de seguir a un usuario (requiere X-User-ID)
- `GET /api/users/:id/followers` - Obtener seguidores
- `GET /api/users/:id/following` - Obtener usuarios seguidos

### Timeline
- `GET /api/timeline` - Obtener timeline personal (requiere X-User-ID)
- `POST /api/timeline/refresh` - Refrescar timeline (requiere X-User-ID)

### Usuarios
- `POST /api/users` - Crear un usuario
- `GET /api/users/:id` - Obtener información de usuario
- `GET /api/users/:id/stats` - Obtener estadísticas de usuario

## Optimizaciones para Escalabilidad

1. **Cache Distribuido**: Redis para timeline y datos frecuentemente accedidos
2. **Índices Optimizados**: Índices compuestos en MySQL para consultas de timeline
3. **Paginación**: Implementación eficiente de paginación con limit/offset
4. **Connection Pooling**: Pool de conexiones para MySQL y Redis

## 🛠️ Desarrollo

### Comandos de Desarrollo

```bash
# Ejecutar en modo desarrollo
go run cmd/server/main.go

# Ejecutar tests
go test ./... -v

# Ejecutar tests con coverage
go test -cover -v ./...

# Build para producción
go build -o bin/server cmd/server/main.go

# Formatear código
go fmt ./...

# Ejecutar linter
golint -set_exit_status ./...
```

## 🐳 Docker

### Usando Docker Compose

```bash
# Levantar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down
```

### Servicios incluidos:
- **🚀 MicroX API**: http://localhost:8080
- **📊 Adminer** (MySQL): http://localhost:8081
- **🔴 Redis Commander**: http://localhost:8082
- **🐘 MySQL**: localhost:3306
- **🔴 Redis**: localhost:6379

### Construir imagen Docker

```bash
docker build -t microx:latest .
```

## 📋 Variables de Entorno

```env
# Servidor
PORT=8080
ENV=development

# Base de datos MySQL
DB_HOST=localhost
DB_PORT=3306
DB_NAME=microx
DB_USER=root
DB_PASSWORD=password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Configuración de la aplicación
MAX_TWEET_LENGTH=280
TIMELINE_CACHE_TTL=3600
```

## 🤝 Contribuir

¡Las contribuciones son bienvenidas! Si quieres contribuir al proyecto, por favor crea un issue o un pull request.

## 🎯 Asunciones del Proyecto

- Todos los usuarios son válidos (no se requiere registro/login)
- No se manejan sesiones complejas
- Autenticación mediante header X-User-ID
- Optimizado para lecturas (timeline, feeds)
- Diseñado para escalar a millones de usuarios 