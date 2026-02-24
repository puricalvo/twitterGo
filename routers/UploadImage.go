package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
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

	// 1. Definimos la ruta y el nombre del archivo en S3
	switch uploadType {
	case "A":
		filename = "avatars/" + IDUsuario + ".jpg"
		usuario.Avatar = filename
	case "B":
		filename = "banners/" + IDUsuario + ".jpg"
		usuario.Banner = filename
	case "T":
		// Añadimos .jpg al final para que el archivo sea válido
		filename = "tweetImage/" + IDUsuario + "_" + time.Now().Format("20060102150405") + ".jpg"
		
	}

	contentType := request.Headers["content-type"]
	if contentType == "" {
		contentType = request.Headers["Content-Type"]
	}

	mediaType, params, err := mime.ParseMediaType(contentType)

	if strings.HasPrefix(mediaType, "multipart/") {
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

		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		if err != io.EOF {
			if p.FileName() != "" {
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				// --- CAMBIO AQUÍ ---
                fileContentType := p.Header.Get("Content-Type")
                
                // Si viene vacío o es el tipo genérico, intentamos detectarlo por el contenido
                if fileContentType == "" || fileContentType == "application/octet-stream" {
                    // Usamos el buffer para detectar el tipo real (jpeg, png, etc)
                    fileContentType = http.DetectContentType(buf.Bytes())
                }

                // Si es un Tweet ("T"), nos aseguramos de que S3 sepa que es una imagen
                if uploadType == "T" && (fileContentType == "application/octet-stream" || fileContentType == "text/plain; charset=utf-8") {
                    fileContentType = "image/jpeg"
                }

				cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				client := s3.NewFromConfig(cfg)

				_, err = client.PutObject(ctx, &s3.PutObjectInput{
					Bucket:      bucket,
					Key:         aws.String(filename),
					Body:        buf,
					ContentType: aws.String(fileContentType),
				})

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}
			}
		}

		// 2. SOLO actualizamos la BD si es Avatar (A) o Banner (B)
		if uploadType == "A" || uploadType == "B" {
			status, err := bd.ModificoRegistro(usuario, IDUsuario)
			if err != nil || !status {
				r.Status = 400
				r.Message = "Error al modificar registro del usuario " + err.Error()
				return r
			}
		} else {
			fmt.Println("Imagen de tweet subida a S3: " + filename)
		}

	} else {
		r.Message = "Debe enviar una imagen con el 'Content-Type' multipart"
		r.Status = 400
		return r
	}

	// 3. Respuesta de éxito: El nombre del archivo va en r.Message
	r.Status = 200
	r.Message = filename
	return r
}