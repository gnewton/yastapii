package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var taxonUnitsMap map[int64]*TaxonUnitType

func cacheTaxonUnits(db *gorm.DB) error {
	taxonUnitsMap = make(map[int64]*TaxonUnitType, 200)
	var tut []TaxonUnitType
	db.Find(&tut)
	for i, _ := range tut {
		fmt.Printf("%v\n", tut[i])
		taxonUnitsMap[tut[i].Rank_id] = &tut[i]
	}
	return nil
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
