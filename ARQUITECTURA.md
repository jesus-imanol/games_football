# 🏛️ Arquitectura Limpia - Explicación Detallada

## Principios de Clean Architecture

```
┌─────────────────────────────────────────────────────┐
│                   MAIN.GO                           │
│              (Entry Point)                          │
└───────────────────┬─────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────┐
│            DEPENDENCIES                             │
│        (Inyección de Dependencias)                  │
│                                                     │
│  ┌───────────────────────────────────────────┐     │
│  │  Infrastructure → Application → Domain    │     │
│  └───────────────────────────────────────────┘     │
└───────────────────┬─────────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
        ▼                       ▼
┌──────────────┐         ┌──────────────┐
│ CONTROLLERS  │         │   ROUTERS    │
│ (Handlers)   │◄────────│   (Routes)   │
└──────┬───────┘         └──────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│         USE CASES                    │
│      (Business Logic)                │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│      REPOSITORIES                    │
│      (Interfaces)                    │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│         ADAPTERS                     │
│    (MySQL, WebSocket Hub)            │
└──────────────────────────────────────┘
```

## Capas del Proyecto

### 1️⃣ Domain Layer (Capa de Dominio)

**Responsabilidad:** Contiene la lógica de negocio pura, entidades y reglas del dominio.

**Ubicación:** `src/retas/domain/`

**Características:**
- ✅ Sin dependencias externas
- ✅ No conoce detalles de infraestructura
- ✅ Define interfaces (contratos)
- ✅ Entidades con métodos de negocio

**Archivos:**

#### Entities (Entidades)
```go
// entities/Reta.go
type Reta struct {
    ID                string
    ZonaID            string
    Titulo            string
    // ... campos de negocio
}

func NewReta(...) (*Reta, error) {
    // Lógica de creación y validación
}
```

#### Repositories (Interfaces)
```go
// repositories/reta_repository.go
type IRetaRepository interface {
    UnirseReta(retaID, usuarioID, nombreJugador string) (int, []entities.Jugador, error)
    CrearReta(reta *entities.Reta) (*entities.Reta, *entities.Jugador, error)
}
```

**Regla de Oro:** El dominio NO depende de nada más. Todo depende del dominio.

---

### 2️⃣ Application Layer (Capa de Aplicación)

**Responsabilidad:** Casos de uso, orquestación de la lógica de negocio.

**Ubicación:** `src/retas/application/`

**Características:**
- ✅ Depende solo de la capa de dominio
- ✅ Orquesta entidades y repositorios
- ✅ Implementa flujos de negocio
- ✅ No conoce detalles de BD o HTTP

**Archivos:**

```go
// UnirseReta_usecase.go
type UnirseRetaUseCase struct {
    retaRepo repositories.IRetaRepository  // Dependencia de INTERFAZ
}

func (uc *UnirseRetaUseCase) Execute(retaID, usuarioID, nombreJugador string) (int, []entities.Jugador, error) {
    // Orquestación: delega al repositorio
    return uc.retaRepo.UnirseReta(retaID, usuarioID, nombreJugador)
}
```

**Ventajas:**
- 🔄 Fácil de testear (mock del repositorio)
- 🔄 Lógica reutilizable
- 🔄 Independiente de infraestructura

---

### 3️⃣ Infrastructure Layer (Capa de Infraestructura)

**Responsabilidad:** Implementaciones concretas, adaptadores, I/O.

**Ubicación:** `src/retas/infraestructure/`

**Características:**
- ✅ Implementa interfaces del dominio
- ✅ Maneja detalles técnicos (BD, HTTP, WebSocket)
- ✅ Puede cambiar sin afectar el dominio

#### A. Adapters (Implementaciones)

```go
// adapters/MySQL.go
type MySQLRetaRepository struct {
    db *sql.DB
}

// Implementa IRetaRepository
func (repo *MySQLRetaRepository) UnirseReta(...) {
    // Código específico de MySQL
    tx, err := repo.db.Begin()
    // SELECT FOR UPDATE ...
    tx.Commit()
}
```

```go
// adapters/websocket_hub.go
type Hub struct {
    clients map[string]map[*Client]bool
    // ...
}

func (h *Hub) BroadcastToZone(zonaID string, message interface{}) {
    // Lógica específica de WebSocket
}
```

#### B. Controllers (Handlers)

```go
// controllers/WebSocket_controller.go
type WebSocketController struct {
    hub              *adapters.Hub
    unirseUseCase    *application.UnirseRetaUseCase  // Dependencia
    crearRetaUseCase *application.CrearRetaUseCase
}

func (wsc *WebSocketController) HandleWebSocket(c *gin.Context) {
    // 1. Recibe request HTTP/WebSocket
    // 2. Parsea y valida
    // 3. Llama al use case
    // 4. Devuelve respuesta
}
```

#### C. Routers (Rutas)

```go
// routers/retas_router.go
func RetasRouter(r *gin.Engine, wsController *controllers.WebSocketController) {
    retasGroup := r.Group("/ws")
    {
        retasGroup.GET("/retas", wsController.HandleWebSocket)
    }
}
```

#### D. Dependencies (Inyección de Dependencias)

```go
// dependencies_retas/dependencies.go
func InitRetas(r *gin.Engine) {
    // 1. Crear instancias concretas
    db, _ := core.NewMySQL()
    hub := adapters.NewHub()
    
    // 2. Crear repositorio (implementación)
    retaRepo := adapters.NewMySQLRetaRepository(db)
    
    // 3. Crear use cases (inyectar repositorio)
    unirseUseCase := application.NewUnirseRetaUseCase(retaRepo)
    crearRetaUseCase := application.NewCrearRetaUseCase(retaRepo)
    
    // 4. Crear controller (inyectar use cases)
    wsController := controllers.NewWebSocketController(hub, unirseUseCase, crearRetaUseCase)
    
    // 5. Registrar rutas
    routers.RetasRouter(r, wsController)
}
```

