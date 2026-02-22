
---
```go
# ⚙️ TwitterGo - API (Backend)

Este es el motor de **TwitterGo**, una API robusta y escalable desarrollada en **Go** y diseñada para funcionar en una arquitectura Serverless.

## 🚀 Despliegue y Arquitectura
La API está construida para ser altamente eficiente y está desplegada utilizando las siguientes herramientas de **AWS**:
* **AWS Lambda**: Funciones serverless para el procesamiento de peticiones.
* **API Gateway**: Punto de entrada para las solicitudes HTTP.
* **MongoDB Atlas**: Base de datos NoSQL en la nube.

## ✨ Funcionalidades principales
* **Gestión de Usuarios**: Registro, login y generación de JWT.
* **Tweets**: Creación, lectura y eliminación de mensajes.
* **Relaciones**: Sistema de "Seguir" y "Dejar de seguir" a otros usuarios.
* **Perfil**: Gestión de datos de usuario y subida de avatares/banners.

## 🛠️ Tecnologías
* **Lenguaje**: Go (Golang)
* **Frameworks/Librerías**: `aws-lambda-go`, `mongo-driver`, `jwt-go`.
* **Contenedores**: Docker (para el despliegue en Lambda).

## 🔗 Cliente (Frontend)
Puedes ver la interfaz que consume esta API aquí:  
👉 [https://puricalvo.github.io/twitterGo-client](https://puricalvo.github.io/twitterGo-client)
```
---
