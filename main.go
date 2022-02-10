package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/terrpan/clientdb/internal/controllers"
	"github.com/terrpan/clientdb/internal/dbclient"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

// func Homehandler is dummy func for returning "I'm alive"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I'm alive"))
}

// commonMiddleware is a middleware for setting content type on on all requests
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {

	dbclient.DbConnect()

	r := mux.NewRouter()

	r.Use(commonMiddleware)
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/api/clients", controllers.GetClients).Methods("GET")
	r.HandleFunc("/api/clients/{id}", controllers.GetClientbyId).Methods("GET")
	r.HandleFunc("/api/clients", controllers.AddClient).Methods("POST")
	r.HandleFunc("/api/clients/{id}", controllers.UpdateClient).Methods("PUT")
	r.HandleFunc("/api/clients/{id}", controllers.DeleteClient).Methods("DELETE")
	r.HandleFunc("/api/services", controllers.GetServices).Methods("GET")
	r.HandleFunc("/api/services/{id}", controllers.GetServiceById).Methods("GET")
	r.HandleFunc("/api/services", controllers.AddService).Methods("POST")
	r.HandleFunc("/api/services/{id}", controllers.UpdateService).Methods("PUT")
	r.HandleFunc("/api/services/{id}", controllers.DeleteService).Methods("DELETE")
	r.Handle("/", r)

	// setup the cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		ExposedHeaders:   []string{"Content-Type", "Accept", "X-Total-Count"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler: c.Handler(r),
		Addr:    ":8080",
	}

	log.Fatal(srv.ListenAndServe())
}
