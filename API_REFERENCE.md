# API Reference — Games Football Backend

## Conexión

| Campo       | Valor                        |
|-------------|------------------------------|
| Protocolo   | WebSocket seguro (`wss://`)  |
| Host        | `apigamesfotball.chuy7x.space` |
| Endpoint WS Retas | `/ws/retas`                    |
| Endpoint WS Chat  | `/ws/retas/chat`               |
| Base REST   | `https://apigamesfotball.chuy7x.space` |
| URL WS Retas| `wss://apigamesfotball.chuy7x.space/ws/retas` |
| URL WS Chat | `wss://apigamesfotball.chuy7x.space/ws/retas/chat` |

> En producción usa el dominio `apigamesfotball.chuy7x.space`.

---

## Health Check (HTTP)

```
GET /health
```

**Respuesta:**
```json
{
  "status": "online",
  "message": "API Games Football está en línea ✓",
  "version": "1.0.0",
  "endpoints": {
    "websocket": "/ws/retas",
    "websocket_chat": "/ws/retas/chat"
  }
}
```

---

## Módulo de Usuarios (REST HTTP)

### 1. Registrar usuario

```
POST /api/usuarios/register
Content-Type: application/json
```

**Body:**
```json
{
  "username": "jesus-imanol",
  "password": "miPassword123",
  "nombre": "Jesús Imanol"
}
```

| Campo      | Tipo   | Obligatorio | Descripción                     |
|------------|--------|:-----------:|---------------------------------|
| `username` | string | ✅          | Nombre de usuario (único)       |
| `password` | string | ✅          | Contraseña (se hashea con bcrypt) |
| `nombre`   | string | ✅          | Nombre real del jugador         |

**Respuesta exitosa (201):**
```json
{
  "status": "success",
  "mensaje": "Usuario registrado exitosamente",
  "usuario": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "jesus-imanol",
    "nombre": "Jesús Imanol"
  }
}
```

**Errores posibles:**

| Código | `mensaje`                                | Causa                           |
|--------|------------------------------------------|---------------------------------|
| 400    | `"Campos requeridos: username, password, nombre"` | Faltan campos en el body |
| 409    | `"el username ya está registrado"`       | Username duplicado              |

---

### 2. Login

```
POST /api/usuarios/login
Content-Type: application/json
```

**Body:**
```json
{
  "username": "jesus-imanol",
  "password": "miPassword123"
}
```

| Campo      | Tipo   | Obligatorio | Descripción             |
|------------|--------|:-----------:|-------------------------|
| `username` | string | ✅          | Nombre de usuario       |
| `password` | string | ✅          | Contraseña en texto plano (se compara contra el hash bcrypt almacenado) |

**Respuesta exitosa (200):**
```json
{
  "status": "success",
  "mensaje": "Login exitoso",
  "usuario": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "jesus-imanol",
    "nombre": "Jesús Imanol"
  }
}
```

**Errores posibles:**

| Código | `mensaje`                                | Causa                                |
|--------|------------------------------------------|--------------------------------------|
| 400    | `"Campos requeridos: username, password"` | Faltan campos en el body            |
| 401    | `"credenciales inválidas"`               | Usuario no existe o password incorrecta |

> **Seguridad:** Las contraseñas se almacenan hasheadas con **bcrypt** (cost 10). El servidor nunca guarda ni retorna la contraseña en texto plano.

---

## Módulo de Retas (WebSocket)

### Flujo general de conexión

```
1. App registra usuario → POST /api/usuarios/register
2. App hace login       → POST /api/usuarios/login  (obtiene el id y nombre del usuario)
3. App abre conexión WS → wss://apigamesfotball.chuy7x.space/ws/retas
4. App envía JSON con { "zona_id": "...", "accion": "..." }
5. El servidor hace broadcast a todos los clientes de esa zona_id
6. App recibe JSON con el resultado en tiempo real
7. Para chat en vivo   → wss://apigamesfotball.chuy7x.space/ws/retas/chat (endpoint dedicado)
```

