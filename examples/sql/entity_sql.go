// Code generated by 'fieldr -type Entity -src ../util/const_template.go -output entity_sql.go -const _upsert -const _selectByID -const _deleteByID -const _pk'; DO NOT EDIT.

package sql

const (
	sql_Upsert     = "INSERT id,name,surname INTO " + tableName + " VALUES ($1,$2,$3) DO ON CONFLICT id UPDATE SET name=$2,surname=$3 RETURNING id"
	sql_selectByID = "SELECT id,name,surname FROM " + tableName + " WHERE id = $1"
	sql_deleteByID = "DELETE FROM " + tableName + " WHERE id = $1"
	entity__pk     = "id"
)
