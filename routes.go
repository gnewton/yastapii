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

var mm map[string]*Extras

func addHandlers(router *mux.Router, db *gorm.DB) {
	m := new(Manager)
	m.db = db
	router.HandleFunc("/node/{"+ID+":[0-9]+}", m.GetNode).Methods("GET")
	router.HandleFunc("/node", m.GetNodes).Methods("GET")

	router.HandleFunc("/taxon/{"+ID+":[0-9]+}", m.GetSingleTaxon).Methods("GET")
	router.HandleFunc("/taxon?filter[", m.GetTaxonByQuery).Methods("GET")
	//router.HandleFunc("/taxon/{"+ID+":[0-9]+}", m.GetTaxon).Methods("GET")
	//router.HandleFunc("/taxon", m.GetTaxa).Queries(PAGE_OFFSET, "{"+PAGE_OFFSET+"}").Methods("GET")
	//router.HandleFunc("/taxon", m.GetTaxa).Queries("foo", "{foo}").Methods("GET")

	router.HandleFunc("/taxon", m.GetAllTaxons).Methods("GET")
	mm = make(map[string]*Extras)
}

type Manager struct {
	db *gorm.DB
}

func (m *Manager) GetTaxonByQuery(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	log.Println("params", params)
}

func (m *Manager) GetAllTaxons(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAllTaxa")
	w.Header().Set("Content-Type", jsonapi.MediaType)

	params := mux.Vars(r)
	log.Printf("params: %v\n", params)

	var err error
	extras := new(Extras)
	extras.meta = new(Meta)
	links := new(Links)
	extras.links = links
	links.offset, links.limit, err = makeOffsetLimit(r.URL.Query())

	if err != nil {
		log.Println(err)
		return
	}
	log.Println("============", links.offset, links.limit)

	if err != nil {
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Paging number conversion error",
			Detail: err.Error(),
			Status: "404",
			//Meta:   &map[string]interface{}{"field": "some_field", "error": "bad type", "expected": "string", "received": "float64"},
		}})
		return
	}

	tu := yl.GetTaxonomicUnitAllOffsetLimit(m.db, extras.links.offset, extras.links.limit)
	taxons := convertItisTaxonomicUnits(tu, true)
	extras.meta.Count = int64(len(taxons))
	extras.meta.TotalCount = numTaxons

	addr := addressAsString(taxons)
	log.Println("GetTaxa addr", addr)
	mm[addr] = extras
	if err := jsonapi.MarshalPayloadWithoutIncluded(w, taxons); err != nil {
		http.Error(w, err.Error(), 500)
	}
	delete(mm, addr)
}

func (taxon Taxon) JSONAPILinks() *jsonapi.Links {

	//func (taxon []Taxon) JSONAPILinks() *jsonapi.Links {

	return &jsonapi.Links{
		"this": jsonapi.Link{
			Href: fmt.Sprintf("http://localhost:8080/taxon/%d", taxon.ID),
		},
	}
}

type Taxons []*Taxon

func (taxa Taxons) JSONAPIMeta() *jsonapi.Meta {
	log.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
	addr := addressAsString(taxa)
	extras, ok := mm[addr]
	if !ok {
		log.Fatal("Missing address")
	}
	a := fmt.Sprintf("%p", &taxa)
	log.Println(a)
	return &jsonapi.Meta{
		"Version":     1,
		"Software":    "yastapii",
		"Author":      "Glen Newton",
		"count":       extras.meta.Count,
		"total_count": extras.meta.TotalCount,
	}
}

func (taxa Taxons) JSONAPILinks() *jsonapi.Links {
	addr := addressAsString(taxa)
	extras, ok := mm[addr]
	if !ok {
		log.Fatal("Missing address")
	}
	// FIXX with RWMutex

	return makeOffsetLimits(extras.links.offset, extras.links.limit)
}

func (m *Manager) GetSingleTaxon(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSingleTaxon")
	w.Header().Set("Content-Type", jsonapi.MediaType)
	params := mux.Vars(r)
	s, ok := params[ID]
	if !ok {
		http.Error(w, "Missing parameter", 500)
	}
	log.Printf("params: %v\n", params)
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

	//if err := jsonapi.MarshalPayloadWithoutIncluded(w, taxon); err != nil {
	if err := jsonapi.MarshalPayload(w, taxon); err != nil {
		http.Error(w, err.Error(), 500)
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
	if ok {
		_, err := strconv.Atoi(idStr)
		if err != nil {
			// 400
			w.Header().Set("Content-Type", jsonapi.MediaType)
			http.Error(w, "Bad parameter", 500)
			return
		}
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
