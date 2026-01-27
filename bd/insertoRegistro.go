package bd

import (
	"context"

	"github.com/puricalvo/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertoRegistro(u models.Usuario) (string, bool, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("usuarios")

	u.Password, _ = EncriptarPassword(u.Password)

	result, err := col.InsertOne(ctx, u)
	if err != nil {
		return  "", false, err
	}

	ObjID := result.InsertedID.(primitive.ObjectID)
	return ObjID.Hex(), true, nil

	/* ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil */
}