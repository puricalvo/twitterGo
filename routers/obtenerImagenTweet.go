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

  nombre := request.QueryStringParameters["nombre"]
  if nombre == "" {
    r.Message = "El parámetro nombre es obligatorio"
    return r
  }

  fmt.Println("Buscando imagen de Tweet: " + nombre)

  svc := s3.NewFromConfig(awsgo.Cfg)

  file, err := downloadFromS3(ctx, svc, nombre)
  if err != nil {
    r.Status = 500
    r.Message = "La imagen no existe en S3: " + err.Error()
    return r
  }

  encoded := base64.StdEncoding.EncodeToString(file.Bytes())
  //data := file.Bytes()

  //mime := mimeFromKey(nombre)

  r.CustomResp = &events.APIGatewayProxyResponse{
    StatusCode:      200,
    Body:         encoded,   //base64.StdEncoding.EncodeToString(data),
    IsBase64Encoded: true,
    Headers: map[string]string{
      "Content-Type":                "image/jpeg", // <-- aquí estaba tu problema
      "Access-Control-Allow-Origin": "*",
    },
  }

  r.Status = 200
  r.Message = "Imagen OK!!!!"
  return r
}

