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

var (
	servicesCollection *mongo.Collection = dbclient.GetCollection(dbclient.DB, "services")
)

type ServiceBase struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ServiceName        string             `json:"service_name" bson:"service_name" validate:"required"`
	ServiceType        string             `json:"service_type" bson:"service_type" validate:"required"`
	ServiceOwner       string             `json:"service_owner" bson:"service_owner" validate:"required"`
	ServiceDescription string             `json:"service_description" bson:"service_description"`
	ServiceStatus      string             `json:"service_status" bson:"service_status" validate:"required"`
	AttachedToClient   []Clients          `json:"attached_to_client" bson:"attached_to_client"`
	InvoiceFrequency   string             `json:"invoice_frequency" bson:"invoice_frequency"`
	InvoiceAmount      float64            `json:"invoice_amount" bson:"invoice_amount"`
	ManagementFee      float64            `json:"management_fee" bson:"management_fee"`
	CreatedOn          time.Time          `json:"created_on" bson:"created_on,omitempty"`
	ModifiedOn         time.Time          `json:"modified_on" bson:"modified_on,omitempty"`
}

type ServiceResponse struct {
	ID                 primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	ServiceName        string                  `json:"service_name" bson:"service_name"`
	ServiceType        string                  `json:"service_type" bson:"service_type"`
	ServiceOwner       string                  `json:"service_owner" bson:"service_owner"`
	ServiceDescription string                  `json:"service_description" bson:"service_description"`
	ServiceStatus      string                  `json:"service_status" bson:"service_status"`
	Client             []ServiceClientResponse `json:"client" bson:"client"`
	InvoiceFrequency   string                  `json:"invoice_frequency" bson:"invoice_frequency"`
	InvoiceAmount      float64                 `json:"invoice_amount" bson:"invoice_amount"`
	ManagementFee      float64                 `json:"management_fee" bson:"management_fee"`
	CreatedOn          time.Time               `json:"created_on" bson:"created_on,omitempty"`
	ModifiedOn         time.Time               `json:"modified_on" bson:"modified_on,omitempty"`
}

type ServiceClientResponse struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	ClientName string `json:"client_name" bson:"client_name"`
}

type Clients struct {
	ClientID primitive.ObjectID `json:"client_id" bson:"_id"`
}

// func GetServices returns all registered services from db
func GetServices(w http.ResponseWriter, r *http.Request) {
	services := []ServiceResponse{}

	// aggregate the services with the clients collection to get the client name and id for each service using mongodb $lookup and $project
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
				"_id":                 1,
				"service_name":        1,
				"service_type":        1,
				"service_owner":       1,
				"service_description": 1,
				"service_status":      1,
				"invoice_frequency":   1,
				"invoice_amount":      1,
				"management_fee":      1,
				"created_on":          1,
				"modified_on":         1,
				"client._id":          1,
				"client.client_name":  1,
			},
		},
	}

	// execute the pipeline
	cursor, err := servicesCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		response := "Failed to get services"
		log.Error(response + err.Error())
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response))
		return
	}

	// iterate through the cursor and add the results to the services slice
	for cursor.Next(context.TODO()) {
		var result ServiceResponse
		cursor.Decode(&result)
		services = append(services, result)
	}

	// Count documents in array and return x-total-count header
	w.Header().Add("X-Total-Count", (strconv.Itoa(len(services))))

	// Return the slice
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)

}

// func GetServiceById returns a single service from db
func GetServiceById(w http.ResponseWriter, r *http.Request) {
	var service ServiceResponse
	// get the id from the url
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// aggregate the services collection with the clients collection to get the client name and client id for each contact
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
				"_id":                 1,
				"service_name":        1,
				"service_type":        1,
				"service_owner":       1,
				"service_description": 1,
				"service_status":      1,
				"invoice_frequency":   1,
				"invoice_amount":      1,
				"management_fee":      1,
				"created_on":          1,
				"modified_on":         1,
				"client._id":          1,
				"client.client_name":  1,
			},
		},
	}

	// execute the pipeline
	cursor, err := servicesCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		response := "Failed to get service"
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response))
		return
	}

	// make sure cursor is not empty
	if !cursor.Next(context.TODO()) {
		response := "Service not found"
		log.Error(response)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(response))
		return
	}

	// Decode the document into the client struct
	err = cursor.Decode(&service)
	if err != nil {
		log.Error("failed to decode cursor", err.Error())
		return
	}

	// Return the service
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}

// func AddService adds a new service to the db
func AddService(w http.ResponseWriter, r *http.Request) {
	var service ServiceBase

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		response := "Invalid request payload: "
		log.Error(response, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// validate the body to ensure all required fields are present
	if validationErr := validate.Struct(service); validationErr != nil {
		response := "Body missing required fields: " + validationErr.Error()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Don't allow duplicate service names
	count, err := servicesCollection.CountDocuments(context.TODO(), bson.M{"service_name": service.ServiceName})
	if err != nil {
		response := "Failed to check if service exists"
		log.Error(response, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if count > 0 {
		response := "Service already exists"
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// set the created on and modified on fields
	service.CreatedOn = time.Now()
	service.ModifiedOn = time.Now()

	// insert the service into the collection
	result, err := servicesCollection.InsertOne(context.TODO(), service)
	if err != nil {
		response := "Failed to insert service: "
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Info("Service created ", result.InsertedID.(primitive.ObjectID).Hex())

	// return the service
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.InsertedID)
}

// func UpdateService updates an existing service
func UpdateService(w http.ResponseWriter, r *http.Request) {
	var service ServiceBase
	// get the id from the url
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		response := "Invalid request payload: " + err.Error()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// validate the body to ensure all required fields are present
	if validationErr := validate.Struct(service); validationErr != nil {
		response := "Body missing required fields: " + validationErr.Error()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// verify that the id exists in the collection
	count, err := servicesCollection.CountDocuments(context.TODO(), bson.M{"_id": id})
	if err != nil {
		response := "Failed to check if service id exists"
		log.Error(response, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if count == 0 {
		response := "Service does not exist, id: " + id.Hex()
		log.Error(response)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// bump the timestamp
	service.ModifiedOn = time.Now()

	// update the service in the collection
	result, err := servicesCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": service})
	if err != nil {
		response := "Failed to update service: " + id.Hex()
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Info("Client updated, id: ", id.Hex())

	var updatedService ServiceBase
	if result.MatchedCount == 1 {
		err := servicesCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&updatedService)
		if err != nil {
			response := "Failed to retrieve updated service"
			log.Error(response, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// return the service
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedService)
}

// func DeleteService removes a registered service in the db
func DeleteService(w http.ResponseWriter, r *http.Request) {
	var service ServiceBase
	// get the id from the url
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	// find the service in the collection
	err := servicesCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&service)
	idString := id.Hex()
	if err != nil {
		response := "Failed to find service: " + idString
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// delete the service from the collection
	_, err = servicesCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		response := "Failed to delete service: " + idString
		log.Error(response + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// return the service
	response := "Client deleted, id: " + idString
	log.Info(response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
