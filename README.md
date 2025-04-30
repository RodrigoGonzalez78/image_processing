### Documentacion de la api


### 1. Inicio de Sesión de Usuario

**Endpoint:** `/login`  
**Método:** `POST`  
**Descripción:** Autentica al usuario mediante su nombre de usuario y contraseña. Si las credenciales son válidas, genera y retorna un token JWT.

**Formato de solicitud:**
```json
{
  "userName": "juanperez",
  "password": "Pass1234"
}
```

**Requisitos:**
- El campo `userName` no debe estar vacío.
- El campo `password` no debe estar vacío.

**Respuestas:**

| Código | Descripción                                                                 |
|--------|-----------------------------------------------------------------------------|
| 201    | Autenticación exitosa. Se retorna un token JWT.                            |
| 400    | Nombre de usuario o contraseña inválidos, campos vacíos, o error de lógica.|
| 500    | Error interno al generar el token de autenticación.                        |

**Respuesta exitosa (`201 Created`):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR..."
}
```

### 2. Registro de Usuario

**Endpoint:** `/register`  
**Método:** `POST`  
**Descripción:** Registra un nuevo usuario validando el nombre de usuario, la longitud mínima de la contraseña y su unicidad. La contraseña es almacenada de forma segura usando hash.

**Formato de solicitud:**
```json
{
  "userName": "juanperez",
  "password": "Pass1234"
}
```

**Requisitos:**
- El campo `userName` no debe estar vacío.
- El campo `password` debe tener al menos 8 caracteres.
- El `userName` debe ser único en la base de datos.

**Respuestas:**

| Código | Descripción                                                                 |
|--------|-----------------------------------------------------------------------------|
| 201    | Usuario registrado exitosamente.                                            |
| 400    | Error de validación, datos incompletos o nombre de usuario ya existente.   |
| 500    | Error interno al registrar el usuario o al encriptar la contraseña.        |


### 3. Subida de Imagen

**Endpoint:** `/upload`  
**Método:** `POST`  
**Descripción:** Permite a un usuario autenticado subir una imagen. La imagen es validada por tipo, almacenada en el sistema de archivos y registrada en la base de datos con sus metadatos.

**Requiere Autenticación:** Sí (JWT en el contexto, validado por middleware)

**Formato de solicitud:**
- Tipo de contenido: `multipart/form-data`
- Campo requerido: `image` (archivo de imagen)

**Requisitos:**
- El campo `image` debe contener un archivo con extensión válida (`jpg`, `jpeg`, `png`, `gif`).
- El usuario debe estar autenticado.
- El archivo debe tener una extensión válida y ser una imagen reconocida.

**Respuestas:**

| Código | Descripción                                                                 |
|--------|-----------------------------------------------------------------------------|
| 201    | Imagen subida exitosamente. Retorna URL y metadatos de la imagen.          |
| 400    | Archivo no válido, faltante o con formato no permitido.                    |
| 401    | Usuario no autenticado o token inválido.                                   |
| 500    | Error interno al guardar el archivo o registrar en base de datos.          |

**Respuesta exitosa (`201 Created`):**
```json
{
  "message": "Imagen subida exitosamente",
  "image": {
    "url": "http://localhost:8080/images/juanperez/7a9dcd3e-c34e-4ab1-a1b2-0a14e2a0f527.png",
    "name": "7a9dcd3e-c34e-4ab1-a1b2-0a14e2a0f527.png",
    "size": 204567,
    "format": "png",
    "width": 800,
    "height": 600
  }
}
```

### 4. Obtener Imagen Subida

**Endpoint:** `/images/{usuario}/{nombre_archivo}`  
**Método:** `GET`  
**Descripción:** Retorna la imagen solicitada por el nombre de archivo y usuario, siempre que exista en el sistema de archivos.

**Parámetros de ruta:**
- `{usuario}`: Nombre de usuario asociado a la imagen.
- `{nombre_archivo}`: Nombre exacto del archivo con extensión (por ejemplo: `foto.png`).

**Requisitos:**
- La imagen debe existir en el directorio `uploads/{usuario}/`.
- No se permite acceder a directorios directamente (solo archivos específicos).

**Respuestas:**

| Código | Descripción                                       |
|--------|---------------------------------------------------|
| 200    | Imagen servida correctamente.                     |
| 403    | Se intentó acceder a un directorio.               |
| 404    | Archivo o directorio no encontrado.               |

**Ejemplo de solicitud:**
```http
GET /images/juanperez/7a9dcd3e-c34e-4ab1-a1b2-0a14e2a0f527.png
```

**Respuesta exitosa (`200 OK`):**
- La imagen se retorna directamente como contenido binario.
- El encabezado `Content-Type` será determinado automáticamente (por ejemplo: `image/png`).



### 5. Obtener Información de Imagen por ID

**Endpoint:** `/image/{id}`  
**Método:** `GET`  
**Descripción:** Retorna los metadatos e información de una imagen específica cargada por el usuario autenticado. Requiere token JWT.

**Encabezados requeridos:**
- `Authorization: Bearer <token>` (JWT generado en el login)

**Parámetros de ruta:**
- `{id}`: ID numérico de la imagen.

**Requisitos:**
- El ID debe ser un número entero válido.
- El usuario autenticado debe ser el propietario de la imagen.

**Respuestas:**

| Código | Descripción                                                              |
|--------|---------------------------------------------------------------------------|
| 200    | Imagen encontrada. Se retorna la información y la URL de acceso.         |
| 400    | El ID de la imagen es inválido.                                          |
| 401    | No autorizado. Token JWT inválido o imagen no pertenece al usuario.      |
| 404    | Imagen no encontrada.                                                    |

**Ejemplo de solicitud:**
```http
GET /image/123
Authorization: Bearer eyJhbGciOi...
```

**Respuesta exitosa (`200 OK`):**
```json
{
  "url": "http://localhost:8080/images/juanperez/7a9dcd3e-c34e-4ab1-a1b2-0a14e2a0f527.png",
  "name": "7a9dcd3e-c34e-4ab1-a1b2-0a14e2a0f527.png",
  "size": 204800,
  "format": "png",
  "width": 1024,
  "height": 768
}
```

### 6. Obtener Imágenes del Usuario (paginadas)

**Endpoint:** `/user-images`  
**Método:** `GET`  
**Descripción:** Retorna una lista paginada de todas las imágenes subidas por el usuario autenticado. Requiere autenticación mediante JWT.

**Encabezados requeridos:**
- `Authorization: Bearer <token>`

**Parámetros de consulta (query params):**
- `page` (opcional): Número de página (por defecto es `1`).
- `limit` (opcional): Cantidad de imágenes por página (por defecto es `10`).

**Respuestas:**

| Código | Descripción                                                             |
|--------|-------------------------------------------------------------------------|
| 200    | Lista de imágenes del usuario con paginación.                          |
| 401    | Usuario no autenticado o token inválido.                                |
| 500    | Error del servidor al obtener las imágenes.                             |

**Ejemplo de solicitud:**
```http
GET /user-images?page=2&limit=5
Authorization: Bearer eyJhbGciOi...
```

**Respuesta exitosa (`200 OK`):**
```json
{
  "page": 2,
  "limit": 5,
  "total": 23,
  "images": [
    {
      "id": 6,
      "name": "d234e123-cc8b-4a2e-a45a.png",
      "userName": "juanperez",
      "path": "uploads/juanperez/d234e123-cc8b-4a2e-a45a.png",
      "size": 183729,
      "format": "png",
      "width": 1024,
      "height": 768
    },
    ...
  ]
}
```

### 7. Transformar Imagen

**Endpoint:** `/images/{id}/transform`  
**Método:** `POST`  
**Descripción:** Aplica transformaciones a una imagen existente del usuario autenticado y retorna la imagen resultante.  
**Requiere autenticación mediante JWT.**

**Encabezados requeridos:**
- `Authorization: Bearer <token>`
- `Content-Type: application/json`

**Parámetros de ruta:**
- `id`: ID numérico de la imagen a transformar.

**Cuerpo del request (JSON):**
```json
{
  "transformations": {
    "resize": {
      "width": 800,
      "height": 600
    },
    "crop": {
      "x": 10,
      "y": 10,
      "width": 300,
      "height": 200
    },
    "rotate": 90,
    "filters": {
      "grayscale": true,
      "sepia": false
    },
    "format": "jpg"
  }
}
```

**Transformaciones soportadas:**
- **resize**: Cambia el tamaño de la imagen a las dimensiones especificadas.
- **crop**: Recorta la imagen desde una posición (`x`, `y`) con ancho y alto dados.
- **rotate**: Rota la imagen en grados.
- **filters**:
  - `grayscale`: Convierte la imagen a escala de grises.
  - `sepia`: Aplica un efecto sepia simulado.
- **format**: `"jpg"`, `"png"` o `"gif"` (por defecto `"png"`).

**Respuestas:**

| Código | Descripción                                                          |
|--------|----------------------------------------------------------------------|
| 200    | Imagen transformada retornada directamente en el cuerpo.            |
| 400    | Error en los datos enviados (ID inválido, JSON mal formado, etc.).  |
| 401    | Usuario no autenticado.                                              |
| 404    | Imagen no encontrada.                                                |
| 500    | Error interno al procesar la imagen.                                 |

**Ejemplo de solicitud:**
```http
POST /images/6/transform
Authorization: Bearer eyJhbGciOi...
Content-Type: application/json

{
  "transformations": {
    "resize": {
      "width": 600,
      "height": 400
    },
    "rotate": 180,
    "filters": {
      "grayscale": true
    },
    "format": "png"
  }
}
```

**Respuesta exitosa (`200 OK`):**
El cuerpo de la respuesta contiene directamente la imagen transformada, con el encabezado `Content-Type` apropiado (`image/png`, `image/jpeg`, etc.).