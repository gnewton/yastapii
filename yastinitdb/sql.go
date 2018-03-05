package main

// taxonomic_units indexes
const INDEX_TU_PARENT = `CREATE INDEX IF NOT EXISTS "taxonomic_units_taxon_unit_parent" ON "taxonomic_units" ("parent_tsn");`
const INDEX_TU_NAME1 = `CREATE INDEX IF NOT EXISTS "taxonomic_units_unit_name1" ON "taxonomic_units" ("unit_name1");`
const INDEX_TU_NAME2 = `CREATE INDEX IF NOT EXISTS "taxonomic_units_unit_name2" ON "taxonomic_units" ("unit_name2");`
const INDEX_TU_NAME3 = `CREATE INDEX IF NOT EXISTS "taxonomic_units_unit_name3" ON "taxonomic_units" ("unit_name3");`
const INDEX_TU_NAME4 = `CREATE INDEX IF NOT EXISTS "taxonomic_units_unit_name4" ON "taxonomic_units" ("unit_name4");`
const INDEX_TU_ALL = `CREATE INDEX IF NOT EXISTS "taxonomic_units_unit_name_all" ON "taxonomic_units" ("unit_name1","unit_name2","unit_name3","unit_name4");`

//node table
// const CREATE_NODE_TABLE = `CREATE TABLE IF NOT EXISTS "nodes"("id" int8 NOT NULL, "parent" int8, "name" varchar(64), "taxon" int(11), PRIMARY KEY ("id"));`
// const INDEX_NODE_PARENT = `CREATE INDEX IF NOT EXISTS "node_parents" ON "nodes" ("parent");`

//taxonomy table
const CREATE_TAXONOMY_TABLE = `CREATE TABLE IF NOT EXISTS "taxonomys"("id" int8 NOT NULL, "name" varchar(64) NOT NULL, "root_node" int8, PRIMARY KEY ("id"));`

// Select for top level ranks (Kingdoms)
const KINGDOMS = `SELECT tsn, parent_tsn, complete_name FROM taxonomic_units WHERE rank_id = 10;`
