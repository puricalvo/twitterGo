package bd

import (
	"context"
	"time"

	"github.com/puricalvo/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BuscoPerfil(ID string) (models.Usuario, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    
    db := MongoCN.Database(DatabaseName)
    col := db.Collection("usuarios")

    var perfil models.Usuario
    
    // 1. VALIDACIÓN DEL HEX: No ignores el error con "_"
    objID, err := primitive.ObjectIDFromHex(ID)
    if err != nil {
        // Si el ID no es válido, devolvemos el error inmediatamente
        // Esto evita que la función siga y explote
        return perfil, err 
    }

    condicion := bson.M{
        "_id": objID,
    }

    err = col.FindOne(ctx, condicion).Decode(&perfil)
    
    // 2. Manejo de resultados
    if err != nil {
        return perfil, err // Aquí devolvemos el error (incluido ErrNoDocuments)
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