Todos los mensajes son **JSON** tanto de entrada como de salida.

---

### Mensajes que envía el cliente (Frontend → Servidor)

#### 1. Crear una Reta

```json
{
  "accion": "crear",
  "zona_id": "suchiapa_centro",
  "titulo": "Partido del domingo",
  "fecha_hora": "2026-03-01 10:00:00",
  "max_jugadores": 14,
  "creador_id": "550e8400-e29b-41d4-a716-446655440000",
  "creador_nombre": "Jesús Imanol"
}
```

| Campo           | Tipo   | Obligatorio | Descripción                                    |
|-----------------|--------|:-----------:|------------------------------------------------|
| `accion`        | string | ✅          | Siempre `"crear"`                              |
| `zona_id`       | string | ✅          | Identificador de la zona geográfica            |
| `titulo`        | string | ✅          | Nombre del partido                             |
| `fecha_hora`    | string | ✅          | Formato: `"YYYY-MM-DD HH:MM:SS"`               |
| `max_jugadores` | int    | ✅          | Número máximo de jugadores (ej: 14)            |
| `creador_id`    | string | ⬜          | ID del usuario creador (obtenido del login). Si no se envía, el servidor genera un UUID |
| `creador_nombre`| string | ✅          | Nombre del usuario que crea la reta            |

---

#### 2. Unirse a una Reta

```json
{
  "accion": "unirse",
  "zona_id": "suchiapa_centro",
  "reta_id": "uuid-de-la-reta",
  "usuario_id": "550e8400-e29b-41d4-a716-446655440000",
  "nombre": "Jesús Imanol"
}
```

| Campo       | Tipo   | Obligatorio | Descripción                          |
|-------------|--------|:-----------:|--------------------------------------|
| `accion`    | string | ✅          | Siempre `"unirse"`                   |
| `zona_id`   | string | ✅          | Identificador de la zona             |
| `reta_id`   | string | ✅          | UUID de la reta a la que se une      |
| `usuario_id`| string | ✅          | ID del usuario (obtenido del login). **Debe existir en la tabla `usuarios`** |
| `nombre`    | string | ✅          | Nombre del jugador (el servidor lo ignora y usa el nombre real de la BD) |

> **Importante:** El servidor valida que `usuario_id` exista en la tabla `usuarios`. Si no existe, retorna error `"el usuario no existe"`. El nombre mostrado en broadcasts siempre es el registrado en la base de datos, no el enviado por el cliente.

---

#### 3. Enviar mensaje de chat (vía `/ws/retas`)

También se puede enviar un mensaje de chat directamente desde la conexión principal de retas usando la acción `enviar_mensaje`.

```json
{
  "accion": "enviar_mensaje",
  "zona_id": "suchiapa_centro",
  "reta_id": "550e8400-e29b-41d4-a716-446655440000",
  "usuario_id": "u-001",
  "texto": "Llevo balón"
}
```

| Campo       | Tipo   | Obligatorio | Descripción                          |
|-------------|--------|:-----------:|--------------------------------------|
| `accion`    | string | ✅          | Siempre `"enviar_mensaje"`           |
| `zona_id`   | string | ✅          | Identificador de la zona             |
| `reta_id`   | string | ✅          | UUID de la reta                      |
| `usuario_id`| string | ✅          | ID del usuario (obtenido del login)  |
| `texto`     | string | ✅          | Contenido del mensaje (máx 500 chars)|

> **Nota:** Para una experiencia de chat dedicada, se recomienda usar el endpoint `/ws/retas/chat` documentado más abajo.

---

### Mensajes que recibe el cliente (Servidor → Frontend)

> Todos los clientes conectados a la misma `zona_id` reciben estos mensajes en tiempo real (broadcast).

#### Respuesta: retas_zona (al conectarse)

Se envía automáticamente al cliente en cuanto manda su primer mensaje con `zona_id`. Contiene todas las retas existentes de esa zona.

