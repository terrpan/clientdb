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
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
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

var (
	contactsCollection *mongo.Collection = dbclient.GetCollection(dbclient.DB, "contacts")
)

// getContacts returns all contacts
func GetContacts(w http.ResponseWriter, r *http.Request) {
	contacts := []ContactsBase{}

	cursor, err := contactsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		response := "Failed to find contacts: "
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// iterate over the cursor and add each document to the contacts array
	for cursor.Next(context.TODO()) {
		var contact ContactsBase
		err := cursor.Decode(&contact)
		if err != nil {
			response := "Failed to decode contact: "
			log.Error(response + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		contacts = append(contacts, contact)
	}

	// Count docuemtns in array and return x-total-count header
	w.Header().Add("X-Total-Count", (strconv.Itoa(len(contacts))))

	// Return the contacts array
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contacts)

}

// GetContactById returns a contact based on the id.
func GetContactById(w http.ResponseWriter, r *http.Request) {
	var contact ContactsBase

	// retrive the id from the request and convert to ObjectID
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// find the contact in the collection
	err := contactsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&contact)
	if err != nil {
		response := "Failed to find contact: "
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// return the contact
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

	//update the contact in the cplllection
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
