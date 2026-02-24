# 🔐 Manejo de Concurrencia y Transacciones SQL

## Problema de Concurrencia

Cuando múltiples usuarios intentan unirse simultáneamente a una reta que tiene pocos cupos disponibles, pueden ocurrir **condiciones de carrera (race conditions)**:

### ❌ Sin Control de Concurrencia:
```
Tiempo | Usuario A                    | Usuario B
-------|------------------------------|---------------------------
T1     | SELECT jugadores_actuales    |
       | (resultado: 13)              |
T2     |                              | SELECT jugadores_actuales
       |                              | (resultado: 13)
T3     | UPDATE (13 + 1 = 14) ✓       |
T4     |                              | UPDATE (13 + 1 = 14) ✓
       | RESULTADO: 15 jugadores en una reta de max 14 ❌
```

## ✅ Solución: SELECT FOR UPDATE

### Implementación en MySQL.go

```go
// Iniciar transacción
tx, err := repo.db.Begin()

// SELECT FOR UPDATE bloquea la fila
query := "SELECT jugadores_actuales, max_jugadores FROM retas WHERE id = ? FOR UPDATE"
err = tx.QueryRow(query, retaID).Scan(&jugadoresActuales, &maxJugadores)

// Verificar cupo
if jugadoresActuales >= maxJugadores {
    tx.Rollback()  // ⚠️ Rollback: libera el bloqueo sin cambios
    return errors.New("reta llena")
}

// Incrementar contador
updateQuery := "UPDATE retas SET jugadores_actuales = jugadores_actuales + 1 WHERE id = ?"
tx.Exec(updateQuery, retaID)

// Insertar jugador
insertQuery := "INSERT INTO reta_jugadores (id, reta_id, usuario_id, nombre_jugador) VALUES (?, ?, ?, ?)"
tx.Exec(insertQuery, jugadorID, retaID, usuarioID, nombreJugador)

// Commit: libera el bloqueo y aplica los cambios
tx.Commit()
```

### ✅ Con Control de Concurrencia:
```
Tiempo | Usuario A                    | Usuario B
-------|------------------------------|---------------------------
T1     | BEGIN TRANSACTION            |
T2     | SELECT ... FOR UPDATE        |
       | (bloquea la fila)            |
T3     |                              | BEGIN TRANSACTION
T4     |                              | SELECT ... FOR UPDATE
       |                              | (ESPERA... bloqueado por A)
T5     | Verifica cupo (13 < 14) ✓    |
T6     | UPDATE (13 + 1 = 14)         |
T7     | INSERT jugador               |
T8     | COMMIT (libera bloqueo)      |
T9     |                              | (desbloquea y continúa)
T10    |                              | Verifica cupo (14 >= 14) ❌
T11    |                              | ROLLBACK (sin cambios)
       | RESULTADO: Exactamente 14 jugadores ✓
```

## 🔄 Flujo Completo de Transacción

### Caso 1: Unirse Exitoso

```
┌─────────────────────────────────────────────────────┐
│ 1. Cliente envía: {"accion": "unirse", ...}        │
└──────────────────┬──────────────────────────────────┘
                   ▼
┌─────────────────────────────────────────────────────┐
│ 2. WebSocket Controller recibe mensaje             │
└──────────────────┬──────────────────────────────────┘
                   ▼
┌─────────────────────────────────────────────────────┐
│ 3. UnirseRetaUseCase.Execute()                     │
└──────────────────┬──────────────────────────────────┘
                   ▼
┌─────────────────────────────────────────────────────┐
│ 4. MySQLRetaRepository.UnirseReta()                │
│    ┌──────────────────────────────────────────┐    │
│    │ BEGIN TRANSACTION                        │    │
│    │ SELECT ... FOR UPDATE (bloquea)          │    │
│    │ Verificar: jugadores < max ✓             │    │
│    │ UPDATE jugadores_actuales + 1            │    │
│    │ INSERT INTO reta_jugadores               │    │
│    │ COMMIT (libera bloqueo)                  │    │
│    └──────────────────────────────────────────┘    │
└──────────────────┬──────────────────────────────────┘
                   ▼
┌─────────────────────────────────────────────────────┐
│ 5. SELECT lista_jugadores (sin bloqueo)            │
└──────────────────┬──────────────────────────────────┘
                   ▼
┌─────────────────────────────────────────────────────┐
│ 6. Hub.BroadcastToZone() → Todos en zona_id        │
└─────────────────────────────────────────────────────┘
```

### Caso 2: Reta Llena (Rollback)