---

## 🔄 Flujo de Dependencias (Inversión de Control)

```
┌──────────────────────────────────────────────────┐
│               DOMAIN                             │
│  ┌────────────────────────────────────────┐     │
│  │  IRetaRepository (Interface)           │     │
│  │  - UnirseReta()                        │     │
│  │  - CrearReta()                         │     │
│  └────────────────┬───────────────────────┘     │
└───────────────────│──────────────────────────────┘
                    │
                    │ depende de (interface)
                    │
┌───────────────────▼──────────────────────────────┐
│              APPLICATION                         │
│  ┌────────────────────────────────────────┐     │
│  │  UnirseRetaUseCase                     │     │
│  │  {                                     │     │
│  │    retaRepo IRetaRepository  ◄────┐   │     │
│  │  }                                 │   │     │
│  └────────────────────────────────────┼───┘     │
└───────────────────────────────────────┼──────────┘
                                        │
                                        │ inyecta (implementación)
                                        │
┌───────────────────────────────────────▼──────────┐
│             INFRASTRUCTURE                       │
│  ┌────────────────────────────────────────┐     │
│  │  MySQLRetaRepository                   │     │
│  │  IMPLEMENTS IRetaRepository            │     │
│  │  {                                     │     │
│  │    db *sql.DB                          │     │
│  │  }                                     │     │
│  └────────────────────────────────────────┘     │
└──────────────────────────────────────────────────┘
```

**Inversión de Dependencias:**
- ❌ Use Case NO depende de MySQL directamente
- ✅ Use Case depende de la INTERFAZ IRetaRepository
- ✅ MySQL implementa la interfaz
- ✅ Se inyecta en tiempo de ejecución

---

## 🎯 Beneficios de Esta Arquitectura

### 1. Testabilidad
```go
// Mock del repositorio para testing
type MockRetaRepository struct {}

func (m *MockRetaRepository) UnirseReta(...) (int, []entities.Jugador, error) {
    return 5, []entities.Jugador{{ID: "1", Nombre: "Test"}}, nil
}

// Test del use case
func TestUnirseReta(t *testing.T) {
    mockRepo := &MockRetaRepository{}
    useCase := application.NewUnirseRetaUseCase(mockRepo)
    
    // Probar lógica sin tocar la BD real
    result, _, err := useCase.Execute("reta1", "user1", "Test")
    assert.NoError(t, err)
}
```

### 2. Cambio de Tecnología Sin Dolor

**Cambiar MySQL por PostgreSQL:**
```go
// Crear nuevo adapter
type PostgresRetaRepository struct {
    db *sql.DB
}

func (repo *PostgresRetaRepository) UnirseReta(...) {
    // Implementación para Postgres
}

// En dependencies.go, solo cambiar:
// retaRepo := adapters.NewMySQLRetaRepository(db)
retaRepo := adapters.NewPostgresRetaRepository(db)

// ✅ Use cases y controllers no cambian!
```

### 3. Separación de Responsabilidades

| Capa | Responsabilidad | Ejemplo |
|------|----------------|---------|
| Domain | Reglas de negocio | "Una reta no puede tener más de max_jugadores" |
| Application | Orquestación | "Primero verifica cupo, luego inserta" |
| Infrastructure | Detalles técnicos | "Usar transacción SQL con SELECT FOR UPDATE" |

### 4. Código Mantenible

```
Cambio necesario: Agregar validación de edad mínima

❌ Sin Clean Architecture:
- Modificar controller (mezclado con HTTP)
- Modificar repository (mezclado con SQL)
- Difícil de testear

✅ Con Clean Architecture:
- Agregar validación en entities/Reta.go
- Use case la llama automáticamente
- Fácil de testear (unit test)
```

---

## 📦 Comparación con MVC Tradicional

### MVC Tradicional:
```
Controller → Model → Database
     ↓
   View
```
- ❌ Controller conoce detalles de BD
- ❌ Model = tabla de BD
- ❌ Difícil de testear
- ❌ Lógica dispersa

### Clean Architecture:
```
Controller → UseCase → Repository (Interface)
                            ↓
                      Adapter (MySQL)
```
- ✅ Separación clara
- ✅ Testeable
- ✅ Flexible
- ✅ Lógica centralizada

---

## 🚀 Escalabilidad del Diseño

### Agregar Nueva Funcionalidad: "Cancelar Reta"

#### 1. Domain
```go
// repositories/reta_repository.go
type IRetaRepository interface {
    UnirseReta(...)
    CrearReta(...)
    CancelarReta(retaID string) error  // Nueva función
}
```

#### 2. Application
```go
// CancelarReta_usecase.go
type CancelarRetaUseCase struct {
    retaRepo repositories.IRetaRepository
}

func (uc *CancelarRetaUseCase) Execute(retaID string) error {
    return uc.retaRepo.CancelarReta(retaID)
}
```

#### 3. Infrastructure
```go
// adapters/MySQL.go
func (repo *MySQLRetaRepository) CancelarReta(retaID string) error {
    query := "DELETE FROM retas WHERE id = ?"
    _, err := repo.db.Exec(query, retaID)
    return err
}
```

#### 4. Controller
```go
// Agregar handler en WebSocket_controller.go
case "cancelar":
    wsc.handleCancelar(client, wsMsg)
```

**Total:** 4 archivos modificados/creados, sin tocar código existente ✅

---

## 📚 Referencias y Recursos

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)

---

**Implementado en:** Todo el proyecto `games_football_back/`
