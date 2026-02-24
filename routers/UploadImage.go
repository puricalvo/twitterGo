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

/* func UploadImage(ctx context.Context, uploadType string, request events.APIGatewayProxyRequest, claim models.Claim) models.RespApi {

	var r models.RespApi
	r.Status = 400
	IDUsuario := claim.ID.Hex()

	var filename string
	var usuario models.Usuario

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))
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

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	if strings.HasPrefix(mediaType, "multipart/") {

		body, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
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

				cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}

				client := s3.NewFromConfig(cfg)
				_, err = client.PutObject(ctx, &s3.PutObjectInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   buf,
				})

				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}
			}
		}

		status, err := bd.ModificoRegistro(usuario, IDUsuario)
		if err != nil || !status {
			r.Status = 400
			r.Message = "Error al modificar registro del usuario " + err.Error()
			return r
		}
	} else {
		r.Message = "Debe enviar una imagen con el 'Content-Type' de tipo 'multipart/' en el Header"
		r.Status = 400
		return r
	}

	r.Status = 200
	r.Message = "Image Upload OK !"
	return r
} */



 func UploadImage(ctx context.Context, uploadType string, request events.APIGatewayProxyRequest, claim models.Claim) models.RespApi {
    var r models.RespApi
    r.Status = 400
    IDUsuario := claim.ID.Hex()

    var filename string
    var usuario models.Usuario

    bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))

    switch uploadType {
    case "A":
        filename = "avatars/" + IDUsuario + ".jpg"
        usuario.Avatar = filename
    case "B":
        filename = "banners/" + IDUsuario + ".jpg"
        usuario.Banner = filename
    case "T":
        filename = "tweetImage/" + IDUsuario + "_" + time.Now().Format("20060102150405") + ".jpg"
    }

    contentType := request.Headers["content-type"]
    if contentType == "" {
        contentType = request.Headers["Content-Type"]
    }

    mediaType, params, err := mime.ParseMediaType(contentType)
    if err != nil {
        r.Status = 500
        r.Message = err.Error()
        return r
    }

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

                // --- SOLUCIÓN AL BUFFER ---
                data := buf.Bytes() // Extraemos los bytes reales
                fileContentType := p.Header.Get("Content-Type")
                
                if fileContentType == "" || fileContentType == "application/octet-stream" {
                    fileContentType = http.DetectContentType(data)
                }

                if uploadType == "T" && (fileContentType == "application/octet-stream" || strings.HasPrefix(fileContentType, "text/")) {
                    fileContentType = "image/jpeg"
                }

                cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
                if err != nil {
                    r.Status = 500
                    r.Message = err.Error()
                    return r
                }

                client := s3.NewFromConfig(cfg)

                // IMPORTANTE: Usamos bytes.NewReader(data) para que S3 siempre lea desde el principio
                _, err = client.PutObject(ctx, &s3.PutObjectInput{
                    Bucket:      bucket,
                    Key:         aws.String(filename),
                    Body:        bytes.NewReader(data), 
                    ContentType: aws.String(fileContentType),
                })

                if err != nil {
                    r.Status = 500
                    r.Message = err.Error()
                    return r
                }
            }
        }

        if uploadType == "A" || uploadType == "B" {
            status, err := bd.ModificoRegistro(usuario, IDUsuario)
            if err != nil || !status {
                r.Status = 400
                r.Message = "Error al modificar registro: " + err.Error()
                return r
            }
        }

    } else {
        r.Message = "Debe enviar una imagen multipart"
        r.Status = 400
        return r
    }

    // Respuesta final para API Gateway
    r.Status = 200
    r.Message = filename
    r.CustomResp = &events.APIGatewayProxyResponse{
        StatusCode: 200,
        // Usamos fmt.Sprintf para evitar errores de concatenación
        Body:       fmt.Sprintf(`{"status": 200, "message": "%s"}`, filename),
        Headers: map[string]string{
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*",
        },
    }

    return r
}