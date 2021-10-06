package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	//s"github.com/Kartikdubey/goTest/clientserver"
	"github.com/Kartikdubey/goTest/clientserver"
	"github.com/allegro/bigcache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var db *sql.DB
var err, cerr error
var cache *bigcache.BigCache

type Person struct {
	Name   string `json:"name"`
	Dob    string `json:"dob"`
	Salary string `json:"salary"`
	Age    int32  `json:"age"`
}

//getPerson to get all the data
//from FILE

func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("GET HIT")
	//var person Person
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := clientserver.NewServiceOneClient(conn)
	fileType := r.Header.Get("fileType")

	file := &clientserver.File{File: fileType}
	response, err := c.GetData(context.Background(), file)

	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response)

	json.NewEncoder(w).Encode(response)
}

//TO VALIDATE THE RECORDS

func validate(per Person) bool {
	var ageCheck bool
	if per.Age > 0 && per.Age < 110 {
		ageCheck = true
	}

	dob := regexp.MustCompile("^([0]?[1-9]|[1|2][0-9]|[3][0|1])[./-]([0]?[1-9]|[1][0-2])[./-]([0-9]{4}|[0-9]{2})$")
	name := regexp.MustCompile(`^[a-zA-Z]+$`)
	//salary := regexp.MustCompile(`(?!0+(?:\\.0+)?$)[0-9]+(?:\\.[0-9]+)?`)
	/*return (dob.MatchString(per.Dob) && name.MatchString(per.Name) && ageCheck)
	&& salary.MatchString(per.Salary))
	}*/
	return (dob.MatchString(per.Dob) && name.MatchString(per.Name) && ageCheck)
}

//createPerson to add a new record in CSV OR XML FILE
func createPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATE RECORD  HIT")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	var person Person
	err = json.Unmarshal(body, &person)

	// json.Unmarshal Error-Validate Types

	if err != nil {
		log.Panic(err)
	}

	log.Println("Request Came to server createPerson", person)
	fileType := r.Header.Get("fileType")
	//vALIDATE
	val := validate(person)
	fmt.Println("heaader---val ", fileType, val)

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := clientserver.NewServiceOneClient(conn)

	per := &clientserver.Person{Name: person.Name, Age: person.Age, Dob: person.Dob, Salary: person.Salary, Filetype: fileType}
	response, err := c.SendData(context.Background(), per)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(w, response.Body)
}

/*updatePerson to add a record in FILE
func updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Update HIT")

}

*/

func main() {
	fmt.Println("Main Started")

	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/get", getPerson).Methods("GET")
	router.HandleFunc("/send", createPerson).Methods("POST")
	//router.HandleFunc("/delete/{name}", deletePerson).Methods("DELETE")
	//router.HandleFunc("/update/{name}", updatePerson).Methods("PUT")

	http.ListenAndServe(":8000", router)

}
