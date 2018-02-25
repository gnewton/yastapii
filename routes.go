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
	router.HandleFunc("/node/{"+ID+":[0-9]+}", m.GetNode).Methods("GET")
	router.HandleFunc("/node", m.GetNodes).Methods("GET")

	router.HandleFunc("/taxon/{"+ID+":[0-9]+}", m.GetTaxon).Methods("GET")
	//router.HandleFunc("/taxon/{"+ID+":[0-9]+}", m.GetTaxon).Methods("GET")
	//router.HandleFunc("/taxon", m.GetTaxa).Queries(PAGE_OFFSET, "{"+PAGE_OFFSET+"}").Methods("GET")
	//router.HandleFunc("/taxon", m.GetTaxa).Queries("foo", "{foo}").Methods("GET")

	router.HandleFunc("/taxon", m.GetTaxa).Methods("GET")

}

type Manager struct {
	db *gorm.DB
}

func (m *Manager) GetTaxa(w http.ResponseWriter, r *http.Request) {

	log.Println("GetTaxa")
	log.Println(r.URL.Query())
	w.Header().Set("Content-Type", jsonapi.MediaType)

	offset, limit, err := makeOffsetLimit(r.URL.Query())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("============", offset, limit)
	//offset, limit, err := makeOffsetLimit(mux.Vars(r))

	if err != nil {
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Paging number conversion error",
			Detail: err.Error(),
			Status: "404",
			//Meta:   &map[string]interface{}{"field": "some_field", "error": "bad type", "expected": "string", "received": "float64"},
		}})
		return
	}

	tu := yl.GetTaxonomicUnitAllOffsetLimit(m.db, offset, limit)
	log.Println("GetTaxa")
	taxons := convertItisTaxonomicUnits(tu, true)
	log.Println("GetTaxa")
	//taxa := new(Taxa)
	//taxa := new(Taxons)
	//taxa.Taxa = taxons
	//taxa.offset = offset
	//taxa.limit = limit
	//taxa.ID = 42

	//w.WriteHeader(201)
	//if err := jsonapi.MarshalPayloadWithoutIncluded(w, taxa); err != nil {
	//if err := jsonapi.MarshalPayloadWithoutIncluded(w, taxons); err != nil {
	if err := jsonapi.MarshalPayload(w, taxons); err != nil {
		//if err := jsonapi.MarshalPayload(w, taxa); err != nil {
		http.Error(w, err.Error(), 500)
	}
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
	a := fmt.Sprintf("mm%pmm", &(taxa[0]))
	log.Println(a)
	return &jsonapi.Meta{
		"details": "sample details here",
	}
}

func (taxa Taxons) JSONAPILinks() *jsonapi.Links {
	log.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%mm")
	return &jsonapi.Links{
		"next": jsonapi.Link{
			//Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLNext(taxa.offset, taxa.limit)),
			Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLNext(10, 10)),
		},
		"previous": jsonapi.Link{
			//Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLPrevious(taxa.offset, taxa.limit)),
			Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLNext(10, 10)),
		},
		"first": jsonapi.Link{
			//Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLFirst(taxa.limit)),
			Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLNext(10, 10)),
		},
		"this": jsonapi.Link{
			//Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURL(taxa.offset, taxa.limit)),
			Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLNext(10, 10)),
		},

		"last": jsonapi.Link{
			//Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLLast(10000, taxa.limit)),
			Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLNext(10, 10)),
		},
	}
}

func (m *Manager) GetTaxon(w http.ResponseWriter, r *http.Request) {
	log.Println("GetTaxon")
	w.Header().Set("Content-Type", jsonapi.MediaType)
	params := mux.Vars(r)
	log.Println(params)
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
