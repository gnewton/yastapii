package main

import (
	yl "github.com/gnewton/yastapii/lib"
)

func convertItisTaxonomicUnit(tu *yl.TaxonomicUnit) *Taxon {
	if tu == nil {
		return nil
	}
	taxon := new(Taxon)
	taxon.ID = tu.Tsn
	taxon.Name = tu.Unit_name1

	return taxon
}

func convertItisTaxonomicUnits(tu []yl.TaxonomicUnit, paging bool) []*Taxon {
	if tu == nil || len(tu) == 0 {
		return nil
	}
	l := len(tu)

	taxa := make([]*Taxon, l)

	for i := 0; i < l; i++ {
		var taxon Taxon
		taxon.ID = tu[i].Tsn
		taxon.Name = tu[i].Unit_name1
		taxon.paging = paging
		taxa[i] = &taxon

	}

	return taxa
}
