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
/*from FILE
func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("GET HIT")
	var persons []Person
	result, err := db.Query("SELECT * FROM Persons")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var person Person
		err := result.Scan(&person.Age, &person.Name)
		if err != nil {
			panic(err.Error())
		}
		persons = append(persons, person)
	}
	fmt.Println("Response from db", persons)
	json.NewEncoder(w).Encode(persons)
}*/
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
	elliot := Person{
		Name: "Elliot",
		Age:  24}
	_, err := proto.Marshal(&elliot)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}*/
	return (dob.MatchString(per.Dob) && name.MatchString(per.Name) && ageCheck)
}

//createPerson to add a new record in CSV OR XML FILE
func createPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATE HIT")
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
	val := validate(person)
	fmt.Println("heaader---val ", fileType, val)
	fmt.Fprintf(w, "New person was created")

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
	log.Printf("Response from server: %s", response.Name)
}

/*getSpecificPersons to get a particular row from table
func getSpecificPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Get Specific HIT")
	params := mux.Vars(r)
	result, err := db.Query("SELECT pAge,pName FROM Persons WHERE pAge >= ?", params["age"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var pers []Person
	for result.Next() {
		var per Person
		err := result.Scan(&per.Age, &per.Name)
		if err != nil {
			panic(err.Error())
		}
		pers = append(pers, per)
	}
	json.NewEncoder(w).Encode(pers)
}

/*updatePerson to add a record in FILE
func updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Update HIT")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE Persons SET pAge = ? WHERE pName = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	var per Person
	json.Unmarshal(body, &per)
	age := per.Age
	_, err = stmt.Exec(age, params["name"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Person with Name = %s was updated", params["name"])
}

/*deletePerson to delete record
func deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("DELETE HIT")
	params := mux.Vars(r)
	result, err := db.Query("SELECT * FROM Persons WHERE pName = ?", params["name"])
	if err != nil {
		panic(err.Error())
	}
	var per Person
	for result.Next() {
		errs := result.Scan(&per.Age, &per.Name)
		if errs != nil {
			panic(errs.Error())
		}
	}
	cache.Set(strconv.Itoa(per.Age), []byte(per.Name))

	stmt, err := db.Prepare("DELETE FROM Persons WHERE pName = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["name"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Person with name = %s was deleted", params["name"])
}*/

func main() {
	fmt.Println("Main Started")

	defer db.Close()
	router := mux.NewRouter()
	//router.HandleFunc("/get", getPerson).Methods("GET")
	router.HandleFunc("/send", createPerson).Methods("POST")
	//router.HandleFunc("/delete/{name}", deletePerson).Methods("DELETE")
	//router.HandleFunc("/update/{name}", updatePerson).Methods("PUT")
	//router.HandleFunc("/get/{age}", getSpecificPersons).Methods("GET")

	http.ListenAndServe(":8000", router)

}
