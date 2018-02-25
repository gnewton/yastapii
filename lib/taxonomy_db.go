package lib

import (
	"github.com/jinzhu/gorm"
	"time"
)

type TaxonUnitType2 struct {
	KingdomId       int64
	RankId          int64
	RankName        string
	DirParentRankId int64
	ReqParentRankId int64
}

type TaxonUnitType struct {
	Kingdom_id         int64
	Rank_id            int64
	Rank_name          string
	Dir_parent_rank_id int64
	Req_parent_rank_id int64
}

type TaxonomicUnit struct {
	Tsn                int64
	Unit_ind1          []byte
	Unit_name1         string
	Unit_ind2          string
	Unit_name2         string
	Unit_ind3          string
	Unit_name3         string
	Unit_ind4          string
	Unit_name4         string
	Unnamed_taxon_ind  string
	Name_usage         string
	Unaccept_reason    string
	Credibility_rtng   string
	Completeness_rtng  string
	Currency_rating    string
	Phylo_sort_seq     int32
	Initial_time_stamp time.Time
	Parent_tsn         int64
	Taxon_author_id    int64
	Hybrid_author_id   int64
	Kingdom_id         int64
	Rank_id            int64
	Update_date        time.Time
	Uncertain_prnt_ind string
	N_usage            string
	Complete_name      string
}

func Count(db *gorm.DB, o struct{}, where string) int64 {
	var count int64
	db.Find(&o).Count(&count)
	return count
}

func GetTaxonomicUnitByTSN(db *gorm.DB, tsn int64) *TaxonomicUnit {
	var tu TaxonomicUnit
	db.Where("tsn=?", tsn).First(&tu)
	if tu.Tsn == 0 {
		return nil
	}
	return &tu
}

func GetTaxonomicUnitAllOffsetLimit(db *gorm.DB, offset, limit int64) []TaxonomicUnit {
	var tu []TaxonomicUnit
	db.Offset(offset).Limit(limit).Find(&tu)
	return tu
}

func GetTaxonomicUnitChildren(db *gorm.DB, tu *TaxonomicUnit) []*TaxonomicUnit {
	var children []*TaxonomicUnit
	db.Where("parent_tsn=?", tu.Tsn).Find(&children)
	return children
}

func GetTaxonomicUnitAncestors(db *gorm.DB, tu *TaxonomicUnit) []*TaxonomicUnit {
	if tu.Parent_tsn <= 0 {
		return nil
	}
	parent := GetTaxonomicUnitByTSN(db, tu.Parent_tsn)
	//fmt.Println(parent)

	parents := make([]*TaxonomicUnit, 1)

	parents = append(GetTaxonomicUnitAncestors(db, parent), parent)
	return parents

}
