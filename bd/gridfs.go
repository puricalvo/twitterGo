package bd

import (
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GuardarImagen(filename, contentType string, image io.Reader) error {
	bucket, err := gridfs.NewBucket(MongoCN.Database(DatabaseName))
	if err != nil {
		return err
	}

	_, err = bucket.UploadFromStream(filename, image, options.GridFSUpload().SetMetadata(bson.M{
		"contentType": contentType,
	}))
	return err
}

func LeerImagen(filename string) ([]byte, string, error) {
	bucket, err := gridfs.NewBucket(MongoCN.Database(DatabaseName))
	if err != nil {
		return nil, "", err
	}

	stream, err := bucket.OpenDownloadStreamByName(filename)
	if err != nil {
		return nil, "", err
	}
	defer stream.Close()

	contentType, ok := stream.GetFile().Metadata.Lookup("contentType").StringValueOK()
	if !ok || contentType == "" {
		contentType = "image/jpeg"
	}
	image, err := io.ReadAll(stream)
	if err != nil {
		return nil, "", err
	}
	return image, contentType, nil
}
