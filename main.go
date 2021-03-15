package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"./helper"
	"./models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Connection mongoDB with helper class
var collection = helper.ConnectDB()

func main() {

	//Init Router
	r := mux.NewRouter()

	// arrange our route
	r.HandleFunc("/api/test", getTests).Methods("GET")
	r.HandleFunc("/api/test/{id}", getTest).Methods("GET")
	r.HandleFunc("/api/test", createTest).Methods("POST")
	r.HandleFunc("/api/test/{id}", updateTest).Methods("PUT")
	r.HandleFunc("/api/test/{id}", deleteTest).Methods("DELETE")

	// set our port address
	log.Fatal(http.ListenAndServe(":8000", r))

	// import "go.mongodb.org/mongo-driver/mongo"

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(
	//    "mongodb+srv://talentjaza:Monday12@shortener.rvf3k.mongodb.net/shortener?retryWrites=true&w=majority",
	// ))
	// if err != nil { log.Fatal(err) }

	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://talentjaza:Monday12@shortener.rvf3k.mongodb.net/shortener?retryWrites=true&w=majority"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Disconnect(ctx)
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(databases)

	// testDatabase := client.Database("test")
	// userCollection := testDatabase.Collection("user")

	// Create Document
	// userCollection.InsertOne(ctx, bson.D{
	// 	{Key: "name", Value: "P'O"},
	// })

	// Read Document
	// findUser, err := userCollection.Find(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var user []bson.M
	// if err = findUser.All(ctx, &user); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(user)

	// Update Document
	// id, _ := primitive.ObjectIDFromHex("604f839ea6558123c6244194")
	// result, err := userCollection.UpdateOne(
	// 	ctx,
	// 	bson.M{"_id": id},
	// 	bson.D{
	// 		{"$set", bson.D{{"name", "Chonlatit Nerd"}}},
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

	// Delete Document
	// result, err := userCollection.DeleteOne(ctx, bson.M{"name": "Chonlatit Nerd"})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
}

func getTests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var tests []models.Test

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var test models.Test
		// & character returns the memory address of the following variable.
		err := cur.Decode(&test) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		tests = append(tests, test)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(tests) // encode similar to serialize process.
}

func getTest(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var test models.Test
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&test)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(test)
}

func createTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var test models.Test

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&test)
	fmt.Println(test)
	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), test)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var test models.Test

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&test)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"user", test.User},
		}},

		// update := bson.D{
		// 	{"$set", bson.D{
		// 		{"user", bson.D{
		// 			{"firstname", test.User.FirstName},
		// 			{"lastname", test.User.LastName},
		// 		}},
		// 	}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&test)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	test.ID = id

	json.NewEncoder(w).Encode(test)
}

func deleteTest(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}
