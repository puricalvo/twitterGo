package routers

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"
)

func ListaUsuarios(request events.APIGatewayProxyRequest, claim models.Claim) models.RespApi {
	var r models.RespApi
	r.Status = 400

	page := request.QueryStringParameters["page"]
	typeUser := request.QueryStringParameters["type"]
	search := request.QueryStringParameters["search"]
	IDUsuario := claim.ID.Hex()

	if len(page) == 0 {
		page = "1"
	}

	pagTemp, err := strconv.Atoi(page)
	if err != nil {
		r.Message = "Debe enviar el parámetro 'page' como entero mayor a 0 " + err.Error()
		return r
	}

	// 🔹 Creamos contexto con timeout por request
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 🔹 Pasamos el contexto a la función BD
	usuarios, status := bd.LeoUsuariosTodosConContext(ctxTimeout, IDUsuario, int64(pagTemp), search, typeUser)
	if !status {
		// 🔹 Si hay error en bd, devolvemos array vacío en vez de romper
		usuarios = []*models.Usuario{}
	}

	// 🔹 Siempre devolvemos un array JSON
	if usuarios == nil {
		usuarios = []*models.Usuario{}
	}

	respJson, err := json.Marshal(usuarios)
	if err != nil {
		r.Status = 500
		r.Message = "Error al formatear los datos de los usuarios en JSON"
		return r
	}

	r.Status = 200
	r.Message = string(respJson)
	return r
}