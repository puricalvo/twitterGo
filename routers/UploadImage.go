package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"
)

func UploadImage(ctx context.Context, uploadType string, request events.APIGatewayProxyRequest, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400
	IDUsuario := claim.ID.Hex()

	var filename string
	var usuario models.Usuario

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))

	// Define el tipo de archivo según el tipo de subida (A = Avatar, B = Banner)
	switch uploadType {
	case "A":
		filename = "avatars/" + IDUsuario + ".jpg"
		usuario.Avatar = filename
	case "B":
		filename = "banners/" + IDUsuario + ".jpg"
		usuario.Banner = filename
	
	case "T": // <--- Te faltaba esta línea
        filename = "tweetImage/" + IDUsuario + "_" + time.Now().Format("20060102150405")
    }

	// Obtener el tipo de media de la cabecera Content-Type
	contentType := request.Headers["content-type"]
		if contentType == "" {
			contentType = request.Headers["Content-Type"]
		}

		mediaType, params, err := mime.ParseMediaType(contentType)

	// Si el tipo de media es multipart, seguimos con el procesamiento
	if strings.HasPrefix(mediaType, "multipart/") {

		// Manejo si la solicitud está codificada en base64
		var body []byte
		if request.IsBase64Encoded {
			body, err = base64.StdEncoding.DecodeString(request.Body)
			if err != nil {
				r.Status = 500
				r.Message = err.Error()
				return r
			}
		} else {
			body = []byte(request.Body)
		}

		// Crear el lector multipart
		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		// Si encontramos una parte del archivo
		if err != io.EOF {
			if p.FileName() != "" {

				// Creamos un buffer para leer el archivo
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				// Obtener el tipo de contenido real del archivo
				contentType := p.Header.Get("Content-Type")
				if contentType == "" {
					contentType = "application/octet-stream" // Si no hay tipo de contenido, usar un valor por defecto
				}

				// Cargar la configuración de AWS
				cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				// Crear el cliente de S3
				client := s3.NewFromConfig(cfg)

				// Subir el archivo a S3
				_, err = client.PutObject(ctx, &s3.PutObjectInput{
					Bucket:      bucket,
					Key:         aws.String(filename),
					Body:        buf,
					ContentType: aws.String(contentType), // Usamos el Content-Type correcto
				})

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}
			}
		}

		// Actualizamos el registro del usuario en la base de datos
		status, err := bd.ModificoRegistro(usuario, IDUsuario)
		if err != nil || !status {
			r.Status = 400
			r.Message = "Error al modificar registro del usuario " + err.Error()
			return r
		} else {
    // Si es tipo "T" (Tweet), no llamamos a ModificoRegistro.
    // Simplemente hemos subido el archivo a S3 y ya está.
    fmt.Println("Imagen de tweet subida correctamente a S3, saltando actualización de perfil.")
}

	} else {
		r.Message = "Debe enviar una imagen con el 'Content-Type' multipart en el Header"
		r.Status = 400
		return r
	}

	// Si todo ha ido bien, devolvemos un mensaje de éxito
	r.Status = 200
	r.Message = filename
	return r
}
