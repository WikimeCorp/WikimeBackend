package db

import (
	"context"
	"log"

	conflib "github.com/WikimeCorp/WikimeBackend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	animeCollection    *mongo.Collection
	genresCollection   *mongo.Collection
	commentsCollection *mongo.Collection
	googleCollection   *mongo.Collection
	vkCollection       *mongo.Collection
	idBaseCollection   *mongo.Collection
	usersCollection    *mongo.Collection
	tokensCollecton    *mongo.Collection
)

var ctx = context.TODO()

func init() {
	config := conflib.Config

	clientOptions := options.Client().ApplyURI(config.MongoURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	wikimeDB := client.Database(config.DataBaseName)

	animeCollection = wikimeDB.Collection("Anime")
	genresCollection = wikimeDB.Collection("Genres")
	commentsCollection = wikimeDB.Collection("Comments")
	googleCollection = wikimeDB.Collection("Google")
	vkCollection = wikimeDB.Collection("Vk")
	idBaseCollection = wikimeDB.Collection("IdBase")
	usersCollection = wikimeDB.Collection("Users")
	tokensCollecton = wikimeDB.Collection("Tokens")
}
