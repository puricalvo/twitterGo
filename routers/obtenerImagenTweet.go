package routers

import (
	"encoding/base64"
	"net/http"

	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"
)

func ObtenerImagenTweet(request models.APIRequest) models.RespApi {
	r := models.RespApi{Status: 400}
	filename := request.QueryStringParameters["nombre"]
	if filename == "" {
		r.Message = "El parámetro nombre es obligatorio"
		return r
	}

	image, contentType, err := bd.LeerImagen(filename)
	if err != nil {
		r.Status = http.StatusNotFound
		r.Message = "La imagen no existe en MongoDB"
		return r
	}

	r.Status = http.StatusOK
	r.Message = "Imagen OK"
	r.CustomResp = &models.APIResponse{
		StatusCode:      http.StatusOK,
		Body:            base64.StdEncoding.EncodeToString(image),
		IsBase64Encoded: true,
		Headers:         map[string]string{"Content-Type": contentType},
	}
	return r
}
