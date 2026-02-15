package bd

import (
	"context"
	"fmt"
	"time"

	"github.com/puricalvo/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoCN *mongo.Client
var DatabaseName string

func ConectarBD(ctx context.Context) error {
	if MongoCN != nil {
		// Ya conectado, no hacemos nada
		return nil
	}

	user := ctx.Value(models.Key("user")).(string)
	passwd := ctx.Value(models.Key("password")).(string)
	host := ctx.Value(models.Key("host")).(string)
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, passwd, host)

	clientOptions := options.Client().ApplyURI(connStr)

	// Timeout para conexión
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctxTimeout, clientOptions)
	if err != nil {
		fmt.Println("Error conectando a Mongo:", err.Error())
		return err
	}

	err = client.Ping(ctxTimeout, nil)
	if err != nil {
		fmt.Println("Error haciendo ping a Mongo:", err.Error())
		return err
	}

	fmt.Println("Conexión Exitosa con la BD")
	MongoCN = client
	DatabaseName = ctx.Value(models.Key("database")).(string)
	return nil
}

func BaseContectada() bool {
	if MongoCN == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := MongoCN.Ping(ctx, nil)
	return err == nil
}

func BuscoPerfilConContext(ctx context.Context, ID string) (models.Usuario, error) {
	var usuario models.Usuario

	if MongoCN == nil {
		return usuario, fmt.Errorf("Mongo no conectado")
	}

	collection := MongoCN.Database(DatabaseName).Collection("usuarios")

	// 🔹 Convertir ID a ObjectID
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return usuario, fmt.Errorf("ID inválido: %s", ID)
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&usuario)
	if err != nil {
		return usuario, err
	}

	return usuario, nil
}