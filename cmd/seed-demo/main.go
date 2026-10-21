package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const demoDatabase = "twittergo_demo"

type demoProfile struct {
	email     string
	firstName string
	lastName  string
	bio       string
	location  string
	avatar    string
	avatarSVG string
	banner    string
	bannerSVG string
}

type demoPost struct {
	key     string
	author  string
	message string
	age     time.Duration
}

var profiles = []demoProfile{
	{
		email:     "ana.rios@twittergo.invalid",
		firstName: "Ana",
		lastName:  "Ríos",
		bio:       "Diseño, cafeterías tranquilas y paseos junto al mar. Compartiendo pequeñas ideas de cada día.",
		location:  "Valencia",
		avatar:    "avatars/demo-ana.svg",
		avatarSVG: avatarSVG("#e88b67", "#f4d6b8", "#49394a"),
		banner:    "banners/demo-ana.svg",
		bannerSVG: bannerSVG("#f6c59e", "#d88178", "#f7e7cc"),
	},
	{
		email:     "mateo.soler@twittergo.invalid",
		firstName: "Mateo",
		lastName:  "Soler",
		bio:       "Me gusta aprender cosas nuevas, leer con música y descubrir rincones de la ciudad.",
		location:  "Barcelona",
		avatar:    "avatars/demo-mateo.svg",
		avatarSVG: avatarSVG("#6d9eaa", "#e8c6a3", "#34465c"),
		banner:    "banners/demo-mateo.svg",
		bannerSVG: bannerSVG("#a8d1cc", "#5e91a2", "#e3efd9"),
	},
	{
		email:     "lucia.vega@twittergo.invalid",
		firstName: "Lucía",
		lastName:  "Vega",
		bio:       "Fotografía, plantas y recetas improvisadas. Siempre buscando una buena historia.",
		location:  "Sevilla",
		avatar:    "avatars/demo-lucia.svg",
		avatarSVG: avatarSVG("#b37caa", "#e8bd9a", "#3f3448"),
		banner:    "banners/demo-lucia.svg",
		bannerSVG: bannerSVG("#e2b4cf", "#9d79ad", "#f4dfbd"),
	},
}

var posts = []demoPost{
	{key: "ana-1", author: "ana.rios@twittergo.invalid", message: "Un paseo sin prisa y una libreta en el bolso: a veces las mejores ideas aparecen cuando dejamos de buscarlas. ☀️", age: 8 * time.Hour},
	{key: "ana-2", author: "ana.rios@twittergo.invalid", message: "Pequeño recordatorio: descansar también cuenta como avanzar.", age: 30 * time.Hour},
	{key: "ana-3", author: "ana.rios@twittergo.invalid", message: "Hoy probé una cafetería nueva y ya tengo sitio favorito para mis mañanas de lectura. ☕", age: 60 * time.Hour},
	{key: "mateo-1", author: "mateo.soler@twittergo.invalid", message: "Terminé un libro que llevaba meses esperándome en la estantería. Qué gusto volver a leer por puro placer. 📚", age: 5 * time.Hour},
	{key: "mateo-2", author: "mateo.soler@twittergo.invalid", message: "Aprender algo nuevo cada semana: ese es mi reto pequeño para este año.", age: 26 * time.Hour},
	{key: "mateo-3", author: "mateo.soler@twittergo.invalid", message: "La ciudad cambia mucho cuando la recorres a pie y miras hacia arriba.", age: 55 * time.Hour},
	{key: "lucia-1", author: "lucia.vega@twittergo.invalid", message: "Mis plantas han sobrevivido a otra semana de cuidados y conversaciones. Creo que ya somos un equipo. 🌱", age: 3 * time.Hour},
	{key: "lucia-2", author: "lucia.vega@twittergo.invalid", message: "Receta improvisada de hoy: verduras de temporada, aceite de oliva y cero ganas de complicarme. Quedó buenísima.", age: 22 * time.Hour},
	{key: "lucia-3", author: "lucia.vega@twittergo.invalid", message: "La luz de última hora convierte hasta la calle de siempre en un lugar distinto. 📷", age: 48 * time.Hour},
}

func main() {
	if err := seed(); err != nil {
		log.Fatalf("No se pudo preparar el contenido de demo: %v", err)
	}
}

