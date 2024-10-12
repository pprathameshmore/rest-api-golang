package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/pprathameshmore/rest-api-mongodb-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int16  `json:"status"`
	Error  string `json:"error"`
}

type UserController struct {
	session *mongo.Client
}

func NewUserController(s *mongo.Client) *UserController {
	return &UserController{s}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Something went wrong"))
		return
	}

	u := models.User{}

	if err := uc.session.Database("mongogolang").Collection("users").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&u); err != nil {
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 400, Error: "User not found"}})
			return
		}
	}

	uj, err := json.MarshalIndent(u, "", "    ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: 404, Data: uj})
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var u any
	json.NewDecoder(r.Body).Decode(&u)

	doc, err := uc.session.Database("mongogolang").Collection("users").InsertOne(context.TODO(), &u)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 400, Error: "User not found"}})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: 201, Data: doc})
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 400, Error: "User not found"}})
		return
	}

	deleteUser, deleteErr := uc.session.Database("mongogolang").Collection("users").DeleteOne(context.TODO(), bson.M{"_id": objectId})

	if deleteErr != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 400, Error: "User not found"}})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: 200, Data: deleteUser})
}

func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var u []models.User

	cursor, err := uc.session.Database("mongogolang").Collection("users").Find(context.TODO(), bson.M{})
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 400, Error: "Users not found"}})
		return
	}

	if cursor == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 500, Error: "Cursor is nil"}})
		return
	}

	if err = cursor.All(context.TODO(), &u); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: 200, Data: u})

}

func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	var user map[string]string

	oid, err := primitive.ObjectIDFromHex(id)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 400, Error: "User not found"}})
		return
	}

	json.NewDecoder(r.Body).Decode(&user)

	cursor, err := uc.session.Database("mongogolang").Collection("users").UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": bson.M{"age": user["age"], "gender": user["gender"], "name": user["name"]}})

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 400, Error: "Users not found"}})
		return
	}

	if cursor == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonErrorResponse{Error: &ApiError{Status: 500, Error: "Cursor is nil"}})
		return
	}

	json.NewEncoder(w).Encode(Response{Status: 200, Data: cursor})

}
