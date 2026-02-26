# API Reference — Games Football Backend

## Conexión

| Campo       | Valor                        |
|-------------|------------------------------|
| Protocolo   | WebSocket (`ws://`)          |
| Host        | `apigamesfotball.chuy7x.space` |
| Endpoint    | `/ws/retas`                    |
| URL completa| `ws://apigamesfotball.chuy7x.space/ws/retas` |

> En producción usa el dominio `apigamesfotball.chuy7x.space`.

---

## Health Check (HTTP)

```
GET http://apigamesfotball.chuy7x.space/
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

## Flujo general de conexión

```
1. App abre conexión WebSocket → ws://apigamesfotball.chuy7x.space/ws/retas
2. App envía JSON con { "zona_id": "...", "accion": "..." }
3. El servidor hace broadcast a todos los clientes de esa zona_id
4. App recibe JSON con el resultado en tiempo real
```

Todos los mensajes son **JSON** tanto de entrada como de salida.

---

## Mensajes que envía el cliente (Frontend → Servidor)

### 1. Crear una Reta

```json
{
  "accion": "crear",
  "zona_id": "zona_norte",
  "titulo": "Partido del domingo",
  "fecha_hora": "2026-03-01 10:00:00",
  "max_jugadores": 14,
  "creador_nombre": "Carlos"
}
```

| Campo           | Tipo   | Obligatorio | Descripción                                    |
|-----------------|--------|:-----------:|------------------------------------------------|
| `accion`        | string | ✅          | Siempre `"crear"`                              |
| `zona_id`       | string | ✅          | Identificador de la zona geográfica            |
| `titulo`        | string | ✅          | Nombre del partido                             |
| `fecha_hora`    | string | ✅          | Formato: `"YYYY-MM-DD HH:MM:SS"`               |
| `max_jugadores` | int    | ✅          | Número máximo de jugadores (ej: 14)            |
| `creador_id`    | string | ⬜          | ID del creador. Si no se envía, el servidor genera un UUID automáticamente |
| `creador_nombre`| string | ✅          | Nombre del usuario que crea la reta            |

---

### 2. Unirse a una Reta

```json
{
  "accion": "unirse",
  "zona_id": "zona_norte",
  "reta_id": "uuid-de-la-reta",
  "usuario_id": "user_456",
  "nombre": "Pedro"
}
```

| Campo       | Tipo   | Obligatorio | Descripción                          |
|-------------|--------|:-----------:|--------------------------------------|
| `accion`    | string | ✅          | Siempre `"unirse"`                   |
| `zona_id`   | string | ✅          | Identificador de la zona             |
| `reta_id`   | string | ✅          | UUID de la reta a la que se une      |
| `usuario_id`| string | ✅          | ID del usuario que se une            |
| `nombre`    | string | ✅          | Nombre del jugador                   |

---

## Mensajes que recibe el cliente (Servidor → Frontend)

> Todos los clientes conectados a la misma `zona_id` reciben estos mensajes en tiempo real (broadcast).

### Respuesta: nueva_reta (al crear)

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
        "nombre": "Carlos",
        "usuario_id": "user_123",
        "reta_id": "550e8400-e29b-41d4-a716-446655440000"
      }
    ]
  }
}
```

### Respuesta: actualizacion (al unirse)

```json
{
  "status": "actualizacion",
  "reta_id": "550e8400-e29b-41d4-a716-446655440000",
  "jugadores_actuales": 2,
  "lista_jugadores": [
    {
      "id": "uuid-jugador-1",
      "nombre": "Carlos",
      "usuario_id": "user_123",
      "reta_id": "550e8400-e29b-41d4-a716-446655440000"
    },
    {
      "id": "uuid-jugador-2",
      "nombre": "Pedro",
      "usuario_id": "user_456",
      "reta_id": "550e8400-e29b-41d4-a716-446655440000"
    }
  ]
}
```

### Respuesta: error

```json
{
  "status": "error",
  "mensaje": "Campos requeridos: reta_id, usuario_id, nombre"
}
```

**Posibles mensajes de error:**

| Mensaje                                                              | Causa                                    |
|----------------------------------------------------------------------|------------------------------------------|
| `"Formato de mensaje inválido"`                                      | JSON malformado                          |
| `"Acción no reconocida"`                                             | `accion` distinto de `crear` / `unirse`  |
| `"Campos requeridos: reta_id, usuario_id, nombre"`                   | Faltan campos en acción `unirse`         |
| `"Campos requeridos: titulo, fecha_hora, max_jugadores, creador_nombre"` | Faltan campos en acción `crear` |

