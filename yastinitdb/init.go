package main

import (
	"github.com/jinzhu/gorm"
	"log"
)

func dbErrors(db *gorm.DB) error {
	if err := db.Error; err != nil {
		errors := db.GetErrors()
		printDbErrors(errors)
		return err
	}
	return nil
}

func execSql(db *gorm.DB, sql []string) error {
	for i, _ := range sql {
		//tmpdb := db.Exec(sql[i])
		_ = db.Exec(sql[i])
		/*
			if err := tmpdb.Error; err != nil {
				errors := tmpdb.GetErrors()
				printDbErrors(errors)
				return err
			}
		*/

	}
	return nil
}

//
func addIndexes(db *gorm.DB) error {
	newIndexes := []string{INDEX_TU_PARENT, INDEX_TU_NAME1, INDEX_TU_NAME2, INDEX_TU_NAME3, INDEX_TU_NAME4, INDEX_TU_ALL}
	return execSql(db, newIndexes)

	return nil
}
func printDbErrors(dberrs []error) {
	for _, err := range dberrs {
		log.Println(err)
	}
}

// func makeNodeTable(db *gorm.DB) error {
// 	nodeSql := []string{CREATE_NODE_TABLE, INDEX_NODE_PARENT, "PRAGMA journal_mode=OFF", "PRAGMA synchronous=OFF", "PRAGMA locking_mode=EXCLUSIVE", "PRAGMA temp_store=memory"}
// 	return execSql(db, nodeSql)
// }

func initDB(filename string) *gorm.DB {
	db, err := gorm.Open("sqlite3", filename+"?cache=shared&mode=rwc")

	if err != nil {
		log.Fatal(err)
	}
	//db.DB().SetMaxOpenConns(1)
	return db
}
