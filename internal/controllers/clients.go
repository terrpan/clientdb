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
	SlackChannel   string             `json:"slack_channel,omitempty" bson:"slack_channel,omitempty"`
	WebUrl         string             `json:"web_url,omitempty" bson:"web_url,omitempty"`
	CreatedOn      time.Time          `json:"created_on,omitempty" bson:"created_on,omitempty"`
	ModifiedOn     time.Time          `json:"modified_on,omitempty" bson:"modified_on,omitempty"`
}

type ClientResponse struct {
	ID             primitive.ObjectID               `json:"id" bson:"_id,omitempty"`
	ClientName     string                           `json:"client_name" bson:"client_name"`
	SlackChannel   string                           `json:"slack_channel,omitempty" bson:"slack_channel,omitempty"`
	WebUrl         string                           `json:"web_url,omitempty" bson:"web_url,omitempty"`
	MangedServices []ClientsManagedServicesResponse `json:"managed_services" bson:"managed_services"`
	ClientContacts []ClientsContactResponse         `json:"client_contacts" bson:"client_contacts"`
	CreatedOn      time.Time                        `json:"created_on,omitempty" bson:"created_on,omitempty"`
	ModifiedOn     time.Time                        `json:"modified_on,omitempty" bson:"modified_on,omitempty"`
}

type ClientsManagedServicesResponse struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ServiceName      string             `json:"service_name" bson:"service_name" validate:"required"`
	ServiceType      string             `json:"service_type" bson:"service_type" validate:"required"`
	ServiceStatus    string             `json:"service_status" bson:"service_status" validate:"required"`
	InvoiceFrequency string             `json:"invoice_frequency" bson:"invoice_frequency"`
	InvoiceAmount    float64            `json:"invoice_amount" bson:"invoice_amount"`
	ManagementFee    float64            `json:"management_fee" bson:"management_fee"`
}

type ClientsContactResponse struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName   string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName    string             `json:"last_name" bson:"last_name" validate:"required"`
	FullName    string             `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Email       string             `json:"email" bson:"email" validate:"required,email"`
	PhoneNumber string             `json:"phone_number,omitempty" bson:"phone_number"`
	Role        string             `json:"role,omitempty" bson:"role"`
}

var (
	validate                            = *validator.New()
	clientsCollection *mongo.Collection = dbclient.GetCollection(dbclient.DB, "clients")
)

// getClient returns all clients
func GetClients(w http.ResponseWriter, r *http.Request) {
	clients := []ClientResponse{}

	// join manged_services from services collection using mongoDB's $lookup and $project to get the required fields
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
		{
			"$project": bson.M{
				"_id":                               1,
				"client_name":                       1,
				"slack_channel":                     1,
				"web_url":                           1,
				"created_on":                        1,
				"modified_on":                       1,
				"managed_services._id":              1,
				"managed_services.service_name":     1,
				"managed_services.service_type":     1,
				"managed_services.InvoiceFrequency": 1,
				"managed_services.InvoiceAmount":    1,
				"managed_services.ManagementFee":    1,
				"client_contacts._id":               1,
				"client_contacts.first_name":        1,
				"client_contacts.last_name":         1,
				"client_contacts.full_name":         1,
				"client_contacts.email":             1,
				"client_contacts.phone_number":      1,
				"client_contacts.role":              1,
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
		var client ClientResponse
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
	var client ClientResponse
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// aggregate the client and it's services and contacts using client_id, and use $project to get the required fields
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
		{
			"$project": bson.M{
				"_id":                               1,
				"client_name":                       1,
				"slack_channel":                     1,
				"web_url":                           1,
				"created_on":                        1,
				"modified_on":                       1,
				"managed_services._id":              1,
				"managed_services.service_name":     1,
				"managed_services.service_type":     1,
				"managed_services.InvoiceFrequency": 1,
				"managed_services.InvoiceAmount":    1,
				"managed_services.ManagementFee":    1,
				"client_contacts._id":               1,
				"client_contacts.first_name":        1,
				"client_contacts.last_name":         1,
				"client_contacts.full_name":         1,
				"client_contacts.email":             1,
				"client_contacts.phone_number":      1,
				"client_contacts.role":              1,
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
