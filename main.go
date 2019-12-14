package main

import (

	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type Tamu struct {
	ID		primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	Name	string				`json:"name,omitempty" bson:"name,omitempty"`
	Email	string				`json:"email,omitempty" bson:"email,omitempty"`
	With	string				`json:"with,omitempty" bson:"with,omitempty"`
}

var client *mongo.Client

func CreateTamuEp(response http.ResponseWriter, request *http.Request)  {
	response.Header().Add("content-type", "application/json")
	var tamu Tamu
	json.NewDecoder(request.Body).Decode(&tamu)
	collection := client.Database("mantenan").Collection("tamu")
	ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	result, _ := collection.InsertOne(ctx, tamu)
	json.NewEncoder(response).Encode(result)
}

func GetSemuaTamuEp(response http.ResponseWriter, request *http.Request)  {
	response.Header().Add("content-type","application/json")
	var semuaTamu []Tamu
	collection := client.Database("mantenan").Collection("tamu")
	ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	res, err := collection.Find(ctx,bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	defer res.Close(ctx)
	for res.Next(ctx)  {
		var tamu Tamu
		res.Decode(&tamu)
		semuaTamu = append(semuaTamu,tamu)
	}
	if err := res.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(semuaTamu)
}

func GetTamuEp(response http.ResponseWriter, request *http.Request)  {
	response.Header().Add("content-type","application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var tamu Tamu
	collection := client.Database("mantenan").Collection("tamu")
	ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	err := collection.FindOne(ctx,Tamu{ID:id}).Decode(&tamu)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(tamu)
}

func main()  {
	fmt.Print("Running on *:8080")
	ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://mantenan:mantenan@db-0-7yvm4.mongodb.net/test?retryWrites=true&w=majority")
	client,_= mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/tamu",CreateTamuEp).Methods("POST")
	router.HandleFunc("/semuatamu",GetSemuaTamuEp).Methods("GET")
	router.HandleFunc("/tamu/{id}",GetTamuEp).Methods("GET")
	http.ListenAndServe(":8080",router)
}