```json
{
  "status": "retas_zona",
  "retas": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "titulo": "Partido del domingo",
      "fecha_hora": "2026-03-01 10:00:00",
      "max_jugadores": 6,
      "jugadores_actuales": 1,
      "lista_jugadores": [
        {
          "id": "uuid-jugador",
          "nombre": "Jesús Imanol",
          "usuario_id": "u-001",
          "reta_id": "550e8400-e29b-41d4-a716-446655440000"
        }
      ],
      "historial_chat": [
        {
          "id": "msg-uuid",
          "reta_id": "550e8400-e29b-41d4-a716-446655440000",
          "usuario_id": "u-001",
          "nombre": "Jesús Imanol",
          "texto": "Llevo balón",
          "timestamp": "2026-02-27T04:35:00Z"
        }
      ]
    }
  ]
}
```

> Si no hay retas en esa zona, `retas` llega como array vacío `[]`. Si una reta no tiene mensajes, `historial_chat` llega como array vacío `[]`.

#### Respuesta: nueva_reta (al crear)

```json
{
  "status": "nueva_reta",
  "reta": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "titulo": "Partido del domingo",
    "fecha_hora": "2026-03-01 10:00:00",
    "max_jugadores": 14,
    "jugadores_actuales": 1,
    "lista_jugadores": [
      {
        "id": "uuid-jugador",
        "nombre": "Jesús Imanol",
        "usuario_id": "u-001",
        "reta_id": "550e8400-e29b-41d4-a716-446655440000"
      }
    ]
  }
}
```

#### Respuesta: actualizacion (al unirse)

```json
{
  "status": "actualizacion",
  "reta_id": "550e8400-e29b-41d4-a716-446655440000",
  "jugadores_actuales": 2,
  "lista_jugadores": [
    {
      "id": "uuid-jugador-1",
      "nombre": "Jesús Imanol",
      "usuario_id": "u-001",
      "reta_id": "550e8400-e29b-41d4-a716-446655440000"
    },
    {
      "id": "uuid-jugador-2",
      "nombre": "Carlos Dev",
      "usuario_id": "u-002",
      "reta_id": "550e8400-e29b-41d4-a716-446655440000"
    }
  ]
}
```

#### Respuesta: nuevo_mensaje (al enviar mensaje de chat)

Se envía a **todos** los clientes de la `zona_id` cuando alguien envía un mensaje de chat (ya sea vía `/ws/retas` con acción `enviar_mensaje` o vía `/ws/retas/chat`).

```json
{
  "status": "nuevo_mensaje",
  "reta_id": "550e8400-e29b-41d4-a716-446655440000",
  "mensaje_chat": {
    "id": "msg-uuid",
    "reta_id": "550e8400-e29b-41d4-a716-446655440000",
    "usuario_id": "u-001",
    "nombre": "Jesús Imanol",
    "texto": "Llevo balón",
    "timestamp": "2026-02-27T04:35:00Z"
  }
}
```

#### Respuesta: error

```json
{
  "status": "error",
  "mensaje": "Campos requeridos: reta_id, usuario_id, nombre"
}
```

**Posibles mensajes de error (WebSocket):**

| Mensaje                                                              | Causa                                    |
|----------------------------------------------------------------------|------------------------------------------|
| `"Formato de mensaje inválido"`                                      | JSON malformado                          |
| `"Acción no reconocida"`                                             | `accion` distinto de `crear` / `unirse` / `enviar_mensaje` |
| `"Campos requeridos: reta_id, usuario_id, nombre"`                   | Faltan campos en acción `unirse`         |
| `"Campos requeridos: titulo, fecha_hora, max_jugadores, creador_nombre"` | Faltan campos en acción `crear` |
| `"Campos requeridos: reta_id, usuario_id, texto"`                    | Faltan campos en acción `enviar_mensaje` |
| `"el usuario no existe"`                                             | `usuario_id` no encontrado en la tabla `usuarios` |
| `"el usuario ya está inscrito en esta reta"`                         | Intento de unirse dos veces              |
| `"reta llena"`                                                       | Se alcanzó `max_jugadores`               |
| `"reta no encontrada"`                                               | `reta_id` no existe                      |

