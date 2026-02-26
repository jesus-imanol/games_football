# API Reference — Games Football Backend

## Conexión

| Campo       | Valor                        |
|-------------|------------------------------|
| Protocolo   | WebSocket seguro (`wss://`)  |
| Host        | `apigamesfotball.chuy7x.space` |
| Endpoint WS | `/ws/retas`                    |
| Base REST   | `https://apigamesfotball.chuy7x.space` |
| URL WS completa| `wss://apigamesfotball.chuy7x.space/ws/retas` |

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
    "websocket": "/ws/retas"
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
      ]
    }
  ]
}
```

> Si no hay retas en esa zona, `retas` llega como array vacío `[]`.

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
| `"Acción no reconocida"`                                             | `accion` distinto de `crear` / `unirse`  |
| `"Campos requeridos: reta_id, usuario_id, nombre"`                   | Faltan campos en acción `unirse`         |
| `"Campos requeridos: titulo, fecha_hora, max_jugadores, creador_nombre"` | Faltan campos en acción `crear` |
| `"el usuario no existe"`                                             | `usuario_id` no encontrado en la tabla `usuarios` |
| `"el usuario ya está inscrito en esta reta"`                         | Intento de unirse dos veces              |
| `"reta llena"`                                                       | Se alcanzó `max_jugadores`               |
| `"reta no encontrada"`                                               | `reta_id` no existe                      |

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

---

## Notas importantes

- **Flujo obligatorio:** Primero registrar/login para obtener el `id` del usuario, luego usarlo en los mensajes WebSocket.
- **Contraseñas:** Se almacenan hasheadas con **bcrypt** (cost 10). Nunca se retornan en las respuestas.
- **Validación de usuario:** Al unirse a una reta, el servidor verifica que `usuario_id` exista en la tabla `usuarios`. Si no existe, retorna error.
- **Nombre real:** En los broadcasts (lista de jugadores), el nombre se obtiene de la tabla `usuarios` con un `JOIN`, no del campo enviado por el cliente.
- **`zona_id`** es obligatorio en **todos** los mensajes WebSocket. Es el canal del broadcast — solo los clientes de la misma zona reciben las actualizaciones.
- Un usuario no puede unirse dos veces a la misma reta (restricción `UNIQUE` en base de datos).
- El creador de una reta queda automáticamente inscrito como primer jugador.
- `fecha_hora` debe tener exactamente el formato `"YYYY-MM-DD HH:MM:SS"`.
