package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Person struct {
	Name   string `json:"name"`
	Dob    string `json:"dob"`
	Salary string `json:"salary"`
	Age    int    `json:"age"`
}

//client connection check
func main() {
	fmt.Println("Client Started")
	//appPort := "8000"
	content, err := ioutil.ReadFile("input.json")
	var person Person
	err = json.Unmarshal(content, &person)
	// json.Unmarshal Error
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Successfully Opened users.json", person)

	req, err := http.NewRequest("POST", "http://localhost:8000/send", bytes.NewBuffer(content))
	req.Header.Set("fileType", "XML")
	req.Header.Set("Content-Type", "application/json")

	/*client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	*/
	response, err := http.Get("http://localhost:8000/get")
	if err != nil {
		fmt.Println("HTTP req failed with error", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
