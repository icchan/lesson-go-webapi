package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

type Root struct {
	Status string	`json:"status"` 
	Person []Guy	`xml:"People>Person" json:"people"` 
}

type Payload struct {
	Person Guy		`json:"person"` 
}

type Guy struct {
	Name string 	`json:"name"` 
	Age int 		`json:"age"`
	Pets []string 	`xml:"Pets>Pet" json:"pets"`
}

func main() {
	host := "localhost:1337"
	fmt.Printf("Listening on %s ... \n", host)
	http.HandleFunc("/", index)
	http.HandleFunc("/api", serveRest)
	http.ListenAndServe(host, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the index")
}

func serveRest(w http.ResponseWriter, r *http.Request) {
	// get some dummy data
	payload := getDummyList()

	// default output is json
	formatter := json.MarshalIndent
	contentType := "application/json; charset=utf-8"

	// check if the accept header is xml
	if accept := r.Header.Get("Accept"); accept == "application/xml" {
		contentType = "application/xml; charset=utf-8"
		formatter = xml.MarshalIndent
	}

	// set the content-type
	w.Header().Set("Content-Type", contentType)

	// format the response
	response, err := formatter(payload, " ", "   ")

	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(response))
}

func getDummyList() (Root) {
	var list []Guy
	list = append(list, Guy{"Ian", 29, []string{"chopi", "yopi"}})
	list = append(list, Guy{"Finn", 12, []string{"Jake"}})
	list = append(list, Guy{"Ash", 10, []string{"pikachu", "bulbasaur", "butterfree"}})

	p := Root{"OK", list}
	return p
}