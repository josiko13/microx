# MicroX - Arquitectura & Componentes

## Visión General

MicroX es una plataforma de microblogging desarrollada en Go que implementa una arquitectura Clean con separación clara de responsabilidades y tecnologías modernas para optimizar el rendimiento de la misma.

## Arquitectura de Alto Nivel

```
┌────────────────────────────────────────────────────────────┐
│                        API Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Handlers  │  │ Middleware  │  │   Routes    │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────┐
│                    Business Logic Layer                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │User Service │  │Tweet Service│  │Follow Service│        │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Data Access Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   MySQL     │  │    Redis    │  │ Repositories│          │
│  │  (Primary)  │  │   (Cache)   │  │             │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

## Componentes Principales

### 1. API Layer (Capa de Presentación)

#### Handlers
- **UserHandler**: Gestión de usuarios y estadísticas
- **TweetHandler**: Operaciones CRUD de tweets
- **FollowHandler**: Gestión de relaciones de seguimiento
- **TimelineHandler**: Obtención de timelines personalizados

#### Middleware
- **AuthMiddleware**: Autenticación simple basada en header `X-User-ID`
- **CORS**: Configuración de Cross-Origin Resource Sharing

#### Framework
- **Gin**: Framework web ligero y de alto rendimiento para Go

### 2. Business Logic Layer (Capa de Lógica de Negocio)

#### Services
- **UserService**: Lógica de negocio para usuarios
  - Creación de usuarios
  - Obtención de estadísticas (tweets, seguidores, seguidos)
  
- **TweetService**: Lógica de negocio para tweets
  - Creación de tweets con validaciones
  - Obtención de tweets individuales y por usuario
  - Integración automática con timeline de seguidores
  
- **FollowService**: Lógica de negocio para relaciones de seguimiento
  - Seguir/dejar de seguir usuarios
  - Obtención de seguidores y usuarios seguidos
  - Actualización automática de timelines

- **TimelineService**: Lógica de negocio para timelines
  - Obtención de timeline personalizado
  - Refresh de timeline desde base de datos

### 3. Data Access Layer (Capa de Acceso a Datos)

#### Repositories
- **UserRepository**: Operaciones de base de datos para usuarios
- **TweetRepository**: Operaciones de base de datos para tweets
- **FollowRepository**: Operaciones de base de datos para follows
- **TimelineRepository**: Operaciones de caché para timelines

#### Base de Datos
- **MySQL**: Base de datos principal para datos persistentes
  - Tabla `users`: Información de usuarios
  - Tabla `tweets`: Contenido de tweets
  - Tabla `follows`: Relaciones de seguimiento

- **Redis**: Caché de alto rendimiento para timelines
  - Almacenamiento de timelines personalizados
  - Optimización de lecturas frecuentes
  - Invalidación automática de caché

## Patrones Arquitectónicos utilizados.

### 1. Clean Architecture
- **Separación de responsabilidades**: Cada capa tiene una responsabilidad específica
- **Inversión de dependencias**: Las capas internas no dependen de las externas
- **Interfaces**: Uso de interfaces para desacoplar implementaciones

### 2. Repository Pattern
- **Abstracción de datos**: Los servicios no conocen detalles de implementación
- **Testabilidad**: Fácil mockeo para tests unitarios
- **Flexibilidad**: Cambio de implementación sin afectar lógica de negocio

### 3. Service Layer Pattern
- **Lógica de negocio centralizada**: Reglas de negocio en servicios
- **Reutilización**: Servicios pueden ser usados por múltiples handlers
- **Transacciones**: Manejo de transacciones a nivel de servicio


## Tecnologías Utilizadas

### Backend
- **Go 1.21**: Lenguaje de programación principal
- **Gin**: Framework web para API REST
- **database/sql**: Acceso nativo a base de datos
- **go-redis**: Cliente Redis para Go

### Base de Datos
- **MySQL 8.0**: Base de datos relacional principal
- **Redis 7.0**: Caché en memoria para optimización

### Infraestructura
- **Docker**: Solo para ejecutar Redis en contenedor
- **MySQL**: Instalación local o en servidor

### Herramientas de Desarrollo
- **Postman**: Testing de API
- **Git**: Control de versiones
- **Logging**: Librería estándar `log` de Go

## Configuración y Despliegue

### Variables de Entorno
- `DB_HOST`: Host de MySQL (ej: localhost)
- `DB_PORT`: Puerto de MySQL (ej: 3306)
- `DB_USER`: Usuario de MySQL
- `DB_PASSWORD`: Contraseña de MySQL
- `DB_NAME`: Nombre de base de datos
- `REDIS_HOST`: Host de Redis (ej: localhost)
- `REDIS_PORT`: Puerto de Redis (ej: 6379)
- `REDIS_PASSWORD`: Contraseña de Redis (opcional)
- `SERVER_PORT`: Puerto del servidor (ej: 8080)

### Requisitos de Infraestructura
- **MySQL**: Base de datos instalada y configurada
- **Redis**: Ejecutándose en contenedor Docker o instalación local
- **Go**: Entorno de desarrollo Go 1.21+

### Estructura de Directorios
```
microx/
├── cmd/                    # Puntos de entrada de la aplicación
├── internal/              # Código interno de la aplicación
│   ├── api/              # Handlers y middleware
│   ├── config/           # Configuración de base de datos
│   ├── model/            # Modelos de datos
│   ├── repository/       # Capa de acceso a datos
│   └── service/          # Lógica de negocio
├── migrations/           # Scripts de migración de base de datos
├── docs/                # Documentación
└── config.env           # Variables de entorno
```

## Ventajas de la Arquitectura

### 1. Escalabilidad
- **Separación de responsabilidades**: Fácil escalado independiente de componentes
- **Caché distribuido**: Redis permite escalado horizontal
- **Base de datos optimizada**: MySQL con índices apropiados

### 2. Mantenibilidad
- **Código limpio**: Estructura clara y organizada
- **Interfaces bien definidas**: Fácil modificación sin afectar otras partes
- **Documentación completa**: Guías claras para desarrolladores

### 3. Rendimiento
- **Caché inteligente**: Redis para operaciones frecuentes
- **Consultas optimizadas**: JOINs y índices apropiados
- **Framework eficiente**: Gin para alto rendimiento

### 4. Testabilidad
- **Arquitectura limpia**: Fácil mockeo de dependencias
- **Interfaces**: Testing unitario simplificado
- **Separación de capas**: Testing independiente por capa

## Consideraciones Futuras para Producción

### 1. Seguridad
- **Validación de entrada**: Sanitización de datos de usuario
- **Autenticación**: Sistema simple pero extensible

### 2. Monitoreo
- **Logs estructurados**: Para debugging y análisis
- **Métricas**: Rendimiento y uso de recursos
- **Health checks**: Verificación de estado de servicios

### 3. Backup y Recuperación
- **Backup automático**: Base de datos MySQL
- **Persistencia Redis**: Configuración de persistencia