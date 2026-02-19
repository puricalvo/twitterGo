package bd

import (
	"context"

	"github.com/puricalvo/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuscoPerfil(ID string) (models.Usuario, error) {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("usuarios")

	var perfil models.Usuario
	objID, _ := primitive.ObjectIDFromHex(ID)

	condicion := bson.M{
		"_id": objID,
	}

	err := col.FindOne(ctx, condicion).Decode(&perfil)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			// Usuario sin perfil aún
			return perfil, nil
		}

		// Error real de Mongo
		return perfil, err
	}

	perfil.Password = ""
	return perfil, nil
	}

/* func BuscoPerfilConContext(ctx context.Context, ID string) (models.Usuario, error) {
	var usuario models.Usuario

	if MongoCN == nil {
		return usuario, fmt.Errorf("Mongo no conectado")
	}

	collection := MongoCN.Database(DatabaseName).Collection("usuarios")

	// 🔹 Convertir string a ObjectID
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return usuario, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&usuario)
	if err != nil {
		return usuario, err
	}

	usuario.Password = ""
	return usuario, nil
}
 */

