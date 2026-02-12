package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
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

	bucket := ctx.Value(models.Key("bucketName")).(string)

	switch uploadType {
	case "A":
		filename = "avatars/" + IDUsuario
		usuario.Avatar = filename
	case "B":
		filename = "banners/" + IDUsuario
		usuario.Banner = filename
	}

	// ==========================
	// OBTENER CONTENT-TYPE SEGURO
	// ==========================

	contentType := request.Headers["Content-Type"]

	if contentType == "" {
		contentType = request.Headers["content-type"]
	}

	if contentType == "" && request.MultiValueHeaders != nil {
		if values, ok := request.MultiValueHeaders["Content-Type"]; ok && len(values) > 0 {
			contentType = values[0]
		}
	}

	if contentType == "" && request.MultiValueHeaders != nil {
		if values, ok := request.MultiValueHeaders["content-type"]; ok && len(values) > 0 {
			contentType = values[0]
		}
	}

	if contentType == "" {
		r.Status = 400
		r.Message = "Content-Type no recibido"
		return r
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		r.Status = 500
		r.Message = "Error parseando Content-Type: " + err.Error()
		return r
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		r.Status = 400
		r.Message = "Debe enviar una imagen con Content-Type multipart/"
		return r
	}

	// ==========================
	// DECODIFICAR BODY
	// ==========================

	var body []byte

	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			r.Status = 500
			r.Message = "Error decodificando base64: " + err.Error()
			return r
		}
		body = decoded
	} else {
		body = []byte(request.Body)
	}

	// ==========================
	// LEER MULTIPART
	// ==========================

	boundary, ok := params["boundary"]
	if !ok {
		r.Status = 400
		r.Message = "Boundary no encontrado en Content-Type"
		return r
	}

	mr := multipart.NewReader(bytes.NewReader(body), boundary)

	var fileUploaded bool

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		if p.FileName() == "" {
			continue
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, p); err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		fileBytes := buf.Bytes()

		// Detectar tipo MIME automáticamente
		detectedContentType := http.DetectContentType(fileBytes)

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			r.Status = 500
			r.Message = "Error al cargar configuración de AWS: " + err.Error()
			return r
		}

		client := s3.NewFromConfig(cfg)
		uploader := manager.NewUploader(client)

		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucket),
			Key:          aws.String(filename),
			Body:         bytes.NewReader(fileBytes),
			ContentType:  aws.String(detectedContentType),
			CacheControl: aws.String("no-cache"),
		})

		if err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		fileUploaded = true
		break
	}

	if !fileUploaded {
		r.Status = 400
		r.Message = "No se encontró archivo en el form-data"
		return r
	}

	// ==========================
	// ACTUALIZAR BASE DE DATOS
	// ==========================

	status, err := bd.ModificoRegistro(usuario, IDUsuario)
	if err != nil || !status {
		r.Status = 400
		r.Message = "Error actualizando base de datos"
		return r
	}

	r.Status = 200
	r.Message = "Image Upload OK!"
	return r
}
