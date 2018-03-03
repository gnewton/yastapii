package main

import (
	"fmt"
	yl "github.com/gnewton/yastapii/lib"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var taxonUnitsMap map[int64]*yl.TaxonUnitType
var numTaxons int64

func findMaxCounts(db *gorm.DB) {
	numTaxons = yl.CountTaxonomicUnits(db, "")
}

func cacheTaxonUnits(db *gorm.DB) error {
	taxonUnitsMap = make(map[int64]*yl.TaxonUnitType, 200)
	var tut []yl.TaxonUnitType
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