---

## Módulo de Chat en Vivo (WebSocket dedicado)

### Endpoint

```
wss://apigamesfotball.chuy7x.space/ws/retas/chat
```

> Este es un WebSocket **independiente** del de retas. Está diseñado para la pantalla de chat dentro de una reta específica.

### Flujo de conexión del Chat

```
1. App hace login → POST /api/usuarios/login (obtiene id y nombre)
2. App abre WS   → wss://apigamesfotball.chuy7x.space/ws/retas/chat
3. App envía     → { "reta_id": "...", "zona_id": "..." }  (primer mensaje obligatorio)
4. Servidor responde → historial_chat con todos los mensajes previos
5. App envía     → { "usuario_id": "...", "texto": "..." }  (mensajes de chat)
6. Servidor hace broadcast → nuevo_mensaje a todos los clientes de esa zona
```

### Mensajes que envía el cliente (Frontend → Servidor)

#### 1. Unirse al chat (primer mensaje — obligatorio)

El primer mensaje al conectarse **debe** incluir `reta_id` y `zona_id` para registrarse en el chat de esa reta.

```json
{
  "reta_id": "550e8400-e29b-41d4-a716-446655440000",
  "zona_id": "suchiapa_centro"
}
```

| Campo     | Tipo   | Obligatorio | Descripción                          |
|-----------|--------|:-----------:|--------------------------------------|
| `reta_id` | string | ✅          | UUID de la reta                      |
| `zona_id` | string | ✅          | Identificador de la zona geográfica  |

> Al recibir este mensaje, el servidor registra al cliente en la zona y le envía el historial completo de mensajes de esa reta.

#### 2. Enviar mensaje de chat

Una vez registrado, los mensajes siguientes solo necesitan `usuario_id` y `texto`.

```json
{
  "usuario_id": "u-001",
  "texto": "Llevo balón"
}
```

| Campo       | Tipo   | Obligatorio | Descripción                           |
|-------------|--------|:-----------:|---------------------------------------|
| `usuario_id`| string | ✅          | ID del usuario (obtenido del login)   |
| `texto`     | string | ✅          | Contenido del mensaje (máx 500 chars) |

> **No es necesario** enviar `reta_id` ni `zona_id` en mensajes posteriores al primero. El servidor ya los tiene almacenados en la sesión.

### Mensajes que recibe el cliente (Servidor → Frontend)

#### Respuesta: historial_chat (al conectarse)

Se envía **solo al cliente** que acaba de unirse al chat. Contiene todos los mensajes previos de esa reta ordenados cronológicamente.

```json
{
  "status": "historial_chat",
  "reta_id": "550e8400-e29b-41d4-a716-446655440000",
  "mensajes": [
    {
      "id": "msg-uuid-1",
      "reta_id": "550e8400-e29b-41d4-a716-446655440000",
      "usuario_id": "u-001",
      "nombre": "Jesús Imanol",
      "texto": "¿A qué hora nos vemos?",
      "timestamp": "2026-02-27T04:30:00Z"
    },
    {
      "id": "msg-uuid-2",
      "reta_id": "550e8400-e29b-41d4-a716-446655440000",
      "usuario_id": "u-002",
      "nombre": "Carlos Dev",
      "texto": "A las 10, llego temprano",
      "timestamp": "2026-02-27T04:31:00Z"
    }
  ]
}
```

> Si no hay mensajes previos, `mensajes` llega como array vacío `[]`.

#### Respuesta: nuevo_mensaje (broadcast en tiempo real)

Se envía a **todos** los clientes conectados a la misma `zona_id` cuando alguien envía un mensaje.

