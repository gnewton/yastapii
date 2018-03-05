package main

import (
	"github.com/boltdb/bolt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"runtime"
)

func main() {
	log.Println("#routines: ", runtime.NumGoroutine())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("#routines: ", runtime.NumGoroutine())
	db := initDB("../data/tmp/ITIS.sqlite")
	defer func() {
		log.Println("Closing sqlite")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Closed sqlite")
	}()
	log.Println("#routines: ", runtime.NumGoroutine())
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
	log.Println("#routines: ", runtime.NumGoroutine())
	err = addIndexes(db)
	if err != nil {
		log.Fatal(err)
	}

	// err = makeNodeTable(nodeDb)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("#routines: ", runtime.NumGoroutine())
	log.Println("Start taxonomy")
	makeAllTaxonomy(db, boltdb)
	log.Println("End taxonomy")
	log.Println("#routines: ", runtime.NumGoroutine())
	//log.Fatal("foo")
}
