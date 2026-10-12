package routers

import (
	"encoding/base64"
	"net/http"

	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"
)

func ObtenerImagen(uploadType string, request models.APIRequest) models.RespApi {
	r := models.RespApi{Status: 400}

	id := request.QueryStringParameters["id"]
	if id == "" {
		r.Message = "El parámetro ID es obligatorio"
		return r
	}

	perfil, err := bd.BuscoPerfil(id)
	if err != nil {
		r.Status = http.StatusNotFound
		r.Message = "Usuario no encontrado"
		return r
	}

	filename := perfil.Avatar
	if uploadType == "B" {
		filename = perfil.Banner
	}
	if filename == "" {
		r.Status = http.StatusNotFound
		r.Message = "La imagen no existe"
		return r
	}

	image, contentType, err := bd.LeerImagen(filename)
	if err != nil {
		r.Status = http.StatusNotFound
		r.Message = "La imagen no existe en MongoDB"
		return r
	}

	r.Status = http.StatusOK
	r.CustomResp = &models.APIResponse{
		StatusCode:      http.StatusOK,
		Body:            base64.StdEncoding.EncodeToString(image),
		IsBase64Encoded: true,
		Headers:         map[string]string{"Content-Type": contentType},
	}
	return r
}