```json
{
  "status": "nuevo_mensaje",
  "reta_id": "550e8400-e29b-41d4-a716-446655440000",
  "mensaje_chat": {
    "id": "msg-uuid-3",
    "reta_id": "550e8400-e29b-41d4-a716-446655440000",
    "usuario_id": "u-001",
    "nombre": "Jesús Imanol",
    "texto": "Llevo balón",
    "timestamp": "2026-02-27T04:35:00Z"
  }
}
```

#### Respuesta: error

```json
{
  "status": "error",
  "mensaje": "Campos requeridos: usuario_id, texto"
}
```

**Posibles mensajes de error (Chat WebSocket):**

| Mensaje                                                   | Causa                                         |
|-----------------------------------------------------------|-----------------------------------------------|
| `"Formato de mensaje inválido"`                           | JSON malformado                               |
| `"Primero envía reta_id y zona_id para unirte al chat"`  | Se intentó enviar mensaje sin el primer paso  |
| `"Campos requeridos: usuario_id, texto"`                  | Faltan campos en el mensaje de chat           |
| `"reta_id, usuario_id y texto son requeridos"`            | Campos vacíos                                 |

---

## Objetos de datos

### Usuario

| Campo       | Tipo   | Descripción                        |
|-------------|--------|------------------------------------|
| `id`        | string | UUID del usuario                   |
| `username`  | string | Nombre de usuario (único)          |
| `nombre`    | string | Nombre real del jugador            |

> La contraseña **nunca** se retorna en las respuestas.

### Jugador

| Campo       | Tipo   | Descripción                                    |
|-------------|--------|------------------------------------------------|
| `id`        | string | UUID del registro jugador                      |
| `nombre`    | string | Nombre real (obtenido de la tabla `usuarios`)  |
| `usuario_id`| string | ID del usuario                                 |
| `reta_id`   | string | UUID de la reta                                |

### Reta

| Campo                | Tipo   | Descripción                          |
|----------------------|--------|--------------------------------------|
| `id`                 | string | UUID de la reta                      |
| `titulo`             | string | Nombre del partido                   |
| `fecha_hora`         | string | Fecha y hora `"YYYY-MM-DD HH:MM:SS"` |
| `max_jugadores`      | int    | Cupo máximo de jugadores             |
| `jugadores_actuales` | int    | Cuántos jugadores hay actualmente    |
| `lista_jugadores`    | array  | Lista de objetos `Jugador`           |
| `historial_chat`     | array  | Lista de objetos `Mensaje` (historial del chat en vivo) |

### Mensaje

| Campo       | Tipo   | Descripción                                    |
|-------------|--------|------------------------------------------------|
| `id`        | string | UUID del mensaje                               |
| `reta_id`   | string | UUID de la reta a la que pertenece              |
| `usuario_id`| string | ID del usuario que envió el mensaje             |
| `nombre`    | string | Nombre real del usuario (obtenido con `JOIN`)   |
| `texto`     | string | Contenido del mensaje (máx 500 caracteres)      |
| `timestamp` | string | Fecha/hora de creación en formato ISO 8601      |

---

## Zonas disponibles (datos de prueba)

| `zona_id`          | Nombre            |
|--------------------|-------------------|
| `suchiapa_centro`  | Suchiapa Centro   |
| `suchiapa_norte`   | Suchiapa Norte    |
| `suchiapa_sur`     | Suchiapa Sur      |

## Usuarios de prueba

| `username`      | Password | Nombre        |
|-----------------|----------|---------------|
| `jesus-imanol`  | `123`    | Jesús Imanol  |
| `carlos-dev`    | `123`    | Carlos Dev    |

---

## Ejemplos de implementación

### JavaScript / React Native / Web

