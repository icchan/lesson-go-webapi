package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"encoding/xml"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type Root struct {
	Trainers []Trainer
}
type Trainer struct {
	Id int
	Name string
	Age int
	Hometown string
}

func main() {
	host := "localhost:1337"
	fmt.Printf("Listening on %s ... \n", host)
	http.HandleFunc("/", index) 
	http.HandleFunc("/trainer", serveRest)
	http.ListenAndServe(host, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I wanna be the very best, like no one ever was")
	v := r.URL.Query()

	fmt.Printf("%v\n", v.Get("params"))
}

func serveRest(w http.ResponseWriter, r *http.Request) {
	// get list from db
	payload := getTrainers()

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

/*

CREATE TABLE IF NOT EXISTS `trainer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `age` int(11) NOT NULL,
  `hometown` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=3 ;

INSERT INTO `trainer` (`id`, `name`, `age`, `hometown`) VALUES
(1, 'Ian', 29, 'Sydney'),
(2, 'Ash', 10, 'Pallet Town');

*/

func getTrainers() Root {
	con, err := sql.Open("mysql", "root:@/lesson")
	defer con.Close()

	rows, err := con.Query("select * from trainer")
	if err != nil { panic(err) }

	trainers := make([]Trainer,0,2)
	var (
		id, age int
		name, hometown string
	)

	for rows.Next() {
	    err = rows.Scan(&id, &name, &age, &hometown)
	    if err != nil { /* error handling */}
	    trainers = append(trainers, Trainer{id, name, age, hometown})
	    //fmt.Println(name)
	}

	root := Root{trainers}

	return root
}