```
┌─────────────────────────────────────────────────────┐
│ 1. Cliente envía: {"accion": "unirse", ...}        │
└──────────────────┬──────────────────────────────────┘
                   ▼
┌─────────────────────────────────────────────────────┐
│ 2. MySQLRetaRepository.UnirseReta()                │
│    ┌──────────────────────────────────────────┐    │
│    │ BEGIN TRANSACTION                        │    │
│    │ SELECT ... FOR UPDATE (bloquea)          │    │
│    │ Verificar: jugadores >= max ❌           │    │
│    │ ROLLBACK (libera bloqueo, sin cambios)   │    │
│    └──────────────────────────────────────────┘    │
└──────────────────┬──────────────────────────────────┘
                   ▼
┌─────────────────────────────────────────────────────┐
│ 3. Error: "reta llena"                             │
└──────────────────┬──────────────────────────────────┘
                   ▼
┌─────────────────────────────────────────────────────┐
│ 4. WebSocket envía error SOLO al cliente           │
│    {"status": "error", "mensaje": "Reta llena"}    │
└─────────────────────────────────────────────────────┘
```

## 🛡️ Niveles de Protección

### 1. Bloqueo Pesimista (SELECT FOR UPDATE)
- ✅ Bloquea la fila durante la transacción
- ✅ Otros usuarios esperan hasta que se libere
- ✅ Previene actualizaciones conflictivas
- ⚠️ Puede causar esperas si las transacciones son largas

### 2. Unique Constraint en BD
```sql
UNIQUE KEY unique_usuario_reta (reta_id, usuario_id)
```
- ✅ Previene que un usuario se una dos veces
- ✅ Protección a nivel de base de datos
- ✅ Funciona incluso si el código tiene bugs

### 3. Validación en Código
```go
if jugadoresActuales >= maxJugadores {
    tx.Rollback()
    return errors.New("reta llena")
}
```
- ✅ Lógica de negocio explícita
- ✅ Mensajes de error claros
- ✅ Control antes de escribir en BD

## 📊 Niveles de Aislamiento

La transacción usa el nivel de aislamiento por defecto de MySQL:

```sql
-- En MySQL:
REPEATABLE READ (por defecto)

-- Características:
- Las lecturas dentro de una transacción son consistentes
- SELECT FOR UPDATE adquiere bloqueos de escritura
- Previene lecturas fantasma en rangos bloqueados
```

Para cambiar el nivel:
```sql
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
```

## 🎯 Best Practices Implementadas

### ✅ 1. Transacciones Cortas
```go
// ✓ Hacer BEGIN justo antes de necesitar bloqueo
tx, err := repo.db.Begin()

// ✓ COMMIT/ROLLBACK lo antes posible
defer func() {
    if err != nil {
        tx.Rollback()
    }
}()
```

### ✅ 2. Bloqueo Específico
```go
// ✓ Bloquear solo la fila necesaria
SELECT ... WHERE id = ? FOR UPDATE

// ✗ Evitar escaneos completos
SELECT ... FOR UPDATE  // sin WHERE
```

### ✅ 3. Orden Consistente de Bloqueos
```go
// ✓ Siempre bloquear tablas en el mismo orden:
// 1. retas (tabla padre)
// 2. reta_jugadores (tabla hija)

// Esto previene deadlocks
```

### ✅ 4. Manejo de Errores
```go
if err != nil {
    tx.Rollback()
    return fmt.Errorf("error descriptivo: %w", err)
}
```

## 🔬 Prueba de Concurrencia

Para probar que funciona correctamente:

```javascript
// Abrir la consola del navegador y ejecutar:
const ws = new WebSocket('ws://localhost:8080/ws/retas');

ws.onopen = () => {
  // Simular 20 usuarios uniéndose simultáneamente
  for (let i = 0; i < 20; i++) {
    setTimeout(() => {
      ws.send(JSON.stringify({
        accion: "unirse",
        zona_id: "z1",
        reta_id: "tu_reta_id",
        usuario_id: `user_${i}`,
        nombre: `Usuario ${i}`
      }));
    }, Math.random() * 100);
  }
};
```

**Resultado esperado:**
- Exactamente 14 usuarios registrados en la BD
- 6 usuarios reciben error "Reta llena"
- No hay duplicados
- No se excede el máximo

## 📚 Referencias

- [MySQL SELECT FOR UPDATE](https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html)
- [Go database/sql Transactions](https://go.dev/doc/database/execute-transactions)
- [Transaction Isolation Levels](https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html)

---

**Implementado en:** [MySQL.go](src/retas/infraestructure/adapters/MySQL.go)
