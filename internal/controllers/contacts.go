package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/terrpan/clientdb/internal/dbclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactsBase struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName        string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName         string             `json:"last_name" bson:"last_name" validate:"required"`
	FullName         string             `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Email            string             `json:"email" bson:"email" validate:"required,email"`
	AttachedToClient []Clients          `json:"attached_to_client,omitempty" bson:"attached_to_client"`
	PhoneNumber      string             `json:"phone_number,omitempty" bson:"phone_number"`
	Role             string             `json:"role,omitempty" bson:"role"`
	CreatedOn        time.Time          `json:"created_on" bson:"created_on,omitempty"`
	ModifiedOn       time.Time          `json:"modified_on" bson:"modified_on,omitempty"`
}

type ContactResponse struct {
	ID          primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	FirstName   string                  `json:"first_name" bson:"first_name"`
	LastName    string                  `json:"last_name" bson:"last_name"`
	FullName    string                  `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Email       string                  `json:"email" bson:"email"`
	PhoneNumber string                  `json:"phone_number,omitempty" bson:"phone_number"`
	Client      []ContactClientResponse `json:"client,omitempty" bson:"client"`
	Role        string                  `json:"role,omitempty" bson:"role"`
	CreatedOn   time.Time               `json:"created_on" bson:"created_on,omitempty"`
	ModifiedOn  time.Time               `json:"modified_on" bson:"modified_on,omitempty"`
}

type ContactClientResponse struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ClientName string             `json:"client_name" bson:"client_name"`
}

var (
	contactsCollection *mongo.Collection = dbclient.GetCollection(dbclient.DB, "contacts")
)

// getContacts returns all contacts
func GetContacts(w http.ResponseWriter, r *http.Request) {
	contacts := []ContactResponse{}

	// join the contacts with the clients collection to get the client name and client id for each contact
	// using $lookup and $project
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "clients",
				"localField":   "attached_to_client._id",
				"foreignField": "_id",
				"as":           "client",
			},
		},
		{
			"$project": bson.M{
				"_id":                1,
				"first_name":         1,
				"last_name":          1,
				"full_name":          1,
				"email":              1,
				"phone_number":       1,
				"role":               1,
				"created_on":         1,
				"modified_on":        1,
				"client._id":         1,
				"client.client_name": 1,
			},
		},
	}

	// execute the pipeline
	cursor, err := contactsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		response := "Failed to get contacts"
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response))
		return
	}

	// iterate through the cursor and decode each document into the contacts slice
	for cursor.Next(context.TODO()) {
		var contact ContactResponse
		cursor.Decode(&contact)
		contacts = append(contacts, contact)
	}

	// Count documents in slice and return x-total-count header
	w.Header().Add("X-Total-Count", (strconv.Itoa(len(contacts))))

	// Return the contacts slice
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contacts)

}

// GetContactById returns a contact based on the id.
func GetContactById(w http.ResponseWriter, r *http.Request) {
	var contact ContactResponse

	// retrive the id from the request and convert to ObjectID
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// aggregate the contacts collection with the clients collection to get the client name and client id for each contact
	// using $lookup and $project
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": id},
		},
		{
			"$lookup": bson.M{
				"from":         "clients",
				"localField":   "attached_to_client._id",
				"foreignField": "_id",
				"as":           "client",
			},
		},
		{
			"$project": bson.M{
				"_id":                1,
				"first_name":         1,
				"last_name":          1,
				"full_name":          1,
				"email":              1,
				"phone_number":       1,
				"role":               1,
				"created_on":         1,
				"modified_on":        1,
				"client._id":         1,
				"client.client_name": 1,
			},
		},
	}

	// execute the pipeline
	cursor, err := contactsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		response := "Bad request"
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response))
		return
	}

	// make sure cursor has a next document
	if !cursor.Next(context.TODO()) {
		response := "No contact found with id: " + id.Hex()
		log.Error(response)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(response))
		return
	}

	// Decode the document into the contact struct
	err = cursor.Decode(&contact)
	if err != nil {
		log.Error("failed to decode cursor", err.Error())
		return
	}

	// Return the contact
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contact)

}

// AddContact adds a new contact to the db
func AddContact(w http.ResponseWriter, r *http.Request) {
	var contact ContactsBase

	//validate the request body
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		response := "Invalid request payload: "
		log.Error(response, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// validate the body to ensure all required fields are present
	if validationErr := validate.Struct(contact); validationErr != nil {
		response := "Body missing required fields: " + validationErr.Error()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Merge Firsname and Lastname into FullName
	contact.FullName = contact.FirstName + " " + contact.LastName

	// set the created on and modified on fields
	contact.CreatedOn = time.Now()
	contact.ModifiedOn = time.Now()

	// insert the contact into the db
	result, err := contactsCollection.InsertOne(context.TODO(), contact)
	if err != nil {
		response := "Failed to insert contact: "
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Info("Created contact: ", result.InsertedID.(primitive.ObjectID).Hex())

	// return the id of the new contact
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.InsertedID)

}

// UpdateContact updates an existing contact in the db
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	var contact ContactsBase

	//retrive the id from the request and convert to ObjectID
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		response := "Invalid request payload: " + err.Error()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// validate the body to ensure all required fields are present
	if validationErr := validate.Struct(contact); validationErr != nil {
		response := "Body missing required fields: " + validationErr.Error()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// verify that the id exists in the collection
	count, err := contactsCollection.CountDocuments(context.TODO(), bson.M{"_id": id})
	if err != nil {
		response := "Failed to check if contact id exists"
		log.Error(response, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	if count == 0 {
		response := "Contact doest not exist, id: " + id.Hex()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// set the modified on field
	contact.ModifiedOn = time.Now()

	//update the contact in the colllection
	result, err := contactsCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": contact})
	if err != nil {
		response := "Failed to update contact: "
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Info("Updated contact: ", id.Hex())

	var updatedContact ContactsBase
	if result.MatchedCount == 1 {
		err := contactsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&updatedContact)
		if err != nil {
			response := "Failed to retrieve updated contact"
			log.Error(response, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// return the updated contact
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedContact)

}

// DeleteContact deletes a contact from the collection
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	var contact ContactsBase
	// retreive the id from the request and convert to ObjectID
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// find the contact in the collection based on the id
	err := contactsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&contact)
	idString := id.Hex()
	if err != nil {
		response := "Failed to find contact: " + idString
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// delete the contact from the collection
	_, err = contactsCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		response := "Failed to delete contact: " + idString
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := "Client deleted, id: " + idString
	log.Info(response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}