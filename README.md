
---
```go
# TwitterGo API (Backend)

API HTTP desarrollada en Go. La base de datos y los archivos se guardan en MongoDB Atlas; las imágenes se almacenan en una colección GridFS de la misma base de datos.

## Requisitos

- Go 1.25.6 o posterior.
- Una base MongoDB accesible desde el equipo (por ejemplo, MongoDB Atlas).
- En Atlas, añadir la IP desde la que se ejecuta la API a la lista de acceso de red.

## Configuración y ejecución local

En PowerShell, desde la carpeta del backend, configura las variables en la terminal actual y arranca el servidor:

```powershell
$env:MONGODB_URI = "mongodb+srv://<usuario>:<contraseña>@<cluster>/?retryWrites=true&w=majority"
$env:MONGODB_DATABASE = "<nombre-de-la-base>"
$env:JWT_SIGN = "<clave-secreta-larga>"
$env:PORT = "8080"
$env:CORS_ORIGIN = "http://localhost:3000"
$env:ALLOW_REGISTRATION = "true"
$env:DEMO_READ_ONLY = "false"
go run .
```

No compartas ni subas `MONGODB_URI` o `JWT_SIGN` al repositorio. Si la contraseña de MongoDB contiene caracteres especiales, deben codificarse para una URI.

`ALLOW_REGISTRATION` controla el registro público. Si se omite, se mantiene activado para conservar el comportamiento actual. Para una instancia de demostración donde las cuentas se creen de forma controlada, configura `ALLOW_REGISTRATION=false` después de crear las cuentas de demo.

`DEMO_READ_ONLY` bloquea las solicitudes que modifican datos, manteniendo disponibles el inicio de sesión y las consultas. Si se omite, se permiten las operaciones habituales. Para una demo pública de solo lectura, configura `DEMO_READ_ONLY=true`.

La API queda disponible en `http://localhost:8080`. El cliente React puede usar ese origen por defecto. Para apuntar a otra dirección al crear el frontend, configura `REACT_APP_API_HOST` antes de ejecutar `npm run dev` o `npm run build`.

Las variables requeridas también están listadas en [`.env.example`](./.env.example). Es una plantilla: copia los nombres y coloca los secretos en el entorno donde se ejecuta la API. No guardes credenciales reales en ese archivo ni lo sustituyas por `api.env`.

## Despliegue en Serv00

La instancia de producción actual se ejecuta como un binario FreeBSD/amd64 en `~/twittergo-api` y se publica mediante el subdominio proxy `api.cadesa.serv00.net` hacia el puerto TCP reservado `6877`.

- Mantén `api.env`, `api.log`, `api.pid` y el binario compilado fuera de Git.
- En Serv00, `api.env` debe definir `MONGODB_URI`, `MONGODB_DATABASE`, `JWT_SIGN`, `PORT=6877` y `CORS_ORIGIN=https://puricalvo.github.io`, con permisos privados (`chmod 600`).
- Para una instancia de demostración separada, usa otro directorio, puerto, dominio y base de datos. En su `api.env`, limita el usuario MongoDB a esa base, configura `ALLOW_REGISTRATION=false` y `DEMO_READ_ONLY=true` una vez creada la cuenta demo. No reutilices la configuración de la instancia personal.
- El cargador opcional [`cmd/seed-demo`](./cmd/seed-demo) prepara perfiles ficticios en español, publicaciones, relaciones e ilustraciones SVG en GridFS, y asigna nombre, biografía, avatar y banner a la cuenta demo. Comprueba que `MONGODB_DATABASE` sea exactamente `twittergo_demo` y necesita `DEMO_USER_EMAIL` con el correo de la cuenta demo existente. Es idempotente; ejecútalo solo en la instancia aislada.
- Para ejecutarlo, compila `go build -trimpath -ldflags="-s -w" -o twittergo-seed-demo-freebsd-amd64 .\cmd\seed-demo` en Windows con `GOOS=freebsd`, `GOARCH=amd64` y `CGO_ENABLED=0`; sube el binario a `~/twittergo-demo-api`. En SSH, entra en esa carpeta, carga `api.env` con `set -a; . ./api.env; set +a`, solicita el correo de la cuenta con `printf 'Correo de la cuenta demo: '; read DEMO_USER_EMAIL; export DEMO_USER_EMAIL` y ejecuta `./twittergo-seed-demo-freebsd-amd64`. El comando no imprime el correo ni secretos; elimina el binario del servidor al terminar.
- Copia [`deploy/serv00/start-api.sh`](./deploy/serv00/start-api.sh) a `~/twittergo-api/start-api.sh` en Serv00. La tarea cron `@reboot` puede llamarlo para arrancar el proceso.
- La configuración del dominio proxy, el puerto reservado y la tarea cron pertenecen a la cuenta Serv00 y deben volver a configurarse allí si se recrea la cuenta.
- Para compilar el binario desde Windows para el servidor, con Go instalado:

```powershell
$env:GOOS = "freebsd"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -trimpath -ldflags="-s -w" -o twittergo-freebsd-amd64 .
Remove-Item Env:GOOS, Env:GOARCH, Env:CGO_ENABLED
```

Sube ese binario a `~/twittergo-api` y reinicia la API según la configuración de Serv00. No copies la URI de Atlas ni la clave JWT al repositorio.

## Imágenes

Los avatares, banners e imágenes de tweets se guardan en GridFS dentro de la base seleccionada en `MONGODB_DATABASE`. Se sirven desde `GET /images/<nombre>`. GridFS evita el límite de 16 MB por documento BSON, pero las imágenes cuentan para el almacenamiento y las transferencias de MongoDB Atlas.

Las imágenes que solo existían en S3 no se copian automáticamente a GridFS. Si aún tienes una copia, tendrás que volver a subirlas; si no, esos archivos antiguos no estarán disponibles.

## Cliente

La interfaz está en [twitterGo-client](https://puricalvo.github.io/twitterGo-client).
```
---
