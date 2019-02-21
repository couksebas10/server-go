package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type Usuario struct {
	Id   string `json:"_id"`
	Name string `json:"name"`
	Edad int    `json:"edad"`
}

func getUsuarios(w http.ResponseWriter, r *http.Request) {

	MongoConnURL := "mongodb://127.0.0.1:27017"
	client, err := mongo.Connect(context.TODO(), MongoConnURL, nil)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	fmt.Printf("Connected to MongoDB!")

	collection := client.Database("prueba").Collection("test")

	findOptions := options.Find()
	var user []Usuario

	response, err := collection.Find(context.TODO(), nil, findOptions)
	if err != nil {
		panic(err)
	}

	for response.Next(context.TODO()) {
		var elem Usuario
		err := response.Decode(&elem)
		if err != nil {
			panic(err)
		}
		user = append(user, elem)
	}

	if err := response.Err(); err != nil {
		panic(err)
	}

	response.Close(context.TODO())

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func main() {

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/getUsuarios", getUsuarios).Methods("GET")

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