```js
const BASE_URL = 'https://apigamesfotball.chuy7x.space';

// ===== 1. Registrar usuario =====
const register = await fetch(`${BASE_URL}/api/usuarios/register`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'nuevo_jugador',
    password: 'miPassword123',
    nombre: 'Nuevo Jugador'
  })
});
const registerData = await register.json();
console.log('Registrado:', registerData.usuario);

// ===== 2. Login =====
const login = await fetch(`${BASE_URL}/api/usuarios/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'jesus-imanol',
    password: '123'
  })
});
const loginData = await login.json();
const usuario = loginData.usuario; // { id, username, nombre }

// ===== 3. Conectar al WebSocket =====
const socket = new WebSocket('wss://apigamesfotball.chuy7x.space/ws/retas');

socket.onopen = () => {
  console.log('Conectado');
};

socket.onmessage = (event) => {
  const data = JSON.parse(event.data);

  if (data.status === 'retas_zona') {
    console.log('Retas de la zona:', data.retas);
  }

  if (data.status === 'nueva_reta') {
    console.log('Nueva reta creada:', data.reta);
  }

  if (data.status === 'actualizacion') {
    console.log('Reta actualizada:', data.reta_id, data.jugadores_actuales);
  }

  if (data.status === 'error') {
    console.error('Error:', data.mensaje);
  }
};

// Crear una reta (usa el id del login)
socket.send(JSON.stringify({
  accion: 'crear',
  zona_id: 'suchiapa_centro',
  titulo: 'Partido del domingo',
  fecha_hora: '2026-03-01 10:00:00',
  max_jugadores: 14,
  creador_id: usuario.id,
  creador_nombre: usuario.nombre
}));

// Unirse a una reta
socket.send(JSON.stringify({
  accion: 'unirse',
  zona_id: 'suchiapa_centro',
  reta_id: '550e8400-e29b-41d4-a716-446655440000',
  usuario_id: usuario.id,
  nombre: usuario.nombre
}));

// Enviar mensaje de chat (vía /ws/retas)
socket.send(JSON.stringify({
  accion: 'enviar_mensaje',
  zona_id: 'suchiapa_centro',
  reta_id: '550e8400-e29b-41d4-a716-446655440000',
  usuario_id: usuario.id,
  texto: 'Llevo balón'
}));
```

### JavaScript — Chat dedicado (`/ws/retas/chat`)

```js
const retaId = '550e8400-e29b-41d4-a716-446655440000';

// ===== Conectar al WebSocket de Chat =====
const chatSocket = new WebSocket('wss://apigamesfotball.chuy7x.space/ws/retas/chat');

chatSocket.onopen = () => {
  console.log('Chat conectado');

  // Primer mensaje: unirse al chat de una reta (enviar reta_id + zona_id)
  chatSocket.send(JSON.stringify({
    reta_id: retaId,
    zona_id: 'suchiapa_centro'
  }));
};

chatSocket.onmessage = (event) => {
  const data = JSON.parse(event.data);

  if (data.status === 'historial_chat') {
    // Historial completo al conectarse
    console.log('Historial:', data.mensajes);
  }

  if (data.status === 'nuevo_mensaje') {
    // Mensaje nuevo en tiempo real
    console.log(`${data.mensaje_chat.nombre}: ${data.mensaje_chat.texto}`);
  }

  if (data.status === 'error') {
    console.error('Error:', data.mensaje);
  }
};

