package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"strings"
	"time"
)

type Assembly struct {
	Item string `json:"item"`
	Quantity int `json:"qty"`
}

func homePage(w http.ResponseWriter, r *http.Request){
	log.Println("homePage",r.Method,r.URL)
	fmt.Fprintf(w,"Welcome to the home page")

}
func returnAssembly(w http.ResponseWriter, r *http.Request){
	log.Println("assembly api",r.Method,r.URL)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w,strings.Join(databases," "))
	coldb := client.Database("widgets").Collection("assembly")
	findOptions := options.Find()
	findOptions.SetLimit(10)
	var results []Assembly
	cur, err := coldb.Find(context.TODO(), bson.D{{}}, findOptions)
	if err !=nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()){
		var elem Assembly
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results,elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Fprintf(w,"Found multiple documents: %+v\n", results)

}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/api",returnAssembly)
	log.Fatal(http.ListenAndServe(":8080", nil))
}



func main() {
	handleRequests()
}