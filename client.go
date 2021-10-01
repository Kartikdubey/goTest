package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Person struct {
	Name   string  `json:"name"`
	Dob    string  `json:"dob"`
	Salary float64 `json:"salary"`
	Age    int     `json:"age"`
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
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json", person)

	req, err := http.NewRequest("POST", "http://localhost:8000/send", bytes.NewBuffer(content))
	req.Header.Set("storage-type", "file")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	/* response, err := http.Get("http://localhost:" + appPort + "/get")
	if err != nil {
		fmt.Println("HTTP req failed with error", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}*/
}
