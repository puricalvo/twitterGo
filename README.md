
---
```go
# ⚙️ TwitterGo - API (Backend)

Este es el motor de **TwitterGo**, una API robusta y escalable desarrollada en **Go** y diseñada para funcionar en una arquitectura Serverless.

## 🚀 Despliegue y Arquitectura
La API está construida para ser altamente eficiente y está desplegada utilizando las siguientes herramientas de **AWS**:
* **AWS Lambda**: Funciones serverless para el procesamiento de peticiones.
* **API Gateway**: Punto de entrada para las solicitudes HTTP.
* **AWS S3**: Almacenamiento de objetos para imágenes de perfiles y tweets.
* **MongoDB Atlas**: Base de datos NoSQL en la nube.

## ✨ Funcionalidades principales
* **Gestión de Usuarios**: Registro, login y generación de JWT.
* **Tweets con Multimedia**: Creación, lectura y eliminación de mensajes con soporte para imágenes almacenadas en S3.
* **Relaciones**: Sistema de "Seguir" y "Dejar de seguir" a otros usuarios.
* **Feed Dinámico**: Listado de tweets de usuarios seguidos con paginación.
* **Perfil**: Gestión de datos de usuario y subida de avatares/banners.

## 🛠️ Tecnologías
* **Lenguaje**: Go (Golang)
* **Frameworks/Librerías**: `aws-lambda-go`, `aws-sdk-go` (v2), `mongo-driver`, `jwt-go`.
* **Contenedores**: Docker (para el despliegue en Lambda).

## 🖼️ Gestión de Imágenes
Para optimizar el rendimiento, la API devuelve rutas relativas de las imágenes. El cliente debe utilizar la siguiente URL base para visualizar el contenido de los tweets:
`https://twitter-go-cadesa.s3.us-east-1.amazonaws.com/`

## 🔗 Cliente (Frontend)
La interfaz que consume esta API está disponible en:  
👉 [https://puricalvo.github.io/twitterGo-client](https://puricalvo.github.io/twitterGo-client)
```
---
