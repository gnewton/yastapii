package main

import (
	"errors"
	"fmt"
	yl "github.com/gnewton/yastapii/lib"
	"log"
	"net/url"
	"strconv"
)

func convertItisTaxonomicUnit(tu *yl.TaxonomicUnit) *Taxon {
	if tu == nil {
		return nil
	}
	taxon := new(Taxon)
	taxon.ID = tu.Tsn
	taxon.Name = tu.Unit_name1
	taxon.NameUsage = tu.Name_usage

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

	//taxa := make([]*Taxon, l)
	taxa := make(Taxons, l)

	for i := 0; i < l; i++ {
		taxon := convertItisTaxonomicUnit(&tu[i])
		taxa[i] = taxon

	}

	return taxa
}

const PAGE_OFFSET = "page[offset]"
const PAGE_LIMIT = "page[limit]"

func makeOffsetLimitURL(offset, limit int64) string {
	return fmt.Sprintf(PAGE_OFFSET+"=%d"+"&"+PAGE_LIMIT+"=%d", offset, limit)
}

func makeOffsetLimitURLNext(offset, limit int64) string {
	//return fmt.Sprintf(PAGE_OFFSET+"=%d"+"&"+PAGE_LIMIT+"=%d", offset+limit, limit)
	return makeOffsetLimitURL(offset+limit, limit)
}

func makeOffsetLimitURLThis(offset, limit int64) string {
	//return fmt.Sprintf(PAGE_OFFSET+"=%d"+"&"+PAGE_LIMIT+"=%d", offset, limit)
	return makeOffsetLimitURL(offset, limit)
}

func makeOffsetLimitURLFirst(limit int64) string {
	//return fmt.Sprintf(PAGE_OFFSET+"=%d"+"&"+PAGE_LIMIT+"=%d", 0, limit)
	return makeOffsetLimitURL(0, limit)
}

func makeOffsetLimitURLLast(offset, limit int64) string {
	//return fmt.Sprintf(PAGE_OFFSET+"=%d"+"&"+PAGE_LIMIT+"=%d", 0, limit)
	return makeOffsetLimitURL(offset, limit)
}

func makeOffsetLimitURLPrevious(offset, limit int64) string {
	if offset-limit < 0 {
		offset = 0
	} else {
		offset = offset - limit
	}
	return makeOffsetLimitURL(offset, limit)
	//return fmt.Sprintf(PAGE_OFFSET+"=%d"+"&"+PAGE_LIMIT+"=%d", offset, limit)
}

func makeOffsetLimit(query url.Values) (int64, int64, error) {
	log.Println(query)
	var offset int64 = 0
	var limit int64 = 10
	var err error

	soffset := query.Get(PAGE_OFFSET)
	slimit := query.Get(PAGE_LIMIT)
	log.Println(soffset, "|", slimit, "|")

	if soffset == "" && slimit == "" {
		return offset, limit, nil
	}
	if ((soffset == "") && (slimit != "")) || ((soffset == "") && (slimit != "")) {
		return -1, -1, errors.New("Both " + PAGE_OFFSET + " and " + PAGE_LIMIT + " must be set")
	}

	offset, err = strconv.ParseInt(soffset, 10, 64)
	if err != nil {
		log.Println(err)
		return -1, -1, nil
	}

	limit, err = strconv.ParseInt(slimit, 10, 64)
	if err != nil {
		log.Println(err)
		return -1, -1, nil
	}

	return offset, limit, nil
}
