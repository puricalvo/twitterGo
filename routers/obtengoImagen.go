package routers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"

)

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

	// 🔹 Construimos URL pública directamente
	bucket := ctx.Value(models.Key("bucketName")).(string)
	imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, filename)

	r.CustomResp = &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       imageURL,
		Headers: map[string]string{
			"Content-Type":                     "text/plain",
			"Access-Control-Allow-Origin":      "http://localhost:3000",
			"Access-Control-Allow-Credentials": "true",
		},
	}

	r.Status = 200
	r.Message = "Imagen OK"
	return r
}
