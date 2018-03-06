package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	yl "github.com/gnewton/yastapii/lib"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Taxonomy struct {
	Id       uint64
	Name     string
	RootNode uint64
}

var channelSize = 5000

type DBInfo struct {
	boltdb *bolt.DB
	tx     *bolt.Tx
	bucket *bolt.Bucket
}

const nodeBucket = "nodeBucket"

//82696||Arthropoda, 194402||Agaricales, 18063||Magnoliopsida, 245||Enterobacteriaceae, 612944||Nectriaceae, 956286||Xanthomonadales, 180539||Carnivora
var otherTaxonomiesStartingTaxons = []uint64{194402, 82696, 18063, 245, 612944, 956286, 180539}

//var otherTaxonomiesStartingTaxons = []uint64{194402}

func makeAllTaxonomy(db *gorm.DB, boltdb *bolt.DB) {
	rand.Seed(time.Now().UnixNano())

	c := make(chan *Node, channelSize)
	done := make(chan struct{}, 0)

	tx, err := boltdb.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	b, err := tx.CreateBucketIfNotExists([]byte(nodeBucket))
	if err != nil {
		log.Fatal(err)
	}
	dbInfo := DBInfo{boltdb, tx, b}
	addNodes(&dbInfo, c, done)

	for i := 0; i < len(otherTaxonomiesStartingTaxons); i++ {

		t := otherTaxonomiesStartingTaxons[i]
		taxon := yl.GetTaxonomicUnitByTSN(db, t)
		root_node := new(Node)
		root_node.Id = counter
		root_node.Name = taxon.Complete_name

		counter += 1
		fmt.Println("\n\n\n\n*******************************************************************\n\n")
		fmt.Println("Taxonomy root:", root_node.Id, root_node.Name)
		fmt.Println(taxon)

		//root_node.Children = traverseTaxonomicUnits(db, c, t, root_node.Id, true)
		root_node.Children = traverseTaxonomicUnits(db, c, t, root_node.Id, false)
		err := addNode2(dbInfo.bucket, root_node)
		if err != nil {
			log.Fatal(err)
		}
	}

	if true {

		kingdoms := yl.GetTaxonomicUnitByFieldValueInt(db, "rank_id", 10)

		root_node := new(Node)
		root_node.Id = counterStart
		root_node.Name = "__root__"
		counter += 1
		fmt.Println("\n\n\n\n*******************************************************************\n\n")
		fmt.Println("Full Taxonomy root:", root_node.Id, root_node.Name)

		for i, _ := range kingdoms {
			k := kingdoms[i]
			fmt.Println(k)

			node := new(Node)
			node.Id = counter
			node.Name = k.Complete_name
			counter += 1

			node.Children = traverseTaxonomicUnits(db, c, k.Tsn, node.Id, false)
			err := addNode2(dbInfo.bucket, node)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Println("Closing c")
	close(c)
	log.Println("Closed c")
	log.Println("Waiting on done")
	_ = <-done

	log.Println("Got done")
	if false {
		tx, err = boltdb.Begin(false)
		b = tx.Bucket([]byte(nodeBucket))
		if b == nil {
			log.Fatal("Bucket does not exist")
		}

		cur := b.Cursor()

		log.Println("Start cursor")
		for k, v := cur.First(); k != nil; k, v = cur.Next() {
			//fmt.Println("\n----------------------------------------------")
			//fmt.Printf("key=%s, value=%s\n", k, v)
			// convert key
			//x, n := binary.Uvarint(k)
			_, n := binary.Uvarint(k)
			if n <= 0 {
				log.Fatal(errors.New("Error converting key: UVarint returns value <=0:" + strconv.Itoa(n)))
			}
			if n != len(k) {
				//fmt.Println("Uvarint did not consume all of in", n)
			}
			//fmt.Println(x)
			// convert Node struct
			var node Node
			buff := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buff)
			err = dec.Decode(&node)
			if err != nil {
				log.Fatal("decode error 1:", err)
			}
			//fmt.Printf("\n+++ %+v", node)
		}
		log.Println("End cursor")
	}

	tx, err = boltdb.Begin(false)
	b = tx.Bucket([]byte(nodeBucket))
	if b == nil {
		log.Fatal("Bucket does not exist")
	}

	log.Println("Start traverse")
	traverseNodes(b, 134, 0)
	//traverseNodes(b, 357841, 0)
}

var txCount = 0

const txSize = 35000

func addNodes(db *DBInfo, c chan *Node, done chan struct{}) {
	go func() {
		for node := range c {
			txCount += 1
			if txCount%5000 == 0 {
				fmt.Println(txCount)
			}
			if txCount%txSize == 0 {
				var err error
				fmt.Println("\tCommit")
				if err = db.tx.Commit(); err != nil {
					log.Fatal(err)
				}
				fmt.Println("\tDONE Commit")

				db.tx, err = db.boltdb.Begin(true)
				if err != nil {
					log.Fatal(err)
				}
				db.bucket = db.tx.Bucket([]byte(nodeBucket))
				if db.bucket == nil {
					log.Fatal("Bucket does not exist")
				}

			}
			err := addNode2(db.bucket, node)
			if err != nil {
				log.Fatal(err)
			}

		}
		if err := db.tx.Commit(); err != nil {
			log.Fatal(err)
		}
		log.Println("Sending done")
		done <- struct{}{}

	}()

}

func traverseNodes(b *bolt.Bucket, nodeId uint64, depth int) {
	node, err := getNode(b, nodeId)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("\n"+spaces(depth)+"+++ %+v", node)
	for i := 0; i < len(node.Children); i++ {
		traverseNodes(b, node.Children[i], depth+2)
	}
}

func spaces(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s = s + " "
	}
	return s

}

func traverseTaxonomicUnits(db *gorm.DB, c chan *Node, taxonParent uint64, parent uint64, random bool) []uint64 {
	children := yl.GetTaxonomicUnitChildrenById(db, taxonParent, []string{"tsn", "complete_name", "parent_tsn"})
	var actualChildren []uint64 = make([]uint64, 0)
	for i, _ := range children {
		if random {
			if rand.Intn(100) < 5 {
				continue
			}
		}
		child := children[i]

		node := new(Node)
		node.Id = counter
		actualChildren = append(actualChildren, node.Id)

		node.ActualParent = parent
		node.TaxonParent = child.Parent_tsn
		node.Taxon = child.Tsn
		node.Name = child.Complete_name
		counter += 1
		node.Children = traverseTaxonomicUnits(db, c, child.Tsn, node.Id, random)
		c <- node
	}
	return actualChildren
}