func seed() error {
	uri := os.Getenv("MONGODB_URI")
	databaseName := os.Getenv("MONGODB_DATABASE")
	demoEmail := strings.ToLower(strings.TrimSpace(os.Getenv("DEMO_USER_EMAIL")))
	if uri == "" || databaseName != demoDatabase || demoEmail == "" {
		return fmt.Errorf("se requiere MONGODB_URI, MONGODB_DATABASE=%s y DEMO_USER_EMAIL", demoDatabase)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	db := client.Database(demoDatabase)
	users := db.Collection("usuarios")
	seedUser, err := findUser(ctx, users, demoEmail)
	if err != nil {
		return fmt.Errorf("no se encontró la cuenta demo indicada")
	}

	profileIDs := make(map[string]string, len(profiles))
	for _, profile := range profiles {
		id, err := ensureProfile(ctx, users, profile)
		if err != nil {
			return err
		}
		profileIDs[profile.email] = id
	}

	for _, profile := range profiles {
		if err := ensureImage(ctx, db, profile.avatar, profile.avatarSVG); err != nil {
			return err
		}
		if err := ensureImage(ctx, db, profile.banner, profile.bannerSVG); err != nil {
			return err
		}
	}

	tweets := db.Collection("tweet")
	for _, post := range posts {
		authorID, ok := profileIDs[post.author]
		if !ok {
			return fmt.Errorf("autor de publicación demo no configurado")
		}
		if err := ensurePost(ctx, tweets, post, authorID); err != nil {
			return err
		}
	}

	if err := ensurePost(ctx, tweets, demoPost{
		key:     "welcome-" + seedUser.ID.Hex(),
		message: "¡Bienvenidos a mi demo de TwitterGo! Este espacio contiene perfiles y publicaciones de ejemplo. Puedes explorar la interfaz en modo de solo lectura.",
		age:     2 * time.Hour,
	}, seedUser.ID.Hex()); err != nil {
		return err
	}

	relations := db.Collection("relacion")
	for _, id := range profileIDs {
		if err := ensureRelation(ctx, relations, seedUser.ID.Hex(), id); err != nil {
			return err
		}
	}
	for _, follower := range profileIDs {
		for _, followed := range profileIDs {
			if follower != followed {
				if err := ensureRelation(ctx, relations, follower, followed); err != nil {
					return err
				}
			}
		}
	}

	fmt.Printf("Contenido de demo listo: %d perfiles ficticios, %d publicaciones de ejemplo e ilustraciones originales.\n", len(profiles), len(posts)+1)
	return nil
}

func findUser(ctx context.Context, users *mongo.Collection, email string) (struct {
	ID primitive.ObjectID `bson:"_id"`
}, error) {
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user, err
}

func ensureProfile(ctx context.Context, users *mongo.Collection, profile demoProfile) (string, error) {
	passwordBytes := make([]byte, 32)
	if _, err := rand.Read(passwordBytes); err != nil {
		return "", err
	}
	passwordHash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	fields := bson.M{
		"nombre":          profile.firstName,
		"apellidos":       profile.lastName,
		"fechaNacimiento": now.AddDate(-28, 0, 0),
		"email":           profile.email,
		"password":        string(passwordHash),
		"avatar":          profile.avatar,
		"banner":          profile.banner,
		"biografia":       profile.bio,
		"ubicacion":       profile.location,
		"sitioweb":        "",
		"demoSeed":        true,
	}
	_, err = users.UpdateOne(ctx, bson.M{"email": profile.email}, bson.M{"$setOnInsert": fields}, options.Update().SetUpsert(true))
	if err != nil {
		return "", err
	}

	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := users.FindOne(ctx, bson.M{"email": profile.email}).Decode(&user); err != nil {
		return "", err
	}
	return user.ID.Hex(), nil
}

func ensurePost(ctx context.Context, tweets *mongo.Collection, post demoPost, authorID string) error {
	_, err := tweets.UpdateOne(ctx, bson.M{"demoSeedKey": post.key}, bson.M{
		"$setOnInsert": bson.M{
			"userid":      authorID,
			"mensaje":     post.message,
			"fecha":       time.Now().UTC().Add(-post.age),
			"imagen":      "",
			"demoSeedKey": post.key,
		},
	}, options.Update().SetUpsert(true))
	return err
}

func ensureRelation(ctx context.Context, relations *mongo.Collection, followerID, followedID string) error {
	_, err := relations.UpdateOne(ctx, bson.M{
		"usuarioid":         followerID,
		"usuariorelacionid": followedID,
	}, bson.M{"$setOnInsert": bson.M{
		"usuarioid":         followerID,
		"usuariorelacionid": followedID,
	}}, options.Update().SetUpsert(true))
	return err
}

func ensureImage(ctx context.Context, db *mongo.Database, filename, svg string) error {
	files := db.Collection("fs.files")
	count, err := files.CountDocuments(ctx, bson.M{"filename": filename})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return err
	}
	_, err = bucket.UploadFromStream(filename, strings.NewReader(svg), options.GridFSUpload().SetMetadata(bson.M{
		"contentType": "image/svg+xml",
	}))
	return err
}

func avatarSVG(background, skin, hair string) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 160 160"><rect width="160" height="160" rx="80" fill="%s"/><path d="M26 160c5-34 25-52 54-52s49 18 54 52" fill="#f8f3ec"/><path d="M48 75c0-29 13-47 34-47 23 0 35 20 34 48-1 27-16 42-36 42S48 102 48 75" fill="%s"/><path d="M47 72c-4-31 12-53 37-53 25 0 41 17 38 51-9-10-16-21-20-31-12 15-30 25-55 33" fill="%s"/><circle cx="69" cy="78" r="3" fill="#493c39"/><circle cx="96" cy="78" r="3" fill="#493c39"/><path d="M73 94c6 5 13 5 19 0" fill="none" stroke="#a75e62" stroke-width="3" stroke-linecap="round"/></svg>`, background, skin, hair)
}

func bannerSVG(start, end, accent string) string {
	sum := sha256.Sum256([]byte(start + end + accent))
	id := fmt.Sprintf("g%x", sum[:4])
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 900 260"><defs><linearGradient id="%s" x1="0" y1="0" x2="1" y2="1"><stop stop-color="%s"/><stop offset="1" stop-color="%s"/></linearGradient></defs><rect width="900" height="260" fill="url(#%s)"/><circle cx="720" cy="42" r="150" fill="%s" opacity=".38"/><circle cx="825" cy="226" r="190" fill="%s" opacity=".28"/><path d="M0 220c170-100 280 80 470-18s265-93 430-32v90H0z" fill="%s" opacity=".48"/></svg>`, id, start, end, id, accent, start, accent)
}
