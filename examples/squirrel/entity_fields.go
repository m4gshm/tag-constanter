// Code generated by 'fieldr'; DO NOT EDIT.

package squirrel

type (
	entityTag       string
	entityTagValue  string
	entityTagValues []entityTagValue
)

const (
	entityTag_db = entityTag("db")

	entityTagValue_db_ID      = entityTagValue("ID")
	entityTagValue_db_Name    = entityTagValue("NAME")
	entityTagValue_db_Surname = entityTagValue("SURNAME")
	entityTagValue_db_ts      = entityTagValue("TS")
)

var (
	entity_TagValues_db = entityTagValues{entityTagValue_db_ID, entityTagValue_db_Name, entityTagValue_db_Surname}
)

func (v *Entity) getFieldValuesByTag(tag entityTag) []interface{} {
	switch tag {
	case entityTag_db:
		return []interface{}{v.ID, v.Name, v.Surname}
	}
	return nil
}

func (v *Entity) getFieldValuesByTagDb() []interface{} {
	return []interface{}{v.ID, v.Name, v.Surname}
}

func (v entityTagValues) strings() []string {
	strings := make([]string, len(v))
	for i, val := range v {
		strings[i] = string(val)
	}
	return strings
}