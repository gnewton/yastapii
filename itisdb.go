package main

import (
	"fmt"
	yl "github.com/gnewton/yastapii/lib"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func initDB(filename string) *gorm.DB {
	db, err := gorm.Open("sqlite3", filename)

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (node Node) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": jsonapi.Link{
			Href: fmt.Sprintf("node/%d", node.Id),
		},
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	router := mux.NewRouter()

	db := initDB("./data/ITIS.sqlite")
	defer db.Close()
	db.LogMode(true)
	err := cacheTaxonUnits(db)
	if err != nil {
		log.Fatal(err)
	}
	findMaxCounts(db)

	fmt.Println("===============================")
	addHandlers(router, db)
	log.Fatal(http.ListenAndServe(":8080", router))

	// never run
	fmt.Println("===============================")
	var tu []yl.TaxonomicUnit
	n := db.Offset(22111).Limit(100).Find(&tu)
	errors := n.GetErrors()
	if errors != nil && len(errors) > 0 {
		fmt.Println("mmmmmmmmm", errors)
		log.Fatal("foo")
	}

	for i, _ := range tu {
		fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%% ")
		fmt.Println(" ")
		rank, ok := taxonUnitsMap[uint64(tu[i].Rank_id)]
		if !ok {
			log.Fatal("LLm2")
		}
		fmt.Println(tu[i].Tsn, " ", tu[i].Unit_name1, "   ", rank.Rank_name)
		//bb := getTaxonomicUnitByTSN(db, tu[i].Tsn)
		//fmt.Printf("** %v\n", bb)
		parents := yl.GetTaxonomicUnitAncestors(db, &tu[i])
		fmt.Println("Parents")
		for j, _ := range parents {
			rank, ok := taxonUnitsMap[uint64(parents[j].Rank_id)]
			if !ok {
				log.Fatal("LL")
			}
			fmt.Println("\t", parents[j].Tsn, " ", parents[j].Unit_name1, "   ", rank.Rank_name)
		}
		children := yl.GetTaxonomicUnitChildren(db, &tu[i], nil)
		fmt.Println("Children")
		for j, _ := range children {
			rank, ok := taxonUnitsMap[uint64(children[j].Rank_id)]
			if !ok {
				log.Fatal("MM")
			}
			fmt.Println("\t +++", children[j].Tsn, " ", children[j].Unit_name1, " ", children[j].Unit_name2, " ", children[j].Unit_name3, " ", children[j].Unit_name4, "   ", rank.Rank_name)
		}
	}
}
