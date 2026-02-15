package routers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"
)

func VerPerfil(request events.APIGatewayProxyRequest) models.RespApi {
	var r models.RespApi
	r.Status = 400

	fmt.Println("Entré en VerPerfil")
	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El parámetro ID es obligatorio"
		return r
	}

	perfil, err := bd.BuscoPerfil(ID)
	if err != nil {
		// Diferenciar si no hay documentos vs error real
		if err.Error() == "mongo: no documents in result" {
			r.Status = 404
			r.Message = "Usuario no encontrado"
		} else {
			r.Status = 500
			r.Message = "Ocurrió un error al buscar el usuario: " + err.Error()
		}
		return r
	}

	respJson, err := json.Marshal(perfil)
	if err != nil {
		r.Status = 500
		r.Message = "Error al formatear los datos de los usuarios como JSON: " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = string(respJson)
	return r
}