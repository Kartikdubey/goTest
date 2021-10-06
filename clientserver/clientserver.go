package clientserver

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"golang.org/x/net/context"
)

type Server struct {
}

type Record struct {
	Name   string `xml:"name"`
	Dob    string `xml:"dob"`
	Salary string `xml:"salary"`
	Age    int    `xml:"age"`
}

func (s *Server) SendData(ctx context.Context, person *Person) (*Person, error) {
	fmt.Println("Received message body from client--", person)

	if person.Filetype == "CSV" {
		file, err := os.Create("records.csv")
		defer file.Close()
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		w := csv.NewWriter(file)
		defer w.Flush()

		row := []string{person.Name, person.Dob, person.Salary, strconv.Itoa(int(person.Age))}

		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		return &Person{Name: "ok-write-to record,csv"}, nil
	} else {
		v := &Record{

			Name:   person.Name,
			Dob:    person.Dob,
			Salary: person.Salary,
			Age:    int(person.Age),
		}
		filename := "newrecords.xml"
		file, _ := os.Create(filename)
		xmlWriter := io.Writer(file)

		enc := xml.NewEncoder(xmlWriter)
		enc.Indent("  ", "    ")
		if err := enc.Encode(v); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		return &Person{Name: "ok-write-to record,xml"}, nil
	}
}
