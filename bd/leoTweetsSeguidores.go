package bd

import (
	"context"

	"github.com/puricalvo/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
func LeoTweetsSeguidores(ID string, pagina int) ([]models.DevuelvoTweetsSeguidores, bool) {
    ctx := context.TODO()
    db := MongoCN.Database(DatabaseName)
    col := db.Collection("relacion")

    skip := (pagina - 1) * 20

    // Intentamos convertir el ID a ObjectID por si acaso
    objID, _ := primitive.ObjectIDFromHex(ID)

    condiciones := make([]bson.M, 0)
    
    // CAMBIO AQUÍ: Usamos un $or para que busque tanto si es ObjectID como si es String
    // Esto asegura que encuentre el registro sea cual sea el formato en la BD
    condiciones = append(condiciones, bson.M{
        "$match": bson.M{
            "$or": []bson.M{
                {"usuarioid": ID},    // Busca como string (lo que recibe la función)
                {"usuarioid": objID}, // Busca como objeto de Mongo
            },
        },
    })

    condiciones = append(condiciones, bson.M{
        "$lookup": bson.M{
            "from":         "tweet",
            "localField":   "usuariorelacionid",
            "foreignField": "userid",
            "as":           "tweet",
        }})
    
    condiciones = append(condiciones, bson.M{"$unwind": "$tweet"})
    condiciones = append(condiciones, bson.M{"$sort": bson.M{"tweet.fecha": -1}})
    condiciones = append(condiciones, bson.M{"$skip": skip})
    condiciones = append(condiciones, bson.M{"$limit": 20})

    var result []models.DevuelvoTweetsSeguidores

    cursor, err := col.Aggregate(ctx, condiciones)
    if err != nil {
        return result, false
    }

    // Aquí Go intentará llenar tu modelo con strings. 
    // Si los IDs en Mongo son ObjectIDs, Go suele ser capaz de convertirlos 
    // a string automáticamente durante el cursor.All si el tag bson está bien.
    err = cursor.All(ctx, &result)
    if err != nil {
        return result, false
    }
    return result, true
}
