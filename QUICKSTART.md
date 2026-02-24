# 🚀 Quick Start - API Games Football

Guía rápida para poner en marcha la API en menos de 5 minutos.

## Paso 1: Configurar Base de Datos

```bash
# Crear la base de datos y tablas
mysql -u root -p < database_schema.sql
```

## Paso 2: Configurar Variables de Entorno

```bash
# Copiar archivo de ejemplo
copy .env.example .env

# Editar .env con tus credenciales (usando notepad o tu editor favorito)
notepad .env
```

Configuración mínima requerida:
```env
DB_USER=root
DB_PASSWORD=tu_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=games_football
```

## Paso 3: Instalar Dependencias

```bash
go mod download
```

## Paso 4: Ejecutar la API

```bash
go run main.go
```

Deberías ver:
```
Módulo de Retas inicializado correctamente
[GIN-debug] Listening and serving HTTP on :8080
```

## Paso 5: Probar con el Cliente HTML

Abre el archivo `cliente_websocket.html` en tu navegador:
```bash
start cliente_websocket.html
```

O accede a: `file:///ruta/completa/cliente_websocket.html`

## 🎮 Prueba Rápida

### 1. Crear una Reta
En el cliente HTML:
- Deja los valores por defecto
- Haz clic en **"Crear Reta"**
- Verás el mensaje de confirmación con el ID de la reta

### 2. Unirse a la Reta
En el cliente HTML (o abre otra pestaña):
- El ID de la reta se auto-completará
- Cambia el "Tu ID de Usuario" a `user_003`
- Cambia "Tu Nombre" a otro nombre (ej: "Pedro")
- Haz clic en **"Unirse a Reta"**
- Verás el broadcast con la lista actualizada de jugadores

## 📡 Endpoint WebSocket

```
ws://localhost:8080/ws/retas
```

## 📝 Mensajes de Ejemplo

### Crear Reta:
```json
{
  "accion": "crear",
  "zona_id": "zona_norte",
  "titulo": "Reta del Viernes",
  "fecha_hora": "2026-02-28 18:00:00",
  "max_jugadores": 14,
  "creador_id": "user_001",
  "creador_nombre": "Imanol"
}
```

### Unirse a Reta:
```json
{
  "accion": "unirse",
  "zona_id": "zona_norte",
  "reta_id": "uuid_de_la_reta",
  "usuario_id": "user_002",
  "nombre": "Carlos"
}
```

## ⚠️ Solución de Problemas

### Error: "Error loading .env file"
- Asegúrate de haber creado el archivo `.env` (sin extensión .txt)
- Verifica que esté en la raíz del proyecto (junto a main.go)

### Error: "Error al conectar a la base de datos"
- Verifica que MySQL esté corriendo
- Confirma las credenciales en el archivo `.env`
- Asegúrate de haber ejecutado el script `database_schema.sql`

### Error: WebSocket no conecta
- Verifica que la API esté corriendo en el puerto 8080
- Revisa la consola del navegador (F12) por errores
- Asegúrate de usar `ws://localhost:8080` (no `wss://`)

### Error: "Reta no encontrada"
- Verifica que el ID de la reta sea correcto
- Asegúrate de que la reta exista en la base de datos

## 🔍 Verificar que Todo Funciona

```sql
-- Ver las retas creadas
SELECT * FROM retas;

-- Ver los jugadores registrados
SELECT * FROM reta_jugadores;
```

## 🎯 Siguiente Paso

Lee el [README.md](README.md) completo para entender la arquitectura y todas las características.

---

**¿Necesitas ayuda?** Revisa los logs en la consola donde ejecutaste `go run main.go`
