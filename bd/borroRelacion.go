package bd

import (
	"context"

	"github.com/puricalvo/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
)

func BorroRelacion(t models.Relacion) (bool, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("relacion")

	filtro := bson.M{
		"usuarioid":         t.UsuarioID,
		"usuariorelacionid": t.UsuarioRelacionID,
	}

	_, err := col.DeleteMany(ctx, filtro)
	if err != nil {
		return false, err
	}

	return true, nil
}