package main

import (
	"fmt"
	yl "github.com/gnewton/yastapii/lib"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

const ID = "id"

func addHandlers(router *mux.Router, db *gorm.DB) {
	m := new(Manager)
	m.db = db
	router.HandleFunc("/node/{"+ID+"}", m.GetNode).Methods("GET")
	router.HandleFunc("/node", m.GetNodes).Methods("GET")

	router.HandleFunc("/taxon/{"+ID+"}", m.GetTaxon).Methods("GET")
	router.HandleFunc("/taxon", m.GetTaxa).Methods("GET")

}

type Manager struct {
	db *gorm.DB
}

func (m *Manager) GetTaxa(w http.ResponseWriter, r *http.Request) {
	log.Println("GetTaxa")
	w.Header().Set("Content-Type", jsonapi.MediaType)

	tu := yl.GetTaxonomicUnitAllOffsetLimit(m.db, 0, 100)
	log.Println("GetTaxa")
	taxons := convertItisTaxonomicUnits(tu, true)
	log.Println("GetTaxa")
	taxa := new(Taxa)
	taxa.Taxa = taxons
	//taxa.ID = 42

	//w.WriteHeader(201)
	if err := jsonapi.MarshalPayloadWithoutIncluded(w, taxa); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (m *Manager) GetTaxon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	params := mux.Vars(r)
	s, ok := params[ID]
	if !ok {
		http.Error(w, "Missing parameter", 500)
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		http.Error(w, "FIXXX", 500)
	}
	tu := yl.GetTaxonomicUnitByTSN(m.db, id)

	if tu == nil {
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Record does not exist",
			Detail: "Requesting record:" + s,
			Status: "404",
			//Meta:   &map[string]interface{}{"field": "some_field", "error": "bad type", "expected": "string", "received": "float64"},
		}})
		return

	}
	taxon := convertItisTaxonomicUnit(tu)

	w.WriteHeader(201)

	if err := jsonapi.MarshalPayloadWithoutIncluded(w, taxon); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (taxa Taxa) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"next": jsonapi.Link{
			Href: fmt.Sprintf("taxon/%d", 43),
		},
		"previous": jsonapi.Link{
			Href: fmt.Sprintf("taxon/%d", 55),
		},
	}
}

func (m *Manager) GetNodes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var node Node
	idStr, ok := params[ID]
	if !ok {
		http.Error(w, "Bad parameter", 500)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 400
		log.Fatal(err)
		http.Error(w, "Bad parameter", 500)
	}
	log.Println(id, ok)

	//tu :=

	//node.ID =

	taxonomy := Taxonomy{
		ID:   455,
		Name: "ITIS",
	}
	taxon := Taxon{
		Name:     "foo",
		ID:       123,
		RankName: "Genus",
	}
	parentNode := Node{
		Id:       87,
		Taxonomy: &taxonomy,
	}

	node = Node{
		//Name:     "this is a name",
		Id:       445,
		Taxon:    &taxon,
		Taxonomy: &taxonomy,
		Parent:   &parentNode,
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(201)

	if err := jsonapi.MarshalPayloadWithoutIncluded(w, &node); err != nil {
		http.Error(w, err.Error(), 500)
	}

}

func (m *Manager) GetNode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var node Node
	idStr, ok := params[ID]
	log.Println("params", r)
	log.Println("params", params)
	log.Println("----", idStr, ok)
	if ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			// 400
			w.Header().Set("Content-Type", jsonapi.MediaType)
			http.Error(w, "Bad parameter", 500)
			return
		}
		log.Println(id, ok)
	}
	//node.ID =

	taxonomy := Taxonomy{
		ID:   455,
		Name: "ITIS",
	}
	taxon := Taxon{
		Name:     "foo",
		ID:       123,
		RankName: "Genus",
	}
	parentNode := Node{
		Id:       87,
		Taxonomy: &taxonomy,
	}

	node = Node{
		//Name:     "this is a name",
		Id:       445,
		Taxon:    &taxon,
		Taxonomy: &taxonomy,
		Parent:   &parentNode,
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(201)

	if err := jsonapi.MarshalPayloadWithoutIncluded(w, &node); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
