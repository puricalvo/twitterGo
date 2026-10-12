package bd

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoCN *mongo.Client
var DatabaseName string

func ConectarBD(uri, database string) error {
	if uri == "" {
		return fmt.Errorf("MONGODB_URI no está configurada")
	}
	if database == "" {
		return fmt.Errorf("MONGODB_DATABASE no está configurada")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	fmt.Println("Conexión Exitosa con la BD")
	fmt.Printf("Base de datos seleccionada: %s\n", database)
	MongoCN = client
	DatabaseName = database
	return nil
}

func BaseContectada() bool {
	if MongoCN == nil {
		return false
	}
	return MongoCN.Ping(context.TODO(), nil) == nil
}
