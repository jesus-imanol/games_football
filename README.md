# API Games Football - WebSocket Real-Time

API de gestión de retas (partidos) de fútbol con WebSockets y base de datos MySQL/PostgreSQL, implementada con **Arquitectura Limpia** en Go.

## 🏗️ Arquitectura

El proyecto sigue los principios de Clean Architecture con separación clara de capas:

```
games_football_back/
├── main.go                          # Punto de entrada
├── go.mod                           # Dependencias
├── .env.example                     # Variables de entorno
├── database_schema.sql              # Schema de BD
└── src/
    ├── core/
    │   └── db_mysql.go             # Conexión a BD
    └── retas/
        ├── domain/                  # CAPA DE DOMINIO
        │   ├── entities/
        │   │   ├── Reta.go
        │   │   ├── Jugador.go
        │   │   └── WebSocketMessage.go
        │   └── repositories/
        │       └── reta_repository.go
        ├── application/              # CAPA DE APLICACIÓN
        │   ├── UnirseReta_usecase.go
        │   └── CrearReta_usecase.go
        └── infraestructure/          # CAPA DE INFRAESTRUCTURA
            ├── adapters/
            │   ├── MySQL.go          # Implementación del repo con transacciones
            │   └── websocket_hub.go  # Hub de WebSocket por zona
            ├── controllers/
            │   └── WebSocket_controller.go
            ├── routers/
            │   └── retas_router.go
            └── dependencies_retas/
                └── dependencies.go   # Inyección de dependencias
```

## 🚀 Instalación

### 1. Clonar y configurar

```bash
cd games_football_back
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
```

Edita el archivo `.env` con tus credenciales de base de datos:

```env
DB_USER=root
DB_PASSWORD=tu_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=games_football
```

### 3. Crear la base de datos

Ejecuta el script SQL:

```bash
mysql -u root -p < database_schema.sql
```

### 4. Instalar dependencias

```bash
go mod download
```

### 5. Ejecutar la API

```bash
go run main.go
```

La API estará disponible en: `http://localhost:8080`

## 📡 WebSocket Endpoint

**Endpoint:** `ws://localhost:8080/ws/retas`

### Conexión

Los clientes se conectan al WebSocket y envían mensajes JSON con dos posibles acciones: **"unirse"** o **"crear"**.

## 📝 Mensajes del Cliente

### 1. Acción: UNIRSE a una reta

Un jugador se une a una reta existente.

**Mensaje JSON:**
```json
{
  "accion": "unirse",
  "usuario_id": "123",
  "nombre": "Imanol",
  "reta_id": "r1",
  "zona_id": "z1"
}
```

**Lógica:**
1. Usa transacción SQL con `SELECT ... FOR UPDATE` para bloquear la fila
2. Verifica si hay cupo disponible (`jugadores_actuales < max_jugadores`)
3. Si está llena → rollback y error solo al cliente
4. Si hay cupo → incrementa contador e inserta en `reta_jugadores`
5. Hace commit y obtiene lista actualizada de jugadores
6. **Broadcast a todos** los clientes de esa `zona_id`

**Respuesta Broadcast (a todos en la zona):**
```json
{
  "status": "actualizacion",
  "reta_id": "r1",
  "jugadores_actuales": 11,
  "lista_jugadores": [
    {"id": "j1", "usuario_id": "123", "nombre": "Imanol"},
    {"id": "j2", "usuario_id": "456", "nombre": "Carlos"},
    ...
  ]
}
```

**Respuesta Error (solo al cliente):**
```json
{
  "status": "error",
  "mensaje": "Reta llena"
}
```

### 2. Acción: CREAR una nueva reta

Un usuario crea una nueva reta y automáticamente se convierte en el primer jugador.

**Mensaje JSON:**
```json
{
  "accion": "crear",
  "zona_id": "z1",
  "titulo": "Gran Reta Nocturna",
  "fecha_hora": "2026-02-24 20:00:00",
  "max_jugadores": 14,
  "creador_id": "123",
  "creador_nombre": "Imanol"
}
```

