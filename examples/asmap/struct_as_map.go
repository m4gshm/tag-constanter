// Code generated by 'fieldr'; DO NOT EDIT.

package asmap

type (
	StructField string
)

const (
	StructFieldID                  StructField = "ID"
	StructFieldTS                              = "TS"
	StructFieldName                            = "Name"
	StructFieldSurname                         = "Surname"
	structFieldNoExport                        = "noExport"
	StructFieldNoTag                           = "NoTag"
	StructFieldIgnoredInTagMap                 = "IgnoredInTagMap"
	StructFieldAddress                         = "Address"
	StructFieldFlatNoPrefixCardNum             = "FlatNoPrefix.CardNum"
	StructFieldFlatNoPrefixBank                = "FlatNoPrefix.Bank"
	StructFieldFlatPrefixCardNum               = "FlatPrefix.CardNum"
	StructFieldFlatPrefixBank                  = "FlatPrefix.Bank"
)

func (v *Struct) AsMap() map[StructField]interface{} {
	return map[StructField]interface{}{
		StructFieldID:                  v.ID,
		StructFieldTS:                  v.TS,
		StructFieldName:                v.Name,
		StructFieldSurname:             v.Surname,
		StructFieldNoTag:               v.NoTag,
		StructFieldIgnoredInTagMap:     v.IgnoredInTagMap,
		StructFieldAddress:             v.Address.AsMap(),
		StructFieldFlatNoPrefixCardNum: v.FlatNoPrefix.CardNum,
		StructFieldFlatNoPrefixBank:    v.FlatNoPrefix.Bank,
		StructFieldFlatPrefixCardNum:   v.FlatPrefix.CardNum,
		StructFieldFlatPrefixBank:      v.FlatPrefix.Bank,
	}
}
