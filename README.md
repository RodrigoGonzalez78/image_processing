# 🖼️ API REST para Procesamiento de Imágenes

Esta API permite aplicar transformaciones a imágenes previamente subidas por el usuario, tales como redimensionamiento, recorte, rotación y filtros (escala de grises, sepia). Las imágenes se almacenan en **MinIO** y se procesan en tiempo real utilizando **Go**.

![Diagrama de arquitectura](/assets/Diagrama.png)

---

## 🚀 Ejecución del Proyecto

Para iniciar la aplicación localmente utilizando Docker:

```bash
sudo docker compose up --build
```

Esto levantará tanto el backend como MinIO en una red compartida.

---

## 📚 Documentación de la API

La documentación Swagger está disponible en:

```
http://localhost:8080/swagger/index.html
```

Desde allí podés explorar y probar los endpoints disponibles, incluyendo autenticación, carga y transformación de imágenes.

---

## 🔐 Autenticación

La API utiliza autenticación JWT. Para realizar peticiones protegidas, agregá el token en el encabezado:

```
Authorization: Bearer YOUR_TOKEN_HERE
```

> ⚠️ Reemplazá `YOUR_TOKEN_HERE` con tu token JWT válido y `{id}` con el ID de la imagen que querés transformar.

---

## 📦 Ejemplos de Transformaciones con `curl`

### 🧪 Ejemplo 1: Redimensionar a 800x600 y devolver en JPG

```bash
curl -X POST http://localhost:8080/images/123/transform \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -H "Accept: image/jpeg" \
  -d '{
    "transformations": {
      "resize": { "width": 800, "height": 600 },
      "crop": { "width": 0, "height": 0, "x": 0, "y": 0 },
      "rotate": 0,
      "format": "jpg",
      "filters": { "grayscale": false, "sepia": false }
    }
  }' --output resized.jpg
```

---

### 🧪 Ejemplo 2: Recortar 400x300 desde posición (100, 50)

```bash
curl -X POST http://localhost:8080/images/123/transform \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -H "Accept: image/png" \
  -d '{
    "transformations": {
      "resize": { "width": 0, "height": 0 },
      "crop": { "width": 400, "height": 300, "x": 100, "y": 50 },
      "rotate": 0,
      "format": "png",
      "filters": { "grayscale": false, "sepia": false }
    }
  }' --output cropped.png
```

---

### 🧪 Ejemplo 3: Rotar 90 grados y aplicar escala de grises

```bash
curl -X POST http://localhost:8080/images/123/transform \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -H "Accept: image/jpeg" \
  -d '{
    "transformations": {
      "resize": { "width": 0, "height": 0 },
      "crop": { "width": 0, "height": 0, "x": 0, "y": 0 },
      "rotate": 90,
      "format": "jpeg",
      "filters": { "grayscale": true, "sepia": false }
    }
  }' --output rotated_grayscale.jpg
```

---

### 🧪 Ejemplo 4: Redimensionar a 320x240 y aplicar filtro sepia en GIF

```bash
curl -X POST http://localhost:8080/images/123/transform \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -H "Accept: image/gif" \
  -d '{
    "transformations": {
      "resize": { "width": 320, "height": 240 },
      "crop": { "width": 0, "height": 0, "x": 0, "y": 0 },
      "rotate": 0,
      "format": "gif",
      "filters": { "grayscale": false, "sepia": true }
    }
  }' --output sepia.gif
```

---

## 🧩 Funcionalidades Clave

* 📤 Subida y almacenamiento de imágenes en MinIO
* 🔄 Transformación en tiempo real usando `imaging`
* 🔐 Autenticación JWT
* 🧪 Swagger UI para testing de endpoints
* 🐳 Contenerización con Docker

