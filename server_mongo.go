package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Usuario struct {
	Name string `json:"name"`
	Edad int    `json:"edad"`
}

func getUsuarios(w http.ResponseWriter, r *http.Request) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("prueba").Collection("test")

	var user []Usuario
	cur, err := collection.Find(context.Background(), bson.D{})

	if err != nil {
		fmt.Println("Error de find")
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Usuario
		err := cur.Decode(&elem)
		if err != nil {
			panic(err)
		}
		user = append(user, elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Error agregando")
		panic(err)
	}

	cur.Close(context.TODO())

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error parse json")
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println("Error desconectando")
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func postUsuarios(w http.ResponseWriter, r *http.Request) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}

	collection := client.Database("prueba").Collection("test")

	if err != nil {
		panic(err)
	}

	res, err := collection.InsertOne(context.TODO(), r.Body)
	fmt.Printf("Se inserto", res.InsertedID)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")

}

func main() {

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/getUsuarios", getUsuarios).Methods("GET")
	r.HandleFunc("/postUsuarios", postUsuarios).Methods("POST")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("listening http://localhost:8080...")
	server.ListenAndServe()

}
