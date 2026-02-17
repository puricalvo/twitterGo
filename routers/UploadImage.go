package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"strings"

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

	switch uploadType {
	case "A":
		filename = "avatars/" + IDUsuario + ".jpg"
		usuario.Avatar = filename
	case "B":
		filename = "banners/" + IDUsuario + ".jpg"
		usuario.Banner = filename
	}

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	if strings.HasPrefix(mediaType, "multipart/") {

		// Manejar si API Gateway manda base64 o no
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

				// Obtener Content-Type real del archivo
				contentType := p.Header.Get("Content-Type")
				if contentType == "" {
					contentType = "application/octet-stream"
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
					ContentType: aws.String(contentType),
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
		r.Message = "Debe enviar una imagen con el 'Content-Type' multipart en el Header"
		r.Status = 400
		return r
	}

	r.Status = 200
	r.Message = "Image Upload OK !"
	return r
}
