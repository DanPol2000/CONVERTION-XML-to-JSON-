package main

import (
	_ "bytes"
	"encoding/csv"
	_ "flag"
	"fmt"
	_ "io"
	"log"
	"os"
	_"strings"
	"database/sql"
	_ "github.com/lib/pq"

	// _ elastic "github.com/elastic/go-elasticsearch/v7"
	// _ "github.com/gocarina/gocsv"
	// _ "github.com/jszwec/csvutil"
)

func init() {
	log.SetFlags(0)
	// log.SetPrefix("ERROR: ")
}

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = ""
    dbname   = "places"
)

type Place struct {
	ID string `db:"Text"`
	Name    string `csv:"Name" db:"Text"`
	Address string `csv:"Address" db:"Text"`
	Phone   string `csv:"Phone" db:"Text"`
	Loc     struct {
		lon string `csv:"Longitude" db:"Text"`
		lat string `csv:"Latitude" db:"Text"`
	}
}

const schemaQr string = 
`DROP TABLE IF EXISTS places;
 CREATE TABLE places (
    index Text,
    name TEXT,
    address TEXT ,
    phone TEXT,
    longitude Text,
    latitude Text
 );`

// func createElastic() {
// 	es, _ := elastic.NewDefaultClient()
// 	// log.Println(es.Info())
// 	res, err := es.Index(
// 		"places",                               // Index name
// 		strings.NewReader(`{"name" : "text"}`), // Document body
// 		es.Index.WithDocumentID("1"),           // Document ID
// 		es.Index.WithRefresh("true"),           // Refresh
// 	)
// 	defer res.Body.Close()
// 	// res, err := es.Delete("places", "1")
// 	if err != nil {
// 		log.Fatalln("ERROR: ", err)
// 	}
// 	es.Search()
// 	log.Println(res)
// }

func main() {
	db := createDb()
	csv_file, _ := os.Open("new.csv")
	reader := csv.NewReader(csv_file)
	reader.Comma = '\t'
	r, _ := reader.ReadAll()
	for i, line := range r {
		if i == 0 {
			continue
		}
		var pl Place
		for j, field := range line {
			if j == 0 {
				pl.Name = field
			}
			if j == 1 {
				pl.Address = field
			}
			if j == 2 {
				pl.Phone = field
			}
			if j == 3 {
				pl.Loc.lat = field
			}
			if j == 4 {
				pl.Loc.lon = field
			}
			_, err := db.Exec("INSERT INTO places (ID, Name, Address, Phone, Longitude, Latitude) VALUES ($1,$2, $3, $4, $5, $6)", pl.ID, pl.Name, pl.Address, pl.Phone, pl.Loc.lat, pl.Loc.lon)
			if err != nil {
				log.Fatalln("ERROR: ", err)
			}
		}
	}
	fmt.Println("Process finished")
	defer db.Close()
}


func createDb() *sql.DB {

	postgr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", postgr)
	if err != nil {
		log.Fatalln("ERROR1: ", err)
	}
	_, err = db.Exec(schemaQr)
	if err != nil {
		log.Fatalln("ERROR3: ", err)
	}
	return db 
}