// Enviar un mensaje de chat
chatSocket.send(JSON.stringify({
  usuario_id: usuario.id,
  texto: 'Llevo balón'
}));
```

### Flutter / Dart

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

const baseUrl = 'https://apigamesfotball.chuy7x.space';

// ===== 1. Registrar usuario =====
final registerRes = await http.post(
  Uri.parse('$baseUrl/api/usuarios/register'),
  headers: {'Content-Type': 'application/json'},
  body: jsonEncode({
    'username': 'nuevo_jugador',
    'password': 'miPassword123',
    'nombre': 'Nuevo Jugador',
  }),
);
print('Register: ${registerRes.body}');

// ===== 2. Login =====
final loginRes = await http.post(
  Uri.parse('$baseUrl/api/usuarios/login'),
  headers: {'Content-Type': 'application/json'},
  body: jsonEncode({
    'username': 'jesus-imanol',
    'password': '123',
  }),
);
final loginData = jsonDecode(loginRes.body);
final usuario = loginData['usuario']; // { id, username, nombre }

// ===== 3. WebSocket =====
final channel = WebSocketChannel.connect(
  Uri.parse('wss://apigamesfotball.chuy7x.space/ws/retas'),
);

channel.stream.listen((message) {
  final data = jsonDecode(message);

  if (data['status'] == 'retas_zona') {
    print('Retas: ${data['retas']}');
  }
  if (data['status'] == 'nueva_reta') {
    print('Nueva reta: ${data['reta']['titulo']}');
  }
  if (data['status'] == 'actualizacion') {
    print('Jugadores: ${data['jugadores_actuales']}');
  }
  if (data['status'] == 'error') {
    print('Error: ${data['mensaje']}');
  }
});

// Crear una reta
channel.sink.add(jsonEncode({
  'accion': 'crear',
  'zona_id': 'suchiapa_centro',
  'titulo': 'Partido del domingo',
  'fecha_hora': '2026-03-01 10:00:00',
  'max_jugadores': 14,
  'creador_id': usuario['id'],
  'creador_nombre': usuario['nombre'],
}));

// Unirse a una reta
channel.sink.add(jsonEncode({
  'accion': 'unirse',
  'zona_id': 'suchiapa_centro',
  'reta_id': '550e8400-e29b-41d4-a716-446655440000',
  'usuario_id': usuario['id'],
  'nombre': usuario['nombre'],
}));
```

### Flutter / Dart — Chat dedicado (`/ws/retas/chat`)

```dart
final retaId = '550e8400-e29b-41d4-a716-446655440000';

// ===== Conectar al WebSocket de Chat =====
final chatChannel = WebSocketChannel.connect(
  Uri.parse('wss://apigamesfotball.chuy7x.space/ws/retas/chat'),
);

chatChannel.stream.listen((message) {
  final data = jsonDecode(message);

  if (data['status'] == 'historial_chat') {
    print('Historial: ${data['mensajes']}');
  }
  if (data['status'] == 'nuevo_mensaje') {
    final msg = data['mensaje_chat'];
    print('${msg['nombre']}: ${msg['texto']}');
  }
  if (data['status'] == 'error') {
    print('Error: ${data['mensaje']}');
  }
});

// Primer mensaje: unirse al chat de una reta
chatChannel.sink.add(jsonEncode({
  'reta_id': retaId,
  'zona_id': 'suchiapa_centro',
}));

// Enviar mensaje de chat
chatChannel.sink.add(jsonEncode({
  'usuario_id': usuario['id'],
  'texto': 'Llevo balón',
}));
```

### Swift (iOS)

```swift
import Foundation

let baseURL = "https://apigamesfotball.chuy7x.space"

// ===== 1. Login =====
var loginReq = URLRequest(url: URL(string: "\(baseURL)/api/usuarios/login")!)
loginReq.httpMethod = "POST"
loginReq.setValue("application/json", forHTTPHeaderField: "Content-Type")
loginReq.httpBody = try! JSONSerialization.data(withJSONObject: [
    "username": "jesus-imanol",
    "password": "123"
])

URLSession.shared.dataTask(with: loginReq) { data, _, _ in
    guard let data = data,
          let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
          let usuario = json["usuario"] as? [String: Any],
          let userId = usuario["id"] as? String,
          let nombre = usuario["nombre"] as? String else { return }

    // ===== 2. Conectar WebSocket =====
    let wsURL = URL(string: "wss://apigamesfotball.chuy7x.space/ws/retas")!
    let task = URLSession.shared.webSocketTask(with: wsURL)
    task.resume()

    func listen() {
        task.receive { result in
            switch result {
            case .success(let message):
                if case .string(let text) = message,
                   let data = text.data(using: .utf8),
                   let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any] {
                    let status = json["status"] as? String
                    print("Status: \(status ?? "")")
                }
                listen()
            case .failure(let error):
                print("Error: \(error)")
            }
        }
    }
    listen()

    // Crear reta
    let msg: [String: Any] = [
        "accion": "crear",
        "zona_id": "suchiapa_centro",
        "titulo": "Partido del domingo",
        "fecha_hora": "2026-03-01 10:00:00",
        "max_jugadores": 14,
        "creador_id": userId,
        "creador_nombre": nombre
    ]
    let jsonData = try! JSONSerialization.data(withJSONObject: msg)
    task.send(.string(String(data: jsonData, encoding: .utf8)!)) { _ in }
}.resume()
```

