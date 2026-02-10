package bd

import (
	"context"

	"github.com/puricalvo/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertoRelacion(t models.Relacion) (bool, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("relacion")

	// 🔎 Comprobamos si ya existe la relación
	filtro := bson.M{
		"usuarioid":         t.UsuarioID,
		"usuariorelacionid": t.UsuarioRelacionID,
	}

	var resultado models.Relacion
	err := col.FindOne(ctx, filtro).Decode(&resultado)
	if err == nil {
		// Ya existe → no insertamos duplicado
		return false, nil
	}

	_, err = col.InsertOne(ctx, t)
	if err != nil {
		return false, err
	}

	return true, nil
}