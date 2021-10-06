package clientserver

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
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

func (s *Server) SendData(ctx context.Context, person *Person) (*Message, error) {
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
		return &Message{Body: "ok-write-to record,csv"}, nil
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
		return &Message{Body: "ok-write-to record,xml"}, nil
	}
}

func (s *Server) GetData(ctx context.Context, file *File) (*Person, error) {
	if file.File == "XML" {
		xmlFile, err := os.Open("users.xml")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Successfully Opened users.xml")
		// defer the closing of our xmlFile so that we can parse it later on
		defer xmlFile.Close()
		byteValue, _ := ioutil.ReadAll(xmlFile)

		// we initialize our Users array
		var record Record
		// we unmarshal our byteArray which contains our
		// xmlFiles content into 'users' which we defined above
		xml.Unmarshal(byteValue, &record)
		return &Person{Name: record.Name, Dob: record.Dob, Age: int32(record.Age), Salary: record.Salary}, nil
	} else {
		f, err := os.Open("records.csv")

		if err != nil {

			log.Fatal(err)
		}

		r := csv.NewReader(f)
		var record Record
		for {

			rec, err := r.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			record.Name = fmt.Sprintf("%s", rec[0])
			record.Dob = fmt.Sprintf("%s", rec[1])
			record.Salary = fmt.Sprintf("%s", rec[2])
			record.Age, _ = strconv.Atoi(rec[3])

		}
		return &Person{Name: record.Name, Dob: record.Dob, Age: int32(record.Age), Salary: record.Salary}, nil

	}
}

//Working
func (s *Server) UpdateData(ctx context.Context, in *Person) (*Message, error) {
	return nil, nil
}
