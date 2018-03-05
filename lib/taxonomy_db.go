package lib

import (
	"github.com/jinzhu/gorm"
	"time"
)

type TaxonUnitType2 struct {
	KingdomId       uint64
	RankId          uint64
	RankName        string
	DirParentRankId uint64
	ReqParentRankId uint64
}

type TaxonUnitType struct {
	Kingdom_id         uint64
	Rank_id            uint64
	Rank_name          string
	Dir_parent_rank_id uint64
	Req_parent_rank_id uint64
}

type TaxonomicUnit struct {
	Tsn                uint64
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
	Parent_tsn         uint64
	Taxon_author_id    uint64
	Hybrid_author_id   uint64
	Kingdom_id         uint64
	Rank_id            uint64
	Update_date        time.Time
	Uncertain_prnt_ind string
	N_usage            string
	Complete_name      string
}

func CountTaxonomicUnits(db *gorm.DB, where string) uint64 {
	var count uint64
	var tmp TaxonomicUnit
	db.Select("tsn").Find(&tmp).Count(&count)
	return count
}

func GetTaxonomicUnitByTSN(db *gorm.DB, tsn uint64) *TaxonomicUnit {
	var tu TaxonomicUnit
	db.Where("tsn=?", tsn).First(&tu)
	if tu.Tsn == 0 {
		return nil
	}
	return &tu
}

func GetTaxonomicUnitAllOffsetLimit(db *gorm.DB, offset, limit uint64) []TaxonomicUnit {
	var tu []TaxonomicUnit
	db.Offset(offset).Limit(limit).Find(&tu)
	return tu
}

func GetTaxonomicUnitChildren(db *gorm.DB, tu *TaxonomicUnit, fields []string) []*TaxonomicUnit {
	return GetTaxonomicUnitChildrenById(db, tu.Tsn, fields)
}

func GetTaxonomicUnitChildrenById(db *gorm.DB, tsn uint64, fields []string) []*TaxonomicUnit {
	var children []*TaxonomicUnit
	if fields == nil || len(fields) == 0 {
		db.Where("parent_tsn=?", tsn).Find(&children)
	} else {
		db.Select(fields).Where("parent_tsn=?", tsn).Find(&children)
	}

	return children
}

func GetTaxonomicUnitByFieldValueInt(db *gorm.DB, fieldName string, fieldValue uint64) []*TaxonomicUnit {
	var children []*TaxonomicUnit
	db.Where(fieldName+"=?", fieldValue).Find(&children)
	return children
}

func GetTaxonomicUnitByFieldValueString(db *gorm.DB, fieldName, fieldValue string) []*TaxonomicUnit {
	var children []*TaxonomicUnit
	db.Where(fieldName+"=\"?\"", fieldValue).Find(&children)
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
