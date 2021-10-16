package as_map

//go:generate fieldr -type Struct -out $GOFILE -wrap -export -AsMap -AsTagMap

type Struct struct {
	ID              int    `toMap:"id"`
	Name            string `toMap:"name"`
	Surname         string `toMap:"surname"`
	noExport        string `toMap:"no_export"` //nolint
	NoTag           string `toMap:""`
	IgnoredInTagMap string
}

type (
	StructField    string
	StructTag      string
	StructTagValue string
)

const (
	StructField_ID                = StructField("ID")
	StructField_Name              = StructField("Name")
	StructField_Surname           = StructField("Surname")
	structField_noExport          = StructField("noExport")
	StructField_NoTag             = StructField("NoTag")
	StructField_IgnoredInTagMap   = StructField("IgnoredInTagMap")
	StructTag_toMap               = StructTag("toMap")
	StructTagValue_toMap_ID       = StructTagValue("id")
	StructTagValue_toMap_Name     = StructTagValue("name")
	StructTagValue_toMap_Surname  = StructTagValue("surname")
	structTagValue_toMap_noExport = StructTagValue("no_export")
	StructTagValue_toMap_NoTag    = StructTagValue("NoTag") //empty tag
)

func (v *Struct) AsMap() map[StructField]interface{} {
	return map[StructField]interface{}{
		StructField_ID:              v.ID,
		StructField_Name:            v.Name,
		StructField_Surname:         v.Surname,
		StructField_NoTag:           v.NoTag,
		StructField_IgnoredInTagMap: v.IgnoredInTagMap,
	}
}

func (v *Struct) AsTagMap(tag StructTag) map[StructTagValue]interface{} {
	switch tag {
	case StructTag_toMap:
		return map[StructTagValue]interface{}{
			StructTagValue_toMap_ID:      v.ID,
			StructTagValue_toMap_Name:    v.Name,
			StructTagValue_toMap_Surname: v.Surname,
			StructTagValue_toMap_NoTag:   v.NoTag,
		}
	}
	return nil
}