**Lógica:**
1. Inserta registro en tabla `retas` (genera UUID)
2. Inserta al creador en `reta_jugadores` como primer confirmado
3. **Broadcast a todos** los clientes de esa `zona_id`

**Respuesta Broadcast (a todos en la zona):**
```json
{
  "status": "nueva_reta",
  "reta": {
    "id": "uuid_generado",
    "titulo": "Gran Reta Nocturna",
    "fecha_hora": "2026-02-24 20:00:00",
    "max_jugadores": 14,
    "jugadores_actuales": 1,
    "lista_jugadores": [
      {"id": "j1", "usuario_id": "123", "nombre": "Imanol"}
    ]
  }
}
```

## 🔐 Manejo de Concurrencia

La API implementa **control de concurrencia optimista** usando:

### SELECT FOR UPDATE
```sql
SELECT jugadores_actuales, max_jugadores 
FROM retas 
WHERE id = ? 
FOR UPDATE
```

Esto **bloquea la fila** durante la transacción, garantizando que:
- ✅ Solo un usuario a la vez puede unirse
- ✅ No se exceda el máximo de jugadores
- ✅ Rollback automático si la reta está llena
- ✅ Sin condiciones de carrera (race conditions)

## 🗄️ Schema de Base de Datos

### Tabla: `retas`
```sql
CREATE TABLE retas (
    id VARCHAR(36) PRIMARY KEY,              -- UUID
    zona_id VARCHAR(50) NOT NULL,            -- Identificador de zona
    titulo VARCHAR(255) NOT NULL,
    fecha_hora DATETIME NOT NULL,
    max_jugadores INT DEFAULT 14,
    jugadores_actuales INT DEFAULT 0,
    creador_id VARCHAR(50) NOT NULL,
    creador_nombre VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Tabla: `reta_jugadores`
```sql
CREATE TABLE reta_jugadores (
    id VARCHAR(36) PRIMARY KEY,              -- UUID
    reta_id VARCHAR(36) NOT NULL,            -- FK a retas
    usuario_id VARCHAR(50) NOT NULL,
    nombre_jugador VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reta_id) REFERENCES retas(id) ON DELETE CASCADE,
    UNIQUE KEY unique_usuario_reta (reta_id, usuario_id)
);
```

## 🧪 Prueba con Cliente WebSocket

Puedes usar una extensión de VS Code como **WebSocket Client** o este código JavaScript:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/retas');

ws.onopen = () => {
  console.log('Conectado');
  
  // Crear una reta
  ws.send(JSON.stringify({
    accion: "crear",
    zona_id: "z1",
    titulo: "Reta del Viernes",
    fecha_hora: "2026-02-28 18:00:00",
    max_jugadores: 14,
    creador_id: "user1",
    creador_nombre: "Juan"
  }));
};

ws.onmessage = (event) => {
  console.log('Mensaje recibido:', JSON.parse(event.data));
};

// Unirse a una reta
setTimeout(() => {
  ws.send(JSON.stringify({
    accion: "unirse",
    usuario_id: "user2",
    nombre: "Pedro",
    reta_id: "id_de_la_reta_creada",
    zona_id: "z1"
  }));
}, 2000);
```

## 📦 Dependencias

- **gin-gonic/gin**: Framework web
- **gorilla/websocket**: WebSockets
- **go-sql-driver/mysql**: Driver MySQL
- **google/uuid**: Generación de UUIDs
- **joho/godotenv**: Variables de entorno

## ✅ Características Implementadas

- ✅ Arquitectura Limpia (Domain, Application, Infrastructure)
- ✅ Inyección de dependencias
- ✅ WebSocket Hub agrupado por `zona_id`
- ✅ Transacciones SQL con `SELECT FOR UPDATE`
- ✅ Manejo de concurrencia y rollbacks
- ✅ Broadcast por zona
- ✅ Mensajes de error individuales
- ✅ UUID para IDs únicos
- ✅ Separación clara de responsabilidades

## 🛠️ Comandos útiles

```bash
# Ejecutar
go run main.go

# Build
go build -o games_football_api

# Ejecutar binary
./games_football_api
```

---

**Desarrollado con 💚 siguiendo Clean Architecture**
