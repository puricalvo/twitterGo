package routers

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/puricalvo/twitterGo/awsgo"
	"github.com/puricalvo/twitterGo/models"
)

func mimeFromKey(key string) string {
  ext := strings.ToLower(filepath.Ext(key))
  switch ext {
  case ".jpg", ".jpeg":
    return "image/jpeg"
  case ".png":
    return "image/png"
  case ".webp":
    return "image/webp"
  case ".gif":
    return "image/gif"
  default:
    return "application/octet-stream"
  }
}

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

  data := file.Bytes()
  mime := mimeFromKey(nombre)

  r.CustomResp = &events.APIGatewayProxyResponse{
    StatusCode:      200,
    Body:            base64.StdEncoding.EncodeToString(data),
    IsBase64Encoded: true,
    Headers: map[string]string{
      "Content-Type":                mime, // <-- aquí estaba tu problema
      "Access-Control-Allow-Origin": "*",
    },
  }

  r.Status = 200
  return r
}

/* func ObtenerImagenTweet(ctx context.Context, request events.APIGatewayProxyRequest) models.RespApi {
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
            "Content-Type":                "image/*",
            "Access-Control-Allow-Origin": "*",
        },
    }
    r.Status = 200
    return r
} */