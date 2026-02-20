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
    
    // 1. Validar error de base de datos
    if err != nil {
        r.Status = 500
        r.Message = "Error al buscar perfil: " + err.Error()
        return r
    }

    // 2. VALIDACIÓN CRÍTICA: ¿El perfil existe realmente?
    // Comprobamos si el ID del perfil devuelto está vacío
    if len(perfil.ID.Hex()) == 0 || perfil.ID.IsZero() {
        r.Status = 400 // O 404
        r.Message = "Usuario no encontrado en la base de datos"
        return r
    }

    // 3. Si llegamos aquí, el perfil existe y es seguro convertirlo
    respJson, err := json.Marshal(perfil)
    if err != nil {
        r.Status = 500
        r.Message = "Error al formatear perfil a JSON: " + err.Error()
        return r
    }

    r.Status = 200
    r.Message = string(respJson)
    return r
}


