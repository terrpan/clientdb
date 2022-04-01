package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/terrpan/clientdb/internal/controllers"
	"github.com/terrpan/clientdb/internal/util"
)

func init() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)

	// set log level depending on env
	switch strings.ToLower(util.LogLevel) {
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}

	log.Info("Log level set to: ", log.GetLevel().String())
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

func logger(next http.Handler) http.Handler {
	if log.GetLevel().String() == "debug" {
		return handlers.CombinedLoggingHandler(os.Stdout, next)
	}
	return next
}

func main() {

	util.DbConnect()
	//print all os.envs
	r := mux.NewRouter()

	r.Use(commonMiddleware, logger)
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/api/clients", controllers.GetClients).Methods("GET")
	r.HandleFunc("/api/clients/{id}", controllers.GetClientbyId).Methods("GET")
	r.HandleFunc("/api/clients", controllers.AddClient).Methods("POST")
	r.HandleFunc("/api/clients/{id}", controllers.UpdateClient).Methods("PUT", "PATCH")
	r.HandleFunc("/api/clients/{id}", controllers.DeleteClient).Methods("DELETE")
	r.HandleFunc("/api/services", controllers.GetServices).Methods("GET")
	r.HandleFunc("/api/services/{id}", controllers.GetServiceById).Methods("GET")
	r.HandleFunc("/api/services", controllers.AddService).Methods("POST")
	r.HandleFunc("/api/services/{id}", controllers.UpdateService).Methods("PATCH", "PUT")
	r.HandleFunc("/api/services/{id}", controllers.DeleteService).Methods("DELETE")
	r.HandleFunc("/api/contacts", controllers.GetContacts).Methods("GET")
	r.HandleFunc("/api/contacts/{id}", controllers.GetContactById).Methods("GET")
	r.HandleFunc("/api/contacts", controllers.AddContact).Methods("POST")
	r.HandleFunc("/api/contacts/{id}", controllers.UpdateContact).Methods("PATCH", "PUT")
	r.HandleFunc("/api/contacts/{id}", controllers.DeleteContact).Methods("DELETE")
	r.Handle("/", r)

	// setup the cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "FETCH"},
		ExposedHeaders:   []string{"Content-Type", "Accept", "X-Total-Count"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler: c.Handler(r),
		Addr:    ":8080",
	}

	log.Fatal(srv.ListenAndServe())
}