### Swift (iOS) — Chat dedicado (`/ws/retas/chat`)

```swift
let retaId = "550e8400-e29b-41d4-a716-446655440000"

// ===== Conectar al WebSocket de Chat =====
let chatURL = URL(string: "wss://apigamesfotball.chuy7x.space/ws/retas/chat")!
let chatTask = URLSession.shared.webSocketTask(with: chatURL)
chatTask.resume()

func listenChat() {
    chatTask.receive { result in
        switch result {
        case .success(let message):
            if case .string(let text) = message,
               let data = text.data(using: .utf8),
               let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any] {
                let status = json["status"] as? String
                if status == "historial_chat" {
                    print("Historial: \(json["mensajes"] ?? [])")
                } else if status == "nuevo_mensaje" {
                    if let msg = json["mensaje_chat"] as? [String: Any] {
                        print("\(msg["nombre"] ?? ""): \(msg["texto"] ?? "")")
                    }
                }
            }
            listenChat()
        case .failure(let error):
            print("Error chat: \(error)")
        }
    }
}
listenChat()

// Primer mensaje: unirse al chat
let joinMsg: [String: Any] = [
    "reta_id": retaId,
    "zona_id": "suchiapa_centro"
]
let joinData = try! JSONSerialization.data(withJSONObject: joinMsg)
chatTask.send(.string(String(data: joinData, encoding: .utf8)!)) { _ in }

// Enviar mensaje de chat
let chatMsg: [String: Any] = [
    "usuario_id": userId,
    "texto": "Llevo balón"
]
let chatData = try! JSONSerialization.data(withJSONObject: chatMsg)
chatTask.send(.string(String(data: chatData, encoding: .utf8)!)) { _ in }
```

---

## Notas importantes

- **Flujo obligatorio:** Primero registrar/login para obtener el `id` del usuario, luego usarlo en los mensajes WebSocket.
- **Contraseñas:** Se almacenan hasheadas con **bcrypt** (cost 10). Nunca se retornan en las respuestas.
- **Validación de usuario:** Al unirse a una reta, el servidor verifica que `usuario_id` exista en la tabla `usuarios`. Si no existe, retorna error.
- **Nombre real:** En los broadcasts (lista de jugadores y mensajes de chat), el nombre se obtiene de la tabla `usuarios` con un `JOIN`, no del campo enviado por el cliente.
- **`zona_id`** es obligatorio en **todos** los mensajes WebSocket. Es el canal del broadcast — solo los clientes de la misma zona reciben las actualizaciones.
- Un usuario no puede unirse dos veces a la misma reta (restricción `UNIQUE` en base de datos).
- El creador de una reta queda automáticamente inscrito como primer jugador.
- `fecha_hora` debe tener exactamente el formato `"YYYY-MM-DD HH:MM:SS"`.
- **Chat en vivo:** Se puede usar desde `/ws/retas` (acción `enviar_mensaje`) o desde el endpoint dedicado `/ws/retas/chat`.
- **Endpoint `/ws/retas/chat`:** Es una conexión WebSocket independiente diseñada para la pantalla de chat. Requiere como primer mensaje `reta_id` + `zona_id`, y los mensajes posteriores solo necesitan `usuario_id` + `texto`.
