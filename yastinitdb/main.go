package main

import (
	"github.com/boltdb/bolt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db := initDB("../data/tmp/ITIS.sqlite")
	defer func() {
		log.Println("Closing sqlite")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Closed sqlite")
	}()

	boltdb, err := bolt.Open("node.boltdb", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing bolt")
		err := boltdb.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Closed bolt")
	}()

	err = addIndexes(db)
	if err != nil {
		log.Fatal(err)
	}
	makeAllTaxonomy(db, boltdb)
}
