package routers

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"strings"
	"time"

	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/models"
)

func UploadImage(uploadType string, request models.APIRequest, claim models.Claim) models.RespApi {
	r := models.RespApi{Status: 400}
	userID := claim.ID.Hex()

	var filename string
	var usuario models.Usuario
	switch uploadType {
	case "A":
		filename = "avatars/" + userID + ".jpg"
		usuario.Avatar = filename
	case "B":
		filename = "banners/" + userID + ".jpg"
		usuario.Banner = filename
	case "T":
		filename = "tweetImage/" + userID + "_" + time.Now().Format("20060102150405") + ".jpg"
	default:
		r.Message = "Tipo de imagen no válido"
		return r
	}

	contentType := request.Headers["Content-Type"]
	if contentType == "" {
		contentType = request.Headers["content-type"]
	}
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
		r.Message = "La solicitud debe ser multipart/form-data"
		return r
	}

	body := []byte(request.Body)
	if request.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			r.Status = 400
			r.Message = "El cuerpo de la solicitud no es base64 válido"
			return r
		}
	}

	reader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	part, err := reader.NextPart()
	if err != nil {
		r.Message = "No se encontró el archivo en la solicitud"
		return r
	}
	defer part.Close()
	if part.FileName() == "" {
		r.Message = "No se encontró el archivo en la solicitud"
		return r
	}

	image, err := io.ReadAll(part)
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}
	imageType := part.Header.Get("Content-Type")
	if imageType == "" {
		imageType = "application/octet-stream"
	}
	if err := bd.GuardarImagen(filename, imageType, bytes.NewReader(image)); err != nil {
		r.Status = 500
		r.Message = "Error guardando imagen en MongoDB: " + err.Error()
		return r
	}

	if uploadType == "A" || uploadType == "B" {
		status, err := bd.ModificoRegistro(usuario, userID)
		if err != nil || !status {
			r.Status = 500
			r.Message = "Error al modificar registro del usuario"
			if err != nil {
				r.Message += ": " + err.Error()
			}
			return r
		}
	}

	r.Status = 200
	if uploadType == "T" {
		r.Message = filename
	} else {
		r.Message = "Image Upload OK !"
	}
	return r
}
