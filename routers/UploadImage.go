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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"

)

func UploadImage(
	ctx context.Context,
	uploadType string,
	request events.APIGatewayProxyRequest,
	claim models.Claim,
) models.RespApi {

	var r models.RespApi
	r.Status = 400

	IDUsuario := claim.ID.Hex()

	var filename string
	var usuario models.Usuario

	// bucket desde context
	bucket := ctx.Value(models.Key("bucketName")).(string)

	switch uploadType {
	case "A":
		filename = "avatars/" + IDUsuario + ".jpg"
		usuario.Avatar = filename
	case "B":
		filename = "banners/" + IDUsuario + ".jpg"
		usuario.Banner = filename
	}

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	fmt.Println("Headers:", request.Headers)
	fmt.Println("IsBase64Encoded:", request.IsBase64Encoded)
	fmt.Println("Body length:", len(request.Body))

	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		r.Status = 400
		r.Message = "Debe enviar una imagen con Content-Type multipart/"
		return r
	}

	


	/* var body []byte

	if request.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(request.Body)
	} else {
		body = []byte(request.Body)
	} */

	 // API Gateway envía el body en base64
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

	if p.FileName() != "" {

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, p); err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		// 🔹 AWS SDK v2 (Lambda)
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		client := s3.NewFromConfig(cfg)
		uploader := manager.NewUploader(client)

		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
			Body:   buf,
		})

		if err != nil {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		
	}

	// Solo actualizamos los campos que tengan valor
	registro := make(map[string]interface{})
	if len(usuario.Avatar) > 0 {
		registro["avatar"] = usuario.Avatar
	}
	if len(usuario.Banner) > 0 {
		registro["banner"] = usuario.Banner
	}

	// Llamada a bd.ModificoRegistro con el struct parcial
	status, err := bd.ModificoRegistro(usuario, IDUsuario)
	if err != nil || !status {
		r.Status = 400
		r.Message = "Error al modificar registro del usuario " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = "Image Upload OK!"
	return r
}