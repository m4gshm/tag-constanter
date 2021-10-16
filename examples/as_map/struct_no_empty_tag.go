package as_map

//go:generate fieldr -type StructNoEmptyTag -out $GOFILE -wrap -export -AsMap -AsTagMap -noEmptyTag

type StructNoEmptyTag struct {
	ID      int    `toMap:"id"`
	Name    string `toMap:"name"`
	Surname string `toMap:"surname"`
	NoTag   string `toMap:""`
}

type (
	StructNoEmptyTagField    string
	StructNoEmptyTagTag      string
	StructNoEmptyTagTagValue string
)

const (
	StructNoEmptyTagField_ID               = StructNoEmptyTagField("ID")
	StructNoEmptyTagField_Name             = StructNoEmptyTagField("Name")
	StructNoEmptyTagField_Surname          = StructNoEmptyTagField("Surname")
	StructNoEmptyTagField_NoTag            = StructNoEmptyTagField("NoTag")
	StructNoEmptyTagTag_toMap              = StructNoEmptyTagTag("toMap")
	StructNoEmptyTagTagValue_toMap_ID      = StructNoEmptyTagTagValue("id")
	StructNoEmptyTagTagValue_toMap_Name    = StructNoEmptyTagTagValue("name")
	StructNoEmptyTagTagValue_toMap_Surname = StructNoEmptyTagTagValue("surname")
)

func (v *StructNoEmptyTag) AsMap() map[StructNoEmptyTagField]interface{} {
	return map[StructNoEmptyTagField]interface{}{
		StructNoEmptyTagField_ID:      v.ID,
		StructNoEmptyTagField_Name:    v.Name,
		StructNoEmptyTagField_Surname: v.Surname,
		StructNoEmptyTagField_NoTag:   v.NoTag,
	}
}

func (v *StructNoEmptyTag) AsTagMap(tag StructNoEmptyTagTag) map[StructNoEmptyTagTagValue]interface{} {
	switch tag {
	case StructNoEmptyTagTag_toMap:
		return map[StructNoEmptyTagTagValue]interface{}{
			StructNoEmptyTagTagValue_toMap_ID:      v.ID,
			StructNoEmptyTagTagValue_toMap_Name:    v.Name,
			StructNoEmptyTagTagValue_toMap_Surname: v.Surname,
		}
	}
	return nil
}
