package routers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"
	// 👈 AÑADIR ESTE IMPORT
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
		r.Message = "Ocurrió un error al intentar buscar el registro " + err.Error()
		return r
	}

	respJson, err := json.Marshal(perfil)
	if err != nil {
		r.Status = 500
		r.Message = "Error al formatear los datos de los usuario como JSON " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = string(respJson)
	return r
}

/* func VerPerfil(request events.APIGatewayProxyRequest) models.RespApi {
	var r models.RespApi
	r.Status = 400

	fmt.Println("Entré en VerPerfil")

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El parámetro ID es obligatorio"
		return r
	}

	fmt.Println("Buscando perfil con ID:", ID)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	perfil, err := bd.BuscoPerfilConContext(ctxTimeout, ID)
	if err != nil {

		fmt.Println("Error BuscoPerfil:", err)

		// ✅ Forma correcta de detectar usuario no encontrado
		if err == mongo.ErrNoDocuments {
			r.Status = 404
			r.Message = "Usuario no encontrado"
		} else {
			r.Status = 500
			r.Message = "Ocurrió un error al buscar el usuario: " + err.Error()
		}
		return r
	}
	fmt.Println("Perfil obtenido OK")
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
 */