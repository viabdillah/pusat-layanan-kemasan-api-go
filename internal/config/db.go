// internal/config/db.go
package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB adalah instance koneksi database kita yang bisa diakses secara global
var DB *mongo.Database

func ConnectDB() {
	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Gagal terhubung ke MongoDB: ", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Gagal melakukan ping ke MongoDB: ", err)
	}

	fmt.Println("MongoDB Terhubung!")
	// 'kemasan-db' adalah nama database kita
	DB = client.Database("kemasan-db")
}
