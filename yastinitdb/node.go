package main

import (
	//"errors"
	"github.com/boltdb/bolt"
	//yl "github.com/gnewton/yastapii/lib"
	//"database/sql"
	//"fmt"
	//"github.com/jinzhu/gorm"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
)

type Node struct {
	Id           uint64
	Children     []uint64
	ActualParent uint64
	//
	Name        string
	TaxonParent uint64
	Taxon       uint64
}

const counterStart = 100

var counter uint64 = counterStart

var n int64 = 0

//func addNode2(db *gorm.DB, node *Node) error {
func addNode2(bucket *bolt.Bucket, node *Node) error {
	keyBytes := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(keyBytes, node.Id)

	var structBytes bytes.Buffer
	enc := gob.NewEncoder(&structBytes)
	err := enc.Encode(*node)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	err = bucket.Put(keyBytes, structBytes.Bytes())
	return err
}
