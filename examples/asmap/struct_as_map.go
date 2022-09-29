package asmap

type StructField string

const (
	BaseStructID StructField = "ID"
	BaseStructTS StructField = "TS"
	Name         StructField = "Name"
	Surname      StructField = "Surname"
	NoTag        StructField = "NoTag"
	Address      StructField = "Address"
	FlatCardNum  StructField = "CardNum"
	FlatBank     StructField = "Bank"
)

func (v *Struct[n]) AsMap() map[StructField]interface{} {
	if v == nil {
		return nil
	}
	m := map[StructField]interface{}{}
	if v.BaseStruct != nil {
		m[BaseStructID] = v.BaseStruct.ID
	}
	if v.BaseStruct != nil && v.BaseStruct.TS != nil {
		m[BaseStructTS] = v.BaseStruct.TS
	}
	m[Name] = v.Name
	m[Surname] = v.Surname
	m[NoTag] = v.NoTag
	if v.Address != nil {
		m[Address] = v.Address.AsMap()
	}
	m[FlatCardNum] = v.Flat.CardNum
	m[FlatBank] = v.Flat.Bank
	return m
}
