package routers

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/puricalvo/twitterGo/awsgo"
	"github.com/puricalvo/twitterGo/models"
)

func ObtenerImagenTweet(ctx context.Context, request events.APIGatewayProxyRequest) models.RespApi {
    var r models.RespApi
    r.Status = 400

    // Extraemos el nombre de la imagen de los parámetros de la URL
    // Ejemplo: .../obtenerImagenTweet?nombre=tweetImage/123_456.jpg
    nombre := request.QueryStringParameters["nombre"]
    if len(nombre) < 1 {
        r.Message = "El parámetro nombre es obligatorio"
        return r
    }

    fmt.Println("Buscando imagen de Tweet: " + nombre)
    
    // Usamos el cliente de S3
    svc := s3.NewFromConfig(awsgo.Cfg)

    // Reutilizamos tu función downloadFromS3 que ya funciona perfectamente
    file, err := downloadFromS3(ctx, svc, nombre)
    if err != nil {
        r.Status = 500
        r.Message = "La imagen no existe en S3: " + err.Error()
        return r
    }

    // Convertimos a base64 para la respuesta
    encoded := base64.StdEncoding.EncodeToString(file.Bytes())

    r.CustomResp = &events.APIGatewayProxyResponse{
        StatusCode:      200,
        Body:            encoded,
        IsBase64Encoded: true,
        Headers: map[string]string{
            "Content-Type":                "image/jpeg",
            "Access-Control-Allow-Origin": "*",
        },
    }
    r.Status = 200
    return r
}