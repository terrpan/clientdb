package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/terrpan/clientdb/internal/dbclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientBase struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ClientName     string             `json:"client_name" bson:"client_name" validate:"required"`
	ClientContacts []ContactsBase     `json:"client_contacts" bson:"client_contacts"` // import ContactsBase from contacts.go
	SlackChannel   string             `json:"slack_channel,omitempty" bson:"slack_channel,omitempty"`
	WebUrl         string             `json:"web_url,omitempty" bson:"web_url,omitempty"`
	MangedServices []ServiceBase      `json:"managed_services" bson:"managed_services"` // import ServiceBase from services.go
	CreatedOn      time.Time          `json:"created_on,omitempty" bson:"created_on,omitempty"`
	ModifiedOn     time.Time          `json:"modified_on,omitempty" bson:"modified_on,omitempty"`
}

var (
	validate                            = *validator.New()
	clientsCollection *mongo.Collection = dbclient.GetCollection(dbclient.DB, "clients")
)

// getClient returns all clients
func GetClients(w http.ResponseWriter, r *http.Request) {
	clients := []ClientBase{}

	// join manged_services from services collection using mongoDB's $lookup
	// https://docs.mongodb.com/manual/reference/operator/aggregation/lookup/
	// https://docs.mongodb.com/manual/core/aggregation-pipeline/
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "services",
				"localField":   "_id",
				"foreignField": "attached_to_client._id",
				"as":           "managed_services",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "contacts",
				"localField":   "_id",
				"foreignField": "attached_to_client._id",
				"as":           "client_contacts",
			},
		},
	}

	// execute the pipeline
	cursor, err := clientsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		response := "Failed to find clients: "
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// iterate through the cursor and add each client to the clients array
	for cursor.Next(context.TODO()) {
		var client ClientBase
		// decode the document into the client struct
		cursor.Decode(&client)
		clients = append(clients, client)
	}

	// count and return the number of services
	w.Header().Add("X-Total-Count", (strconv.Itoa(len(clients))))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clients)

}

// getClientbyId returns a client by id
func GetClientbyId(w http.ResponseWriter, r *http.Request) {
	var client ClientBase
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// aggregate the client and its services and match the client id using $match and $lookup
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": id},
		},
		{
			"$lookup": bson.M{
				"from":         "services",
				"localField":   "_id",
				"foreignField": "attached_to_client._id",
				"as":           "managed_services",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "contacts",
				"localField":   "_id",
				"foreignField": "attached_to_client._id",
				"as":           "client_contacts",
			},
		},
	}

	// execute the pipeline
	cursor, err := clientsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		response := "Bad request"
		log.Error(response, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// count lenth of the cursor and return 404 if no client is found
	if !cursor.Next(context.TODO()) {
		response := "No client not found with id: " + id.Hex()
		log.Error(response)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Decode the document into the client struct
	err = cursor.Decode(&client)
	if err != nil {
		log.Error("failed to decode cursor", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)

}

// updateClient updates a client
func UpdateClient(w http.ResponseWriter, r *http.Request) {
	var client ClientBase
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		log.Error("Error converting id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate the incoming json data
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		response := "Invalid request payload"
		log.Error(response, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// validate the body to ensure all required fields are present
	if validationErr := validate.Struct(client); validationErr != nil {
		response := "Body missing required fields: " + validationErr.Error()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// verify that id exists in the collection
	count, err := clientsCollection.CountDocuments(context.TODO(), bson.M{"_id": id})
	if err != nil {
		response := "Failed to check if client id exists"
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if count == 0 {
		response := "No client not found with id: " + id.Hex()
		log.Error(response)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	client.ModifiedOn = time.Now()

	// Update the client in the collection
	result, err := clientsCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": client})
	if err != nil {
		response := "Failed to update client"
		log.Error(response, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Info("Client updated, id: ", id.Hex())

	// retrieve the updated client
	var updatedClient ClientBase
	if result.MatchedCount == 1 {
		err := clientsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&updatedClient)
		if err != nil {
			response := "Failed to retrieve updated client"
			log.Error(response, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedClient)
}

// deleteClient deletes a client
func DeleteClient(w http.ResponseWriter, r *http.Request) {
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// Check if client exists
	count, err := clientsCollection.CountDocuments(context.TODO(), bson.M{"_id": id})
	if err != nil {
		log.Error("Error checking that client id exists:", err)
		return
	}

	idString := id.Hex()
	if count == 0 {
		response := "Client not found, id: " + idString
		log.Warn(response)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Delete the client from the collection based on id
	_, err = clientsCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		response := "Failed to delete client"
		log.Error(response, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		return
	}

	response := "Client deleted, id: " + idString
	log.Info(response)
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)

}

// addClient adds a new client to the database
func AddClient(w http.ResponseWriter, r *http.Request) {
	var client ClientBase

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		response := "Invalid request payload"
		log.Error(response, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// validate the body to ensure all required fields are present
	if validationErr := validate.Struct(client); validationErr != nil {
		response := "Body missing required fields: " + validationErr.Error()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Don't allow duplicate client names
	count, err := clientsCollection.CountDocuments(context.TODO(), bson.M{"client_name": client.ClientName})
	if err != nil {
		response := "Failed to check if client exists"
		log.Error(response, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if count > 0 {
		response := "Client already exists"
		log.Error(response, client.ClientName)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set the createdOn and modifiedOn fields
	client.CreatedOn = time.Now()
	client.ModifiedOn = time.Now()

	// Insert the new client to collection
	result, err := clientsCollection.InsertOne(context.TODO(), client)
	if err != nil {
		response := "Failed to insert client"
		log.Error(response, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Info("Client added, id:", result.InsertedID.(primitive.ObjectID).Hex())

	// return the id of the new client and 201 status
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.InsertedID)

}