---

## Objetos de datos

### Jugador

| Campo       | Tipo   | Descripción               |
|-------------|--------|---------------------------|
| `id`        | string | UUID del registro jugador |
| `nombre`    | string | Nombre del jugador        |
| `usuario_id`| string | ID del usuario            |
| `reta_id`   | string | UUID de la reta           |

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

## Ejemplos de implementación

### JavaScript / React Native / Web

```js
const socket = new WebSocket('ws://apigamesfotball.chuy7x.space/ws/retas');

socket.onopen = () => {
  console.log('Conectado');
};

socket.onmessage = (event) => {
  const data = JSON.parse(event.data);

  if (data.status === 'nueva_reta') {
    // Agregar la nueva reta a la lista
    console.log('Nueva reta creada:', data.reta);
  }

  if (data.status === 'actualizacion') {
    // Actualizar jugadores de la reta
    console.log('Reta actualizada:', data.reta_id, data.jugadores_actuales);
  }

  if (data.status === 'error') {
    console.error('Error:', data.mensaje);
  }
};

// Crear una reta (creador_id es opcional, el servidor lo genera si no se envía)
socket.send(JSON.stringify({
  accion: 'crear',
  zona_id: 'zona_norte',
  titulo: 'Partido del domingo',
  fecha_hora: '2026-03-01 10:00:00',
  max_jugadores: 14,
  creador_nombre: 'Carlos'
}));

// Unirse a una reta
socket.send(JSON.stringify({
  accion: 'unirse',
  zona_id: 'zona_norte',
  reta_id: '550e8400-e29b-41d4-a716-446655440000',
  usuario_id: 'user_456',
  nombre: 'Pedro'
}));
```

### Flutter / Dart

```dart
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';

final channel = WebSocketChannel.connect(
  Uri.parse('ws://apigamesfotball.chuy7x.space/ws/retas'),
);

// Escuchar mensajes
channel.stream.listen((message) {
  final data = jsonDecode(message);

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
  'zona_id': 'zona_norte',
  'titulo': 'Partido del domingo',
  'fecha_hora': '2026-03-01 10:00:00',
  'max_jugadores': 14,
  'creador_nombre': 'Carlos',
}));

// Unirse a una reta
channel.sink.add(jsonEncode({
  'accion': 'unirse',
  'zona_id': 'zona_norte',
  'reta_id': '550e8400-e29b-41d4-a716-446655440000',
  'usuario_id': 'user_456',
  'nombre': 'Pedro',
}));
```

### Swift (iOS)

```swift
import Foundation

let url = URL(string: "ws://apigamesfotball.chuy7x.space/ws/retas")!
let task = URLSession.shared.webSocketTask(with: url)
task.resume()

// Escuchar mensajes
func listen() {
    task.receive { result in
        switch result {
        case .success(let message):
            if case .string(let text) = message,
               let data = text.data(using: .utf8),
               let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any] {
                let status = json["status"] as? String
                print("Status recibido: \(status ?? "")")
            }
            listen() // seguir escuchando
        case .failure(let error):
            print("Error: \(error)")
        }
    }
}
listen()

// Crear una reta
let msg: [String: Any] = [
    "accion": "crear",
    "zona_id": "zona_norte",
    "titulo": "Partido del domingo",
    "fecha_hora": "2026-03-01 10:00:00",
    "max_jugadores": 14,
    "creador_nombre": "Carlos"
]
let jsonData = try! JSONSerialization.data(withJSONObject: msg)
let jsonString = String(data: jsonData, encoding: .utf8)!
task.send(.string(jsonString)) { _ in }
```

---

## Notas importantes

- **`zona_id`** es obligatorio en **todos** los mensajes. Es el canal por el que se hace el broadcast — solo los clientes de la misma zona reciben las actualizaciones.
- Un usuario no puede unirse dos veces a la misma reta (restricción única en base de datos).
- El creador de una reta queda automáticamente inscrito como primer jugador.
- `fecha_hora` debe tener exactamente el formato `"YYYY-MM-DD HH:MM:SS"`.
- Los UUIDs (`id`, `reta_id`, `usuario_id`) son generados por el servidor automáticamente para `id`, pero `usuario_id` y `creador_id` los provee el cliente (puede ser el ID de tu sistema de autenticación).
