package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/puricalvo/twitterGo/bd"
	"github.com/puricalvo/twitterGo/handlers"
	"github.com/puricalvo/twitterGo/models"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	database := os.Getenv("MONGODB_DATABASE")
	jwtSign := os.Getenv("JWT_SIGN")
	if mongoURI == "" || database == "" || jwtSign == "" {
		log.Fatal("MONGODB_URI, MONGODB_DATABASE y JWT_SIGN son obligatorias")
	}

	if err := bd.ConectarBD(mongoURI, database); err != nil {
		log.Fatal("Error conectando con MongoDB: ", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API escuchando en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveAPI(w, r, jwtSign)
	})))
}

func serveAPI(w http.ResponseWriter, r *http.Request, jwtSign string) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/images/") {
		serveImage(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "No se pudo leer la solicitud", http.StatusBadRequest)
		return
	}

	query := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}

	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	request := models.APIRequest{
		Headers:               headers,
		QueryStringParameters: query,
		Body:                  string(body),
	}
	ctx := context.WithValue(r.Context(), models.Key("path"), r.URL.Path)
	ctx = context.WithValue(ctx, models.Key("method"), r.Method)
	ctx = context.WithValue(ctx, models.Key("body"), string(body))
	ctx = context.WithValue(ctx, models.Key("jwtSign"), jwtSign)

	response := handlers.Manejadores(ctx, request)
	status := response.Status
	if status == 0 {
		status = http.StatusInternalServerError
	}

	if response.CustomResp != nil {
		for key, value := range response.CustomResp.Headers {
			w.Header().Set(key, value)
		}
		if response.CustomResp.IsBase64Encoded {
			decoded, err := base64.StdEncoding.DecodeString(response.CustomResp.Body)
			if err != nil {
				http.Error(w, "No se pudo decodificar la respuesta", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(response.CustomResp.StatusCode)
			if _, err := w.Write(decoded); err != nil {
				log.Printf("Error escribiendo respuesta: %v", err)
			}
			return
		}
		if response.CustomResp.StatusCode != 0 {
			status = response.CustomResp.StatusCode
		}
		w.WriteHeader(status)
		if _, err := fmt.Fprint(w, response.CustomResp.Body); err != nil {
			log.Printf("Error escribiendo respuesta: %v", err)
		}
		return
	}

	w.WriteHeader(status)
	if _, err := fmt.Fprint(w, response.Message); err != nil {
		log.Printf("Error escribiendo respuesta: %v", err)
	}
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	filename := strings.TrimPrefix(r.URL.Path, "/images/")
	if filename == "" || strings.Contains(filename, "..") {
		http.NotFound(w, r)
		return
	}
	image, contentType, err := bd.LeerImagen(filename)
	if err != nil {
		log.Printf("No se pudo recuperar la imagen %q de GridFS: %v", filename, err)
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	if _, err := w.Write(image); err != nil {
		log.Printf("Error escribiendo imagen: %v", err)
	}
}

func setCORSHeaders(w http.ResponseWriter) {
	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
