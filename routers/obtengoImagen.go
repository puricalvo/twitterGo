package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/puricalvo/twitterGo/awsgo"
	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"

)

// Función de headers CORS
func corsHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                     "image/jpeg",
		"Access-Control-Allow-Origin":      "http://localhost:3000",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Content-Type,Authorization",
		"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE,OPTIONS",
	}
}

// Avatar por defecto en base64 (importado o generado)
var AvatarNoFoundBase64 string // 🔹 Aquí deberías poner tu base64 pre-generado o leerlo desde S3/local

func ObtenerImagen(
	ctx context.Context,
	uploadType string,
	request events.APIGatewayProxyRequest,
	claim models.Claim,
) models.RespApi {

	var r models.RespApi
	r.Status = 400

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El parámetro del ID es obligatorio"
		return r
	}

	perfil, err := bd.BuscoPerfil(ID)
	if err != nil {
		r.Message = "Usuario no encontrado " + err.Error()
		return r
	}

	var filename string
	switch uploadType {
	case "A":
		filename = perfil.Avatar
	case "B":
		filename = perfil.Banner
	}

	if filename == "" {
	r.Status = 404
	r.Message = "El usuario no tiene imagen"
	r.CustomResp = &events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "http://localhost:3000",
			"Access-Control-Allow-Credentials": "true",
		},
	}
	return r
}

	var encoded string

	if filename == "" {
		// Si no hay imagen, usamos la imagen por defecto
		encoded = AvatarNoFoundBase64
	} else {
		svc := s3.NewFromConfig(awsgo.Cfg)
		file, err := downloadFromS3(ctx, svc, filename)
		if err != nil {
			// Si falla S3, también usamos imagen por defecto
			fmt.Println("Error descargando S3, usando avatar por defecto:", err)
			encoded = AvatarNoFoundBase64
		} else {
			encoded = base64.StdEncoding.EncodeToString(file.Bytes())
		}
	}

	r.CustomResp = &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            encoded,
		IsBase64Encoded: true,
		Headers:         corsHeaders(),
	}

	r.Status = 200
	r.Message = "Imagen OK"
	return r
}

func downloadFromS3(ctx context.Context, svc *s3.Client, filename string) (*bytes.Buffer, error) {
	bucket := ctx.Value(models.Key("bucketName")).(string)
	obj, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	file, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(file)
	return buffer, nil
}
