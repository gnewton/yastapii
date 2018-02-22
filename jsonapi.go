package main

import (
//"github.com/google/jsonapi"
//"time"
)

type Taxonomy struct {
	ID   int64 `json:primary,taxonomies`
	Name string
}

type Taxon struct {
	ID       int64 `json:primary,taxons`
	Name     string
	RankName string
}

type Node struct {
	Id       int64     `json:primary,nodes`
	Taxon    *Taxon    `json:relation,taxon,omitempty`
	Taxonomy *Taxonomy `json:relation,taxonomy,omitempty`
	Parent   *Node     `json:relation,parent,omitempty`
	Children []*Node   `json:relation,children,omitempty`
}
