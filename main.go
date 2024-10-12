package main

import (
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pprathameshmore/rest-api-mongodb-go/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	r := httprouter.New()
	uc := controllers.NewUserController(getSession())

	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	r.GET("/user", uc.GetUsers)
	r.PATCH("/user/:id", uc.UpdateUser)

	http.ListenAndServe("localhost:4000", r)
}

func getSession() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI("mongodb://localhost:27017/mongogolang"))

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return client
}
