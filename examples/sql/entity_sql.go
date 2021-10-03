// Code generated by 'fieldr -type Entity -src ../util/const_template.go -output entity_sql.go -const _upsert -const _selectByID -const _deleteByID -const _pk'; DO NOT EDIT.

package sql

const (
	entityField_ID      = "ID"
	entityField_Name    = "Name"
	entityField_Surname = "Surname"
	entityField_ts      = "ts"

	entityTag_db = "db"

	entityTagValue_db_ID      = "id"
	entityTagValue_db_Name    = "name"
	entityTagValue_db_Surname = "surname"
	entityTagValue_db_ts      = "ts"
)
const (
	sql_Upsert     = "INSERT id,name,surname INTO " + tableName + " VALUES ($1,$2,$3) DO ON CONFLICT id UPDATE SET name=$2,surname=$3 RETURNING id"
	sql_selectByID = "SELECT id,name,surname FROM " + tableName + " WHERE id = $1"
	sql_deleteByID = "DELETE FROM " + tableName + " WHERE id = $1"
	entity__pk     = "id"
)

var (
	entity_Fields = []string{entityField_ID, entityField_Name, entityField_Surname}

	entity_Tags = []string{entityTag_db}

	entity_FieldTags = map[string][]string{
		entityField_ID:      []string{entityTag_db},
		entityField_Name:    []string{entityTag_db},
		entityField_Surname: []string{entityTag_db},
	}

	entity_TagValues = map[string][]string{
		entityTag_db: []string{entityTagValue_db_ID, entityTagValue_db_Name, entityTagValue_db_Surname},
	}

	entity_TagFields = map[string][]string{
		entityTag_db: []string{entityField_ID, entityField_Name, entityField_Surname},
	}

	entity_FieldTagValue = map[string]map[string]string{
		entityField_ID:      map[string]string{entityTag_db: entityTagValue_db_ID},
		entityField_Name:    map[string]string{entityTag_db: entityTagValue_db_Name},
		entityField_Surname: map[string]string{entityTag_db: entityTagValue_db_Surname},
	}
)

func (v *Entity) getFieldValue(field string) interface{} {
	switch field {
	case entityField_ID:
		return v.ID
	case entityField_Name:
		return v.Name
	case entityField_Surname:
		return v.Surname
	}
	return nil
}

func (v *Entity) getFieldValueByTagValue(tag string) interface{} {
	switch tag {
	case entityTagValue_db_ID:
		return v.ID
	case entityTagValue_db_Name:
		return v.Name
	case entityTagValue_db_Surname:
		return v.Surname
	}
	return nil
}

func (v *Entity) getFieldValuesByTag(tag string) []interface{} {
	switch tag {
	case entityTag_db:
		return []interface{}{v.ID, v.Name, v.Surname}
	}
	return nil
}

func (v *Entity) asMap() map[string]interface{} {
	return map[string]interface{}{
		entityField_ID:      v.ID,
		entityField_Name:    v.Name,
		entityField_Surname: v.Surname,
	}
}

func (v *Entity) asTagMap(tag string) map[string]interface{} {
	switch tag {
	case entityTag_db:
		return map[string]interface{}{
			entityTagValue_db_ID:      v.ID,
			entityTagValue_db_Name:    v.Name,
			entityTagValue_db_Surname: v.Surname,
		}
	}
	return nil
}
