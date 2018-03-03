package main

import (
	"errors"
	"fmt"
	yl "github.com/gnewton/yastapii/lib"
	"github.com/google/jsonapi"
	"log"
	"net/url"
	"strconv"
)

type Extras struct {
	links *Links
	meta  *Meta
}

type Links struct {
	offset, limit uint64
}

type Meta struct {
	Count      uint64
	TotalCount uint64
}

func addressAsString(v interface{}) string {
	tmp := fmt.Sprintf("%p", v)
	return tmp
}

func convertItisTaxonomicUnit(tu *yl.TaxonomicUnit) *Taxon {
	if tu == nil {
		return nil
	}
	taxon := new(Taxon)
	taxon.ID = tu.Tsn
	taxon.Name = tu.Unit_name1
	taxon.NameUsage = tu.Name_usage
	taxon.Source = "ITIS"

	if tu.Rank_id == 220 {
		taxon.Name += " " + tu.Unit_name2
	} else {
		if tu.Rank_id == 230 {
			taxon.Name += " " + tu.Unit_name2 + " " + tu.Unit_name3
		}
		if tu.Rank_id == 240 {
			taxon.Name += " " + tu.Unit_name2 + " " + tu.Unit_name3 + " " + tu.Unit_name4
		}
	}

	log.Println("Taxon rank id", tu.Rank_id)

	rank, ok := taxonUnitsMap[tu.Rank_id]
	if ok {
		taxon.RankName = rank.Rank_name
	}
	log.Println(rank.Rank_name, "[", tu.Unit_name1, "]", "[", tu.Unit_name2, "]", "[", tu.Unit_name3, "]", "[", tu.Unit_name4, "]")

	return taxon
}

//func convertItisTaxonomicUnits(tu []yl.TaxonomicUnit, paging bool) []*Taxon {
func convertItisTaxonomicUnits(tu []yl.TaxonomicUnit, paging bool) Taxons {
	if tu == nil || len(tu) == 0 {
		return nil
	}
	l := len(tu)

	taxa := make(Taxons, l)
	for i := 0; i < l; i++ {
		taxon := convertItisTaxonomicUnit(&tu[i])
		taxa[i] = taxon

	}

	return taxa
}

const PAGE_OFFSET = "page[offset]"
const PAGE_LIMIT = "page[limit]"

func makeOffsetLimits(offset, limit uint64) *jsonapi.Links {
	var links jsonapi.Links
	links = make(map[string]interface{})

	if offset+limit < numTaxons {
		links["next"] = jsonapi.Link{
			Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLNext(offset, limit)),
		}
	}

	if offset-limit >= 0 {
		links["previous"] = jsonapi.Link{
			Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLPrevious(offset, limit)),
		}
	}

	links["first"] = jsonapi.Link{
		Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLFirst(limit)),
	}

	links["this"] = jsonapi.Link{
		Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURL(offset, limit)),
	}

	links["last"] = jsonapi.Link{
		Href: fmt.Sprintf("http://localhost:8080/taxon?"+"%s", makeOffsetLimitURLLast(numTaxons-limit, limit)),
	}

	return &links
}

func makeOffsetLimitURL(offset, limit uint64) string {
	return fmt.Sprintf(PAGE_OFFSET+"=%d"+"&"+PAGE_LIMIT+"=%d", offset, limit)
}

func makeOffsetLimitURLNext(offset, limit uint64) string {
	return makeOffsetLimitURL(offset+limit, limit)
}

func makeOffsetLimitURLThis(offset, limit uint64) string {

	return makeOffsetLimitURL(offset, limit)
}

func makeOffsetLimitURLFirst(limit uint64) string {

	return makeOffsetLimitURL(0, limit)
}

func makeOffsetLimitURLLast(offset, limit uint64) string {

	return makeOffsetLimitURL(offset, limit)
}

func makeOffsetLimitURLPrevious(offset, limit uint64) string {
	if offset-limit < 0 {
		offset = 0
	} else {
		offset = offset - limit
	}
	return makeOffsetLimitURL(offset, limit)
	//return fmt.Sprintf(PAGE_OFFSET+"=%d"+"&"+PAGE_LIMIT+"=%d", offset, limit)
}

func makeOffsetLimit(query url.Values) (uint64, uint64, error) {
	var offset uint64 = 0
	var limit uint64 = 10
	var err error

	soffset := query.Get(PAGE_OFFSET)
	slimit := query.Get(PAGE_LIMIT)

	if soffset == "" && slimit == "" {
		return offset, limit, nil
	}
	if ((soffset == "") && (slimit != "")) || ((soffset == "") && (slimit != "")) {
		return 0, 0, errors.New("Both " + PAGE_OFFSET + " and " + PAGE_LIMIT + " must be set")
	}

	offset, err = strconv.ParseUint(soffset, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	limit, err = strconv.ParseUint(slimit, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return offset, limit, nil
}
