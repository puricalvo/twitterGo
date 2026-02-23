package models

type Tweet struct {
	Mensaje string `bson:"mensaje" json:"mensaje"`
	Imagen  string `bson:"imagen" json:"imagen,omitempty"`
}