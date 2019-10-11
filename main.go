package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/couchbase/gocb"
	"github.com/gorilla/mux"
)

var (
	bucket *gocb.Bucket
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving")

	w.Write([]byte("Welcome to the Pancake Stack"))
}

func cageCards(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting a cage card...")

	vars := mux.Vars(r)

	var result interface{}
	query := gocb.NewN1qlQuery("SELECT * FROM philbucket WHERE ID = '" + vars["id"] + "'")
	results, err := bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		log.Printf("Error retrieving card: %v", err)
	}

	results.One(&result)
	if result == nil {
		w.Write([]byte("Could Not Find Cage Card"))
		return
	}
	final, _ := json.Marshal(result)

	w.Write(final)

}

func main() {
	cluster, err := gocb.Connect("couchbase://cbtest.ovpr.uga.edu/")
	if err != nil {
		log.Fatalf("Error Connecting: %v", err)
	}

	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "kbr03231",
		Password: "!@#$Uga1234",
	})

	bucket, err = cluster.OpenBucket("philbucket", "")
	if err != nil {
		log.Fatalf("Error Getting Bucket: %v", err)
	}
	defer bucket.Close()

	m := mux.NewRouter()
	m.HandleFunc("/cageCards/{id:C[0-9]{8}}", cageCards).Methods("GET")
	m.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", m))
}
