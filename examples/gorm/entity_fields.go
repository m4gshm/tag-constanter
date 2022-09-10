// Code generated by 'fieldr'; DO NOT EDIT.

package gorm

type (
	EntityCol string
)

const (
	ENTITY_COL_ID         = EntityCol("ID")
	ENTITY_COL_NAME       = EntityCol("NAME")
	ENTITY_COL_SURNAME    = EntityCol("SURNAME")
	ENTITY_COL_UPDATED_AT = EntityCol("UPDATED_AT")

	EntityGormID        = "ID"
	EntityGormName      = "NAME"
	EntityGormSurname   = "SURNAME"
	EntityGormUpdatedAt = "UPDATED_AT"

	EntityJsonID        = "id"
	EntityJsonName      = "name"
	EntityJsonSurname   = "_surname"
	EntityJsonUpdatedAt = "updateAt"

	EntityGormJsonID        = "id"
	EntityGormJsonName      = "NAME"
	EntityGormJsonSurname   = "SURNAME"
	EntityGormJsonUpdatedAt = "updateAt"
)

func EntityCols() []EntityCol {
	return []EntityCol{
		ENTITY_COL_ID,
		ENTITY_COL_NAME,
		ENTITY_COL_SURNAME,
		ENTITY_COL_UPDATED_AT,
	}
}
func (c EntityCol) Field() string {
	switch c {
	case ENTITY_COL_ID:
		return "ID"
	case ENTITY_COL_NAME:
		return "Name"
	case ENTITY_COL_SURNAME:
		return "Surname"
	case ENTITY_COL_UPDATED_AT:
		return "UpdatedAt"
	}
	return ""
}
