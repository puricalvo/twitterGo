package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/puricalvo/twitterGo/awsgo"
	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/handlers"
	"github.com/puricalvo/twitterGo/models"
	"github.com/puricalvo/twitterGo/secretmanager"
)

func main() {
	dt := "1970-06-30T00:00:00+00:00" // Ajusta aquí
	t, err := time.Parse("2006-01-02T15:04:05-07:00", dt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(t)

	lambda.Start(EjecutoLambda)
}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAWS()

	// ===== CORS: manejar OPTIONS (preflight) =====
	if request.HTTPMethod == "OPTIONS" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    corsHeaders(),
			Body:       "",
		}, nil
	}
	// =============================================

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno, debe de incluir 'SecretName', 'BucketName', 'UrlPrefix' ",
			Headers:    corsHeaders(),
		}
		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura de Secret " + err.Error(),
			Headers:    corsHeaders(),
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["twittergo"], os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	// Chequeo Conexión a la Base de Datos o Conecto a la Base de Datos
	err = bd.ConectarBD(awsgo.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error conectando la BD " + err.Error(),
			Headers:    corsHeaders(),
		}
		return res, nil
	}

	respAPI := handlers.Manejadores(awsgo.Ctx, request)

	if respAPI.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       respAPI.Message,
			Headers:    corsHeaders(),
		}
		return res, nil
	} else {
		// ===== Aseguramos CORS en CustomResp =====
		if respAPI.CustomResp.Headers == nil {
			respAPI.CustomResp.Headers = map[string]string{}
		}
		for k, v := range corsHeaders() {
			respAPI.CustomResp.Headers[k] = v
		}
		return respAPI.CustomResp, nil
	}
}

func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		fmt.Println("Database leido del secret:", bd.MongoCN)
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}

	return traeParametro
}

// ===== Función de headers CORS =====
func corsHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      "*",
    	"Access-Control-Allow-Headers":     "Content-Type,Authorization",
    	"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE,OPTIONS",
	}
}
