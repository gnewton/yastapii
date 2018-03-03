package main

type Taxonomy struct {
	ID   uint64 `jsonapi:"primary,taxonomy"`
	Name string `jsonapi:"attr,name"`
}

// type Taxa struct {
// 	ID            int64    `jsonapi:"primary,taxons"`
// 	Taxa          []*Taxon `jsonapi:"relation,taxons"`
// 	offset, limit int64
// }

type Taxon struct {
	ID        uint64 `jsonapi:"primary,taxon"`
	Name      string `jsonapi:"attr,name,omitempty"`
	RankName  string `jsonapi:"attr,rank_name,omitempty"`
	NameUsage string `jsonapi:"attr,name_usage,omitempty"`
	Source    string `jsonapi:"attr,source,omitempty"`
}

type Node struct {
	Id uint64 `jsonapi:"primary,node"`
	//Taxon    *Taxon    `jsonapi:"relation,taxon"`
	Taxon    *Taxon    `jsonapi:"relation,taxon,omitempty"`
	Taxonomy *Taxonomy `jsonapi:"relation,taxonomy,omitempty"`
	Parent   *Node     `jsonapi:"relation,parent,omitempty"`
	Children []Node    `jsonapi:"relation,children,omitempty"`
